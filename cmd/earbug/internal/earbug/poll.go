package earbug

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zmb3/spotify/v2"
	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/durationpb"
)

func (s *Server) start() error {
	var token oauth2.Token
	err := json.Unmarshal(s.Store.Token, &token)
	if err != nil {
		return fmt.Errorf("unmarshal stored token: %w", err)
	}

	go func() {
		t := time.NewTicker(s.pollInterval)
		defer t.Stop()
		for {
			s.update()
			<-t.C
		}
	}()

	return nil
}

func (s *Server) update() {
	s.mu.RLock()
	items, err := s.client.PlayerRecentlyPlayedOpt(
		context.Background(),
		&spotify.RecentlyPlayedOptions{
			Limit: 50, // Max
		},
	)
	s.mu.RUnlock()
	if err != nil {
		u := s.auth.AuthURL(earbugState)
		s.log.Err(err).Str("auth_url", u).Msg("get recently played")
		return
	}

	oldPlaybacks, oldTracks := len(s.Store.Playbacks), len(s.Store.Tracks)
	for _, item := range items {
		ts := item.PlayedAt.Format(time.RFC3339Nano)
		if _, ok := s.Store.Playbacks[ts]; !ok {
			s.Store.Playbacks[ts] = &earbugv1.Playback{
				TrackId:     item.Track.ID.String(),
				TrackUri:    string(item.Track.URI),
				ContextType: item.PlaybackContext.Type,
				ContextUri:  string(item.PlaybackContext.URI),
			}
		}

		if _, ok := s.Store.Tracks[item.Track.ID.String()]; !ok {
			t := &earbugv1.Track{
				Id:       item.Track.ID.String(),
				Uri:      string(item.Track.URI),
				Type:     item.Track.Type,
				Name:     item.Track.Name,
				Duration: durationpb.New(item.Track.TimeDuration()),
			}
			for _, artist := range item.Track.Artists {
				t.Artists = append(t.Artists, &earbugv1.Artist{
					Id:   artist.ID.String(),
					Uri:  string(artist.URI),
					Name: artist.Name,
				})
			}
			s.Store.Tracks[item.Track.ID.String()] = t
		}
	}

	newPlaybacks, newTracks := len(s.Store.Playbacks), len(s.Store.Tracks)

	if (newPlaybacks+newTracks)-(oldTracks+oldPlaybacks) > 0 {
		s.log.Debug().
			Int("tracks_total", newTracks).Int("tracks_new", newTracks-oldTracks).
			Int("plays_total", newPlaybacks).Int("plays_new", newPlaybacks-oldPlaybacks).
			Msg("new")
		err = s.Write()
		if err != nil {
			s.log.Err(err).Msg("write to store")
		}
	}
}
