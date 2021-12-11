package runsvr

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
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/internal/stdlog"
	"google.golang.org/grpc"
)

type Runner struct {
	// set via flags
	otlpEndpoint string
	svcAddr      string

	// set in init
	sigc chan os.Signal
	stop chan struct{}
	l    logr.Logger
	t    trace.TracerProvider
	m    metric.MeterProvider

	// set in Http/Grpc
	serve     func(net.Listener) error
	shutdowns []func(context.Context) error
}

func New(flags *flag.FlagSet) *Runner {
	var r Runner

	addr, port := ":8080", os.Getenv("PORT")
	if port != "" {
		addr = ":" + port
	}

	flags.StringVar(&r.otlpEndpoint, "otlp.endpoint", "", "otlp grpc endpont for metrics/traces")
	flags.StringVar(&r.svcAddr, "svc.addr", addr, "address to listen on, flag > PORT > :8080")

	return &r
}

func (r *Runner) init() {
	r.sigc = make(chan os.Signal, 2)
	signal.Notify(r.sigc, syscall.SIGINT, syscall.SIGTERM)
	r.stop = make(chan struct{})

	r.l = funcr.NewJSON(func(obj string) {
		fmt.Fprintln(os.Stderr, obj)
	}, funcr.Options{
		LogTimestamp:    true,
		TimestampFormat: time.RFC3339,
	})

	l := r.l.WithName("runsvr")

	// opentelemetry setup

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		r.l.WithName("otel").Error(errors.New("otel"), err.Error())
	}))

	if r.otlpEndpoint == "" {
		// use global noop
		r.t = otel.GetTracerProvider()
		r.m = global.GetMeterProvider()
		return
	}

	ctx := context.TODO()
	res, err := otelResources(ctx)
	if err != nil {
		l.Error(err, "init otel resources")
		os.Exit(1)
	}
	otlpConn, err := otelGRPC(ctx, r.otlpEndpoint)
	if err != nil {
		l.Error(err, "connect to otlp collector", "endpoint", r.otlpEndpoint)
		os.Exit(1)
	}

	err = otelTracer(ctx, res, otlpConn)
	if err != nil {
		l.Error(err, "setup tracer")
		os.Exit(1)
	}

	err = otelMetrics(ctx, res, otlpConn)
	if err != nil {
		l.Error(err, "setup metrics")
		os.Exit(1)
	}

	err = otelruntime.Start(otelruntime.WithMeterProvider(r.m))
	if err != nil {
		l.Error(err, "setup otel runtime instrumentation")
		os.Exit(1)
	}
}

type HTTPService interface {
	RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error
}

func (r *Runner) HTTP(s HTTPService) {
	r.init()
	ctx := context.TODO()

	mux := http.NewServeMux()
	handler := otelhttp.NewHandler(
		mux,
		"handle",
		otelhttp.WithMeterProvider(r.m),
		otelhttp.WithTracerProvider(r.t),
	)

	svr := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     stdlog.New(r.l.WithName("http"), errors.New("net/http")),
	}

	r.serve = svr.Serve
	r.shutdowns = append(r.shutdowns, svr.Shutdown)

	err := s.RegisterHTTP(ctx, mux, r.l, r.m, r.t, func() { close(r.stop) })
	if err != nil {
		r.l.WithName("runsvr").Error(err, "registering service", "protocol", "http")
		os.Exit(1)
	}

	r.run()
}

type GRPCService interface {
	RegisterGRPC(ctx context.Context, svr *grpc.Server, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error
}

func (r *Runner) GRPC(s GRPCService) {
	r.init()
	ctx := context.TODO()

	svr := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor(otelgrpc.WithTracerProvider(r.t))),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor(otelgrpc.WithTracerProvider(r.t))),
	)
	r.serve = svr.Serve
	r.shutdowns = append(r.shutdowns, func(context.Context) error { svr.GracefulStop(); return nil })

	err := s.RegisterGRPC(ctx, svr, r.l, r.m, r.t, func() { close(r.stop) })
	if err != nil {
		r.l.WithName("runsvr").Error(err, "registering service", "protocol", "grpc")
		os.Exit(1)
	}

	r.run()
}

func (r *Runner) run() {
	l := r.l.WithName("runsvr")

	family := "tcp"
	lis, err := net.Listen(family, r.svcAddr)
	if err != nil {
		l.Error(err, "listen", "family", family, "address", r.svcAddr)
		os.Exit(1)
	}
	l.Info("listening", "address", r.svcAddr)

	errc := make(chan error)
	go func() { errc <- r.serve(lis) }()

	select {
	case err := <-errc:
		l.Error(err, "unexpected server exit")
	case <-r.stop:
		l.Info("service requested shutdown")
	case sig := <-r.sigc:
		l.Info("shutting down due to signal", "signal", sig)
	}

	var wg sync.WaitGroup
	for _, sh := range r.shutdowns {
		wg.Add(1)
		go func(sh func(context.Context) error) {
			defer wg.Done()
			err := sh(context.Background())
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				l.Error(err, "shutting down component")
			}
		}(sh)
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case sig := <-r.sigc:
		l.Error(errors.New("got signal"), "shutdown sequence interrupted", "signal", sig)
		os.Exit(1)
	case <-done:
	}
}
