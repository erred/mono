package runhttp

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/mono/monolith/component"
	"go.seankhliao.com/mono/monolith/o11y"
	"go.seankhliao.com/mono/monolith/run"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var (
	_ run.Runner           = &Server{}
	_ component.Background = &Server{}
)

type Server struct {
	E bool
	M *http.ServeMux
	S *http.Server
	T o11y.Tool
}

func New() *Server {
	return &Server{
		M: http.NewServeMux(),
		S: &http.Server{
			ReadHeaderTimeout: 5 * time.Second,
			IdleTimeout:       120 * time.Second,
			WriteTimeout:      5 * time.Second,
		},
	}
}

func (s *Server) Enabled() bool {
	return s.E
}

func (s *Server) Register(flags *flag.FlagSet) {
	flags.BoolVar(&s.E, "http", true, "enable http server")
	flags.StringVar(&s.S.Addr, "http.addr", ":8080", "http listen address")
}

func (s *Server) Background(ctx context.Context, tp o11y.ToolProvider) {
	s.T = tp.Tool("runhttp")

	handler := otelhttp.NewHandler(s.M, "handle")
	var h2 http2.Server
	s.S.Handler = h2c.NewHandler(handler, &h2)
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.S.Shutdown(context.Background())
	}()

	s.T.Info("starting http", "addr", s.S.Addr)
	err := s.S.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
