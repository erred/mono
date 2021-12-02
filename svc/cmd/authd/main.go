package main

import (
	"context"
	"errors"
	"flag"
	"net"
	"os"
	"regexp"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/go-logr/logr"
	"go.seankhliao.com/mono/svc/o11y"
	"google.golang.org/grpc"
)

var (
	errNotRegistered = errors.New("not registered")
	errNeedsAuth     = errors.New("need auth")
)

func main() {
	oo := o11y.NewOptions(flag.CommandLine)
	var s Server
	flag.StringVar(&s.configFile, "config", "/authd/config.prototext", "path to config file")
	flag.StringVar(&s.addr, "grpc.addr", ":8080", "listen address")
	flag.StringVar(&s.realm, "realm", "authd", "displayed realm")
	flag.Parse()

	ctx := oo.New()
	s.l = logr.FromContextOrDiscard(ctx).WithName("authd")
	initLog := logr.FromContextOrDiscard(ctx).WithName("init")

	err := s.fromConfig()
	if err != nil {
		initLog.Error(err, "parse config")
		os.Exit(1)
	}

	svr := grpc.NewServer()
	envoy_service_auth_v3.RegisterAuthorizationServer(svr, &s)

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		initLog.Error(err, "listen", "addr", s.addr)
		os.Exit(1)
	}

	errCh := make(chan error)
	go func() {
		errCh <- svr.Serve(lis)
	}()

	select {
	case <-errCh:
	case <-ctx.Done():
	}

	svr.GracefulStop()
}

type Server struct {
	addr       string
	configFile string
	realm      string

	allow   map[string][]*regexp.Regexp  // host: path regex
	tokens  map[string]map[string]string // host: token: id
	passwds map[string][]byte            // username: hashed passwd

	l logr.Logger

	envoy_service_auth_v3.UnimplementedAuthorizationServer
}

func (s *Server) Check(ctx context.Context, r *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	h := r.GetAttributes().GetRequest().GetHttp()
	headers := h.GetHeaders()
	host := h.GetHost()
	path := h.GetPath()

	l := s.l.WithValues("host", host, "path", path)

	ok := s.checkAllowlist(host, path, headers)
	if ok {
		l.Info("allowed", "check", "allowlist")
		return okResponse("anonymous", headers)
	}

	id := s.checkTokens(host, headers)
	if id != "" {
		l.Info("allowed", "check", "token")
		return okResponse(id, headers)
	}

	user := s.checkBasic(headers)
	if user != "" {
		l.Info("allowed", "check", "htpasswd")
		return okResponse(user, headers)
	}

	l.Info("denied")
	return deniedResponse(s.realm)
}
