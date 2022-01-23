package authnb

import (
	"context"
	"flag"
	"fmt"

	"github.com/go-logr/logr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	authnbv1 "go.seankhliao.com/mono/apis/authnb/v1"
	"google.golang.org/grpc"
)

type Server struct {
	l logr.Logger
	t trace.Tracer

	storeURL string
	store    *clientv3.Client

	authnbv1.UnimplementedAuthnBServer
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	flags.StringVar(&s.storeURL, "store.url", "http://etcd:2379", "etcd connection url")
	return &s
}

func (s *Server) RegisterGRPC(ctx context.Context, svr *grpc.Server, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.t = t.Tracer("authnb")
	s.l = l.WithName("authnb")

	authnbv1.RegisterAuthnBServer(svr, s)

	var err error
	s.store, err = clientv3.NewFromURL(s.storeURL)
	if err != nil {
		return fmt.Errorf("connect to store: %w", err)
	}

	return nil
}
