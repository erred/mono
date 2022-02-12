package rungrpc

import (
	"context"
	"flag"
	"net"

	"go.seankhliao.com/mono/monolith/component"
	"go.seankhliao.com/mono/monolith/o11y"
	"go.seankhliao.com/mono/monolith/run"
	"google.golang.org/grpc"
)

var (
	_ run.Runner           = &Server{}
	_ component.Background = &Server{}
)

type Server struct {
	Addr string
	E    bool
	S    *grpc.Server
	T    o11y.Tool
}

func New() *Server {
	return &Server{}
}

func (s *Server) Enabled() bool {
	return s.E
}

func (s *Server) Register(flags *flag.FlagSet) {
	flags.BoolVar(&s.E, "grpc", true, "enable grpc server")
	flags.StringVar(&s.Addr, "grpc.addr", ":8000", "grpc listen address")
}

func (s *Server) Background(ctx context.Context, tp o11y.ToolProvider) {
	s.S = grpc.NewServer()
	s.T = tp.Tool("rungrpc")
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.S.GracefulStop()
	}()

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.T.Info("starting grpc", "addr", s.Addr)
	return s.S.Serve(lis)
}
