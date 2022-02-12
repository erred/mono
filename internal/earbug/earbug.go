package earbug

import (
	"net/http"
	"os"
	"strings"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"go.seankhliao.com/mono/internal/envconf"
)

type Server struct {
	host string

	fname string
	Store *earbugv1.Store

	pollInterval time.Duration
	auth         *spotifyauth.Authenticator
}

func New() (*Server, error) {
	pollInterval, err := time.ParseDuration(envconf.String("EARBUG_POLL_INTERVAL", "5m"))
	if err != nil {
		return nil, err
	}

	s := &Server{
		host:         envconf.String("EARBUG_HOST", "earbug.seankhliao.com"),
		fname:        envconf.String("EARBUG_DATA", "/data/earbug.pb"),
		Store:        &earbugv1.Store{},
		pollInterval: pollInterval,
	}

	s.auth = spotifyauth.New(
		spotifyauth.WithRedirectURL("https://"+s.host+"/auth/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadRecentlyPlayed,
		),
		spotifyauth.WithClientID(strings.TrimSpace(os.Getenv("SPOTIFY_ID"))),
		spotifyauth.WithClientSecret(strings.TrimSpace(os.Getenv("SPOTIFY_SECRET"))),
	)

	err = s.initStore()
	if err != nil {
		return nil, err
	}

	err = s.start()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotImplemented)
}
