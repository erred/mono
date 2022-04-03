package earbug

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"go.seankhliao.com/mono/internal/httpsvc"
	"golang.org/x/oauth2"
)

var (
	_           httpsvc.HTTPSvc = &Server{}
	earbugState                 = "earbug_state"
)

type Server struct {
	log zerolog.Logger

	mux *http.ServeMux

	fname string
	Store *earbugv1.Store

	pollInterval time.Duration
	auth         *spotifyauth.Authenticator

	mu     sync.RWMutex
	client *spotify.Client
}

func (s *Server) Init(init *httpsvc.Init) error {
	s.log = init.Log

	init.Flags.StringVar(&s.fname, "earbug.data", "/var/lib/mono/earbug/earbug.pb", "path to data file")
	init.Flags.DurationVar(&s.pollInterval, "earbug.poll-interval", 5*time.Minute, "polling interval")
	var host, spotifyID, spotifySecret string
	init.Flags.StringVar(&host, "earbug.host", "earbug.liao.dev", "host for auth callback")
	init.Flags.StringVar(&spotifyID, "spotify.id", "", "spotify client id")
	init.Flags.StringVar(&spotifySecret, "spotify.secret", "", "spotify cluent secret")

	init.FlagsAfter = func() error {
		s.auth = spotifyauth.New(
			spotifyauth.WithRedirectURL("https://"+host+"/auth/callback"),
			spotifyauth.WithScopes(
				spotifyauth.ScopeUserReadRecentlyPlayed,
			),
			spotifyauth.WithClientID(spotifyID),
			spotifyauth.WithClientSecret(spotifySecret),
		)

		err := s.initStore()
		if err != nil {
			return fmt.Errorf("init store: %w", err)
		}

		var token oauth2.Token
		err = json.Unmarshal(s.Store.Token, &token)
		if err != nil {
			return fmt.Errorf("unmarshal stored token: %w", err)
		}

		s.setClient(&token)

		err = s.start()
		if err != nil {
			return fmt.Errorf("start background: %w", err)
		}
		return nil
	}

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/auth/callback", s.authCallback)
	s.Store = &earbugv1.Store{}

	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.log.Debug().Str("user_agent", r.UserAgent()).Str("referrer", r.Referer()).Str("url", r.URL.String()).Msg("requested")
	s.mux.ServeHTTP(rw, r)
}

func (s *Server) authCallback(rw http.ResponseWriter, r *http.Request) {
	token, err := s.auth.Token(r.Context(), earbugState, r)
	if err != nil {
		s.log.Err(err).Msg("get token")
		http.Error(rw, "token err", http.StatusInternalServerError)
		return
	}

	s.setClient(token)
}

func (s *Server) setClient(token *oauth2.Token) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.client = spotify.New(
		s.auth.Client(
			context.Background(),
			token,
		),
		spotify.WithRetry(true),
	)
}
