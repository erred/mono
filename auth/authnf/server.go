package authnf

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	authnbv1 "go.seankhliao.com/mono/apis/authnb/v1"
	"google.golang.org/grpc"
)

type Server struct {
	l logr.Logger
	t trace.Tracer

	cookieDomain string
	cookieName   string
	cookieTTL    time.Duration

	authnbURL string
	authnb    authnbv1.AuthnBClient
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	flags.StringVar(&s.cookieDomain, "cookie.domain", "seankhliao.com", "domain for cookie")
	flags.StringVar(&s.cookieName, "cookie.name", "__authnf_session", "cookie name")
	flags.DurationVar(&s.cookieTTL, "cookie.ttl", 24*time.Hour, "cookie valid time in seconds")
	flags.StringVar(&s.authnbURL, "authnb.url", "authnb.auth.svc:80", "authhnb connection url")
	return &s
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("authnf")
	s.t = t.Tracer("authnf")

	conn, err := grpc.DialContext(ctx, s.authnbURL,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("grpc client conn: %w", err)
	}
	s.authnb = authnbv1.NewAuthnBClient(conn)

	mux.HandleFunc("/logout", s.handleLogout)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/", s.handleIndex)

	return nil
}
