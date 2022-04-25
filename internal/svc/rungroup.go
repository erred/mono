package svc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/mono/internal/httpmid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"
)

type rungroup struct {
	log    zerolog.Logger
	cancel func()
	done   <-chan struct{}
	sigc   <-chan os.Signal
	wg     sync.WaitGroup
}

func newRungroup(log zerolog.Logger) *rungroup {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, unix.SIGINT, unix.SIGTERM)

	return &rungroup{
		log:    log,
		cancel: cancel,
		done:   ctx.Done(),
		sigc:   sigc,
	}
}

func (r *rungroup) add(run, stop func() error) {
	r.wg.Add(2)
	go func() {
		defer r.wg.Done()
		defer r.cancel()

		err := run()
		if err != nil {
			r.log.Err(err).Msg("service exited with error")
		}
	}()
	go func() {
		defer r.wg.Done()
		<-r.done

		err := stop()
		if err != nil {
			r.log.Err(err).Msg("service stopped with error")
		}
	}()
}

func (r *rungroup) wait() error {
	var err error
	select {
	case sig := <-r.sigc:
		err = fmt.Errorf("signal: %v", sig)
	case <-r.done:
		err = fmt.Errorf("early exit")
	}

	r.cancel()

	wait := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(wait)
	}()

	select {
	case <-r.sigc:
		err = fmt.Errorf("shutdown (%w) interrupted", err)
	case <-wait:
		// gracefully shutdown
	}

	return err
}

//
// Wrap a grpc.Server for rungroup
//
type grpcsvr struct {
	host string
	port string
	svr  *grpc.Server
	log  zerolog.Logger
}

func newGrpcsvr(log zerolog.Logger) *grpcsvr {
	return &grpcsvr{
		svr: grpc.NewServer(
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		),
		log: log,
	}
}

func (s *grpcsvr) register(fset *flag.FlagSet, prefixProto bool) {
	host, port := "host", "port"
	if prefixProto {
		host, port = "grpc.host", "grpc.port"
	}
	fset.StringVar(&s.host, host, "", "grpc listen host (ip address)")
	fset.StringVar(&s.port, port, "8080", "grpc listen port")
}

func (s *grpcsvr) start() error {
	addr := s.host + ":" + s.port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	s.log.Info().Str("addr", addr).Str("protocol", "grpc").Msg("listening")
	return s.svr.Serve(lis)
}

func (s *grpcsvr) stop() error {
	s.svr.GracefulStop()
	return nil
}

//
// Wrap an http.Server for rungroup
//
type httpsvr struct {
	host string
	port string
	svr  *http.Server
	log  zerolog.Logger
	pub  interface{ Publish(string, any) error }
}

func newHttpsvr(log zerolog.Logger) *httpsvr {
	svr := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		// kills connections between lb and us
		// IdleTimeout:       120 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return &httpsvr{
		svr: svr,
		log: log,
	}
}

func (s *httpsvr) register(fset *flag.FlagSet, prefixProto bool) {
	host, port := "host", "port"
	if prefixProto {
		host, port = "http.host", "http.port"
	}
	fset.StringVar(&s.host, host, "", "http listen host (ip address)")
	fset.StringVar(&s.port, port, "8080", "http listen port")
}

func (s *httpsvr) sethandler(handler http.Handler) {
	handler = httpmid.AccessLog(handler, httpmid.AccessLogOut{
		Log: s.log,
		Pub: s.pub,
	})
	handler = otelhttp.NewHandler(handler, "serve")
	handler = h2c.NewHandler(handler, &http2.Server{})
	s.svr.Handler = handler
}

func (s *httpsvr) start() error {
	addr := net.JoinHostPort(s.host, s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	s.log.Info().Str("addr", addr).Str("protocol", "http").Msg("listening")
	err = s.svr.Serve(lis)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *httpsvr) stop() error {
	return s.svr.Shutdown(context.Background())
}
