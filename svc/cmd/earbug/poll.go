package main

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/go-logr/logr"
	"github.com/zmb3/spotify/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/oauth2"
)

// startStoredPoll starts pollers for all stored tokens
func (s *Server) startStoredPoll(ctx context.Context) error {
	l := logr.FromContextOrDiscard(ctx).WithName("poller")

	p := path.Join(s.StorePrefix, "token")
	res, err := s.Store.Get(ctx, p, clientv3.WithPrefix())
	if err != nil {
		l.Error(err, "get tokens", "prefix", p)
		return err
	}

	for _, kv := range res.Kvs {
		user := path.Base(string(kv.Key))
		var tok oauth2.Token
		err = json.Unmarshal(kv.Value, &tok)
		if err != nil {
			l.Error(err, "unmarshal token", "user", user)
			return err
		}
		s.addPollWorker(ctx, user, &tok)
	}
	return nil
}

// addPollWorker starts a poll worker for the user.
// If token is nil, it is retrieved from the db.
func (s *Server) addPollWorker(ctx context.Context, user string, token *oauth2.Token) {
	s.pollWorkerWg.Add(1)
	go s.pollUser(ctx, user, token)
}

// pollUser is a poll worker responsible for updating a user's listening history
func (s *Server) pollUser(ctx context.Context, user string, token *oauth2.Token) {
	defer s.pollWorkerWg.Done()
	l := logr.FromContextOrDiscard(ctx).WithName("poller")
	l.Info("starting", "user", user)

	c := spotify.New(s.Auth.Client(ctx, token), spotify.WithRetry(true))
	t := time.NewTicker(s.PollInterval)
	defer t.Stop()

	for {
		s.updateHistory(ctx, user, c)
		select {
		case <-t.C:
		case <-s.pollWorkerShutdown:
			return
		}
	}
}

// updateHistory pulls a user's recently played tracks and stores it
func (s *Server) updateHistory(ctx context.Context, user string, c *spotify.Client) {
	l := logr.FromContextOrDiscard(ctx).WithName("poll").WithValues("user", user)
	ctx = logr.NewContext(ctx, l)

	items, err := c.PlayerRecentlyPlayedOpt(ctx, &spotify.RecentlyPlayedOptions{
		Limit: 50, // Max
	})
	if err != nil {
		l.Error(err, "get recently played")
		return
	}
	for _, item := range items {
		s.putHistory(ctx, user, item)
	}
}

// putHistory stores a single user listen history
func (s *Server) putHistory(ctx context.Context, user string, item spotify.RecentlyPlayedItem) {
	l := logr.FromContextOrDiscard(ctx)

	playedP := path.Join(s.StorePrefix, "history", user, "playback", item.PlayedAt.Format(time.RFC3339Nano))
	playedB, err := json.Marshal(map[string]interface{}{
		"track_id":  item.Track.ID,
		"track_uri": item.Track.URI,
		"context":   item.PlaybackContext,
	})
	if err != nil {
		l.Error(err, "marshal playback")
		return
	}
	_, err = s.Store.Put(ctx, playedP, string(playedB))
	if err != nil {
		l.Error(err, "put user playback", "key", playedP)
	}

	uniqTrackP := path.Join(s.StorePrefix, "history", user, "track", string(item.Track.ID))
	_, err = s.Store.Put(ctx, uniqTrackP, item.PlayedAt.Format(time.RFC3339Nano))
	if err != nil {
		l.Error(err, "put user track", "key", uniqTrackP)
	}

	trackP := path.Join(s.StorePrefix, "track", string(item.Track.ID))
	trackB, err := json.Marshal(item.Track)
	if err != nil {
		l.Error(err, "marshal track info")
	}
	_, err = s.Store.Put(ctx, trackP, string(trackB))
	if err != nil {
		l.Error(err, "put global track info", "key", trackP)
	}
}
