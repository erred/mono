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
	"go.seankhliao.com/mono/internal/svc"
	"golang.org/x/oauth2"
)

var (
	_           svc.SHTTP = &Server{}
	earbugState           = "earbug_state"
)

type Server struct {
	log zerolog.Logger

	mux *http.ServeMux

	host          string
	spotifyID     string
	spotifySecret string

	fname string
	Store *earbugv1.Store

	pollInterval time.Duration
	auth         *spotifyauth.Authenticator

	mu     sync.RWMutex
	client *spotify.Client
}

func (s *Server) Register(r svc.Register) error {
	r.Flags.StringVar(&s.fname, "earbug.data", "/var/lib/mono/earbug/earbug.pb", "path to data file")
	r.Flags.DurationVar(&s.pollInterval, "earbug.poll-interval", 5*time.Minute, "polling interval")
	r.Flags.StringVar(&s.host, "earbug.host", "earbug.liao.dev", "host for auth callback")
	r.Flags.StringVar(&s.spotifyID, "spotify.id", "", "spotify client id")
	r.Flags.StringVar(&s.spotifySecret, "spotify.secret", "", "spotify cluent secret")
	return nil
}

func (s *Server) Init(init svc.Init) error {
	s.log = init.Logger

	s.auth = spotifyauth.New(
		spotifyauth.WithRedirectURL("https://"+s.host+"/auth/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadRecentlyPlayed,
		),
		spotifyauth.WithClientID(s.spotifyID),
		spotifyauth.WithClientSecret(s.spotifySecret),
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

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/auth/callback", s.authCallback)

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
