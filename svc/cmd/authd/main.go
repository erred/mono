package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"regexp"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/svc/runsvr"
	"google.golang.org/grpc"
)

var (
	errNotRegistered = errors.New("not registered")
	errNeedsAuth     = errors.New("need auth")
)

func main() {
	s := New(flag.CommandLine)
	r := runsvr.New(flag.CommandLine)
	flag.Parse()

	r.GRPC(s)
}

type Server struct {
	// from flags
	configFile string
	realm      string

	// from config file
	allow   map[string][]*regexp.Regexp  // host: path regex
	tokens  map[string]map[string]string // host: token: id
	passwds map[string][]byte            // username: hashed passwd

	// from runsvr
	l logr.Logger
	t trace.Tracer

	// grpc
	envoy_service_auth_v3.UnimplementedAuthorizationServer
}

func New(flags *flag.FlagSet) *Server {
	var s Server

	flags.StringVar(&s.configFile, "config", "/authd/config.prototext", "path to config file")
	flags.StringVar(&s.realm, "realm", "authd", "displayed realm")

	return &s
}

func (s *Server) RegisterGRPC(ctx context.Context, svr *grpc.Server, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("authd")
	s.t = t.Tracer("authd")

	envoy_service_auth_v3.RegisterAuthorizationServer(svr, s)

	err := s.fromConfig()
	if err != nil {
		return fmt.Errorf("from config=%s: %w", s.configFile, err)
	}

	return nil
}
