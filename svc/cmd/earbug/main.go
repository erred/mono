package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

func main() {
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	var server Server
	flag.StringVar(&server.CanonicalURL, "url", "https://earbug.seankhliao.com", "url app is hosted on")
	flag.StringVar(&server.StoreURL, "store", "http://etcd-0.etcd.etcd.svc.cluster.local:2379", "etcd url")
	flag.StringVar(&server.StorePrefix, "store-prefix", "earbug", "key prefix in etcd")
	flag.DurationVar(&server.PollInterval, "poll-interval", 5*time.Minute, "time between spotify polls")
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = server.Handler(ctx)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "init")
		os.Exit(1)
	}

	defer server.pollWorkerWg.Wait()
	defer server.Store.Close()
	defer func() {
		close(server.pollWorkerShutdown)
	}()

	ho.Run()
}

type Server struct {
	CanonicalURL string
	StoreURL     string
	StorePrefix  string
	PollInterval time.Duration

	Auth  *spotifyauth.Authenticator
	Store *clientv3.Client

	pollWorkerShutdown chan struct{}
	pollWorkerWg       sync.WaitGroup
}

func (s *Server) Handler(ctx context.Context) (http.Handler, error) {
	s.Auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(s.CanonicalURL+"/auth/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadRecentlyPlayed,
		),
		spotifyauth.WithClientID(strings.TrimSpace(os.Getenv("SPOTIFY_ID"))),
		spotifyauth.WithClientSecret(strings.TrimSpace(os.Getenv("SPOTIFY_SECRET"))),
	)

	s.pollWorkerShutdown = make(chan struct{})

	var err error
	s.Store, err = clientv3.NewFromURL(s.StoreURL)
	if err != nil {
		return nil, err
	}

	err = s.startStoredPoll(ctx)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/user/", s.authPage)
	mux.HandleFunc("/auth/callback", s.authCallback)
	mux.HandleFunc("/", s.index)
	return mux, nil
}

func (s *Server) index(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := logr.FromContextOrDiscard(ctx).WithName("index")

	if r.URL.Path != "/" {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	ro := &render.Options{
		Data: render.PageData{
			Title: "earbug",
			H1:    "earbug",
			H2:    "spotify logger",
		},
	}
	err := render.Render(ro, rw, strings.NewReader(indexMsg))
	if err != nil {
		l.Error(err, "render index page")
	}
}

const indexMsg = `
# _earbug_

A simple spotify history logger
`
