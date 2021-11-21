package httpsvr

import (
	"context"
	"errors"
	"flag"
	"net"
	"net/http"
	"net/http/pprof"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/mono/internal/stdlog"
)

type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

type Options struct {
	BaseContext context.Context
	Handler     http.Handler
	Shutdowners []Shutdowner

	AdminAddr         string
	HTTPAddr          string
	HTTPServerKeyFile string
	HTTPServerCrtFile string
	HTTPClientCAFile  string
	ShutdownGrace     time.Duration
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.AdminAddr, "admin.addr", "127.0.0.1:8090", "admin listen host:port")
	fs.StringVar(&o.HTTPAddr, "http.addr", ":8080", "http listen host:port")
	fs.StringVar(&o.HTTPServerKeyFile, "http.serverkey", "", "path to https server tls key file")
	fs.StringVar(&o.HTTPServerCrtFile, "http.servercrt", "", "path to https server tls cert file")
	fs.StringVar(&o.HTTPClientCAFile, "http.clientca", "", "path to https client ca cert file")
	fs.DurationVar(&o.ShutdownGrace, "lifecycle.shutdown.grace", 10*time.Second, "time to allow for graceful shutdown")
	return &o
}

func (o *Options) Run() {
	l := logr.FromContextOrDiscard(o.BaseContext).WithName("httpsvr")

	admsvr := &http.Server{
		Addr:              o.AdminAddr,
		Handler:           o.adminHandler(),
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 10,
		ErrorLog:          stdlog.New(l.WithName("adm"), errors.New("http.Server")),
	}

	httpsvr := &http.Server{
		Addr:              o.HTTPAddr,
		Handler:           o.defaultHandler(),
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20,
		ErrorLog:          stdlog.New(l.WithName("http"), errors.New("http.Server")),
		BaseContext: func(net.Listener) context.Context {
			return logr.NewContext(context.Background(), l.WithName("app"))
		},
	}

	o.Shutdowners = append(o.Shutdowners, admsvr, httpsvr)

	ctx, cancel := context.WithCancel(o.BaseContext)
	defer cancel()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(len(o.Shutdowners) + 2)
	defer wg.Wait()

	for _, shutdowner := range o.Shutdowners {
		go waitShutdown(l, ctx, shutdowner, o.ShutdownGrace, &wg)
	}

	go runServer(l.WithValues("server", "adm"), admsvr, &wg, cancel)
	go runServer(l.WithValues("server", "http"), httpsvr, &wg, cancel)
}

func runServer(l logr.Logger, svr *http.Server, wg *sync.WaitGroup, cancel func()) {
	defer wg.Done()
	defer cancel()

	l.Info("starting", "addr", svr.Addr, "tls", svr.TLSConfig != nil)
	err := svr.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error(err, "shutdown")
	}
}

func waitShutdown(l logr.Logger, ctx context.Context, s Shutdowner, timeout time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.Shutdown(ctx)
	if err != nil {
		l.Error(err, "graceful shutdown")
	}
}

// defaultHandler wraps the provided handler in observability (o11y)
// and sets default headers.
func (o *Options) defaultHandler() http.Handler {
	serverName := "go.seankhliao.com/mono/svc/httpsvr"
	bi, ok := debug.ReadBuildInfo()
	if ok {
		serverName = bi.Path
	}
	return otelhttp.NewHandler(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("server", serverName)

			o.Handler.ServeHTTP(rw, r)
		}),
		"app",
	)
}

// adminHandler returns a handler with pprof and readiness endpoints setup
func (o *Options) adminHandler() http.Handler {
	mux := http.NewServeMux()

	// pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// healthchecks
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return mux
}
