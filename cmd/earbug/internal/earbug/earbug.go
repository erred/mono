package earbug

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"go.seankhliao.com/mono/internal/envconf"
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

func (s *Server) Desc() string {
	return `spotify listening history logger`
}

func (s *Server) Help() string {
	return `
EARBUG_DATA
        path to file to store data
EARBUG_HOST
        hostname (for auth callback)
EARBUG_POLL_INTERVAL
        interval between polls
SPOTIFY_ID
        spotify client id
SPOTIFY_SECRET
        spotify client secret
`
}

func (s *Server) Init(log zerolog.Logger) error {
	s.log = log

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/auth/callback", s.authCallback)
	var err error
	pollInterval := envconf.String("EARBUG_POLL_INTERVAL", "5m")
	s.pollInterval, err = time.ParseDuration(pollInterval)
	if err != nil {
		return fmt.Errorf("parse EARBUG_POLL_INTERVAL=%q:  %w", pollInterval, err)
	}

	s.fname = envconf.String("EARBUG_DATA", "/var/lib/earbug/earbug.pb")

	s.Store = &earbugv1.Store{}

	host := envconf.String("EARBUG_HOST", "earbug.liao.dev")
	s.auth = spotifyauth.New(
		spotifyauth.WithRedirectURL("https://"+host+"/auth/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadRecentlyPlayed,
		),
		spotifyauth.WithClientID(strings.TrimSpace(os.Getenv("SPOTIFY_ID"))),
		spotifyauth.WithClientSecret(strings.TrimSpace(os.Getenv("SPOTIFY_SECRET"))),
	)

	err = s.initStore()
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
