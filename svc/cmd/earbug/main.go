package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/runsvr"
)

func main() {
	r := runsvr.New(flag.CommandLine)
	runsvr.Desc(flag.CommandLine, docgo)
	s := New(flag.CommandLine)
	flag.Parse()

	r.HTTP(s)
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	s.pollWorkerMap = make(map[string]struct{})
	flag.StringVar(&s.CanonicalURL, "url", "https://earbug.seankhliao.com", "url app is hosted on")
	flag.StringVar(&s.StoreURL, "store", "http://etcd-0.etcd:2379", "etcd url")
	flag.StringVar(&s.StorePrefix, "store-prefix", "earbug", "key prefix in etcd")
	flag.DurationVar(&s.PollInterval, "poll-interval", 5*time.Minute, "time between spotify polls")
	return &s
}

type Server struct {
	CanonicalURL string
	StoreURL     string
	StorePrefix  string
	PollInterval time.Duration

	l logr.Logger
	t trace.Tracer

	Auth  *spotifyauth.Authenticator
	Store *clientv3.Client

	indexPage []byte
	startTime time.Time

	pollWorkerShutdown chan struct{}
	pollWorkerWg       sync.WaitGroup
	pollWorkerMap      map[string]struct{}
	pollWorkerMu       sync.Mutex
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("earbug")
	s.t = t.Tracer("earbug")

	s.startTime = time.Now()
	var err error
	s.indexPage, err = render.CompactBytes(
		"earbug",
		"spotify listen tracker",
		"https://earbug.seankhliao.com/",
		[]byte(indexMsg),
	)
	if err != nil {
		return fmt.Errorf("render index: %w", err)
	}

	s.Auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(s.CanonicalURL+"/auth/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadRecentlyPlayed,
		),
		spotifyauth.WithClientID(strings.TrimSpace(os.Getenv("SPOTIFY_ID"))),
		spotifyauth.WithClientSecret(strings.TrimSpace(os.Getenv("SPOTIFY_SECRET"))),
	)

	s.pollWorkerShutdown = make(chan struct{})

	s.Store, err = clientv3.NewFromURL(s.StoreURL)
	if err != nil {
		return err
	}

	err = s.startStoredPoll(ctx)
	if err != nil {
		return err
	}

	mux.HandleFunc("/auth/callback", s.handleAuthCallback)
	mux.HandleFunc("/user", s.handleUser)
	mux.HandleFunc("/user/history", s.handleUserHistory)
	mux.HandleFunc("/", s.index)
	return nil
}

func (s *Server) index(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	http.ServeContent(rw, r, "index.html", s.startTime, bytes.NewReader(s.indexPage))
}

const indexMsg = `
### _earbug_

A simple spotify history logger
`
