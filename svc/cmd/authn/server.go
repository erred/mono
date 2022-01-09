package main

import (
	"context"
	"flag"
	"net/http"
	"sync"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/svc/cmd/authn/store"
)

type Server struct {
	l logr.Logger
	t trace.Tracer

	cookieDomain      string
	cookieName        string
	canonicalHostname string
	hashedPasswordsMu sync.RWMutex
	hashedPasswords   map[string][]byte

	etcdUrl string
	store   *store.Store
}

func New(flags *flag.FlagSet) *Server {
	s := Server{
		hashedPasswords: make(map[string][]byte),
	}

	flags.StringVar(&s.cookieName, "cookie", "__authn_session", "name of cookie")
	flags.StringVar(&s.cookieDomain, "domain", "seankhliao.com", "cookie domain")
	flags.StringVar(&s.canonicalHostname, "hostname", "authn.seankhliao.com", "canonical hostname")
	flags.StringVar(&s.etcdUrl, "etcd", "http://etcd-0.etcd:2379", "etcd session store")

	return &s
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("authn")
	s.t = t.Tracer("authn")

	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/login", s.handleApiLogin)
	mux.HandleFunc("/logout", s.handleApiLogout)

	var err error
	s.store, err = store.New(s.etcdUrl, t)
	if err != nil {
		return err
	}

	return nil
}
