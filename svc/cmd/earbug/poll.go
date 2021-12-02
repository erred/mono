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
		select {
		case <-t.C:
			s.updateHistory(ctx, user, c)
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

	ts := item.PlayedAt.Format(time.RFC3339Nano)
	b, err := json.Marshal(item)
	if err != nil {
		l.Error(err, "marhsal recently played")
		return
	}
	p := path.Join(s.StorePrefix, "history", user, ts)
	_, err = s.Store.Put(ctx, p, string(b))
	if err != nil {
		l.Error(err, "put recently played", "path", p)
		return
	}
}
