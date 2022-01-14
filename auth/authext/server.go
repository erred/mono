package authext

import (
	"context"
	"flag"
	"fmt"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/auth/authnbpb"
	"google.golang.org/grpc"
)

type Server struct {
	l logr.Logger
	t trace.Tracer

	headerID    string
	redirectURL string
	cookieName  string

	authnbURL string
	authnb    authnbpb.AuthnBClient

	// grpc
	envoy_service_auth_v3.UnimplementedAuthorizationServer
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	flags.StringVar(&s.authnbURL, "authnb.url", "authnb.auth.svc:80", "connection to authn backend")
	flags.StringVar(&s.headerID, "header.id", "auth-id", "header with user id")
	flags.StringVar(&s.redirectURL, "redirect.url", "https://authnf.seankhliao.com/", "redirect for unauthenticated users")
	flags.StringVar(&s.cookieName, "cookie.name", "__authnf_session", "cookie name")
	return &s
}

func (s *Server) RegisterGRPC(ctx context.Context, svr *grpc.Server, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("authext")
	s.t = t.Tracer("authext")

	conn, err := grpc.DialContext(ctx, s.authnbURL,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("grpc client conn: %w", err)
	}
	s.authnb = authnbpb.NewAuthnBClient(conn)

	envoy_service_auth_v3.RegisterAuthorizationServer(svr, s)

	return nil
}
