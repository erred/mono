package httpsvc

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	natsproto "github.com/nats-io/nats.go/encoders/protobuf"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.seankhliao.com/mono/internal/flagwrap"
	"go.seankhliao.com/mono/internal/httpmid"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sys/unix"
)

type Init struct {
	Flags      *flag.FlagSet
	FlagsAfter func() error
	Log        zerolog.Logger
}

type HTTPSvc interface {
	Init(*Init) error
	http.Handler
}

func Run(s HTTPSvc) {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	fset := flag.NewFlagSet("", flag.ContinueOnError)

	init := Init{
		Flags: fset,
		Log:   log,
	}
	err := s.Init(&init)
	if err != nil {
		log.Err(err).Msg("svc initialize")
		os.Exit(1)
	}

	var host, port, natsURL string
	fset.StringVar(&host, "host", "", "ip address to listen on")
	fset.StringVar(&port, "port", "8080", "port to listen on")
	fset.StringVar(&natsURL, "nats.url", "nats://localhost:4222", "NATS endpoint for reporting")
	err = flagwrap.Parse(fset, os.Args[1:])
	if errors.Is(err, flag.ErrHelp) {
		os.Exit(0)
	} else if err != nil {
		log.Err(err).Msg("parse config")
		os.Exit(1)
	}
	if fset.NArg() != 0 {
		log.Error().Strs("args", fset.Args()).Msg("unexpected args")
		os.Exit(1)
	}
	if init.FlagsAfter != nil {
		err := init.FlagsAfter()
		if err != nil {
			log.Err(err).Msg("svc parse hook")
			os.Exit(1)
		}
	}

	accessLogOpts := httpmid.AccessLogOut{
		Log: log,
	}

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Info().Err(err).Msg("skipping nats setup")
	} else {
		natsConnEncoded, err := nats.NewEncodedConn(natsConn, natsproto.PROTOBUF_ENCODER)
		if err != nil {
			log.Info().Err(err).Msg("setup nats encoded conn")
		} else {
			accessLogOpts.Pub = natsConnEncoded
		}
	}

	handler := http.Handler(s)
	handler = httpmid.AccessLog(handler, accessLogOpts)
	handler = otelhttp.NewHandler(handler, "serve")
	handler = h2c.NewHandler(handler, &http2.Server{})

	svr := http.Server{
		Addr:              host + ":" + port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		log.Info().Str("addr", svr.Addr).Msg("listening")
		svr.ListenAndServe()
	}()
	go func() {
		<-ctx.Done()
		svr.Shutdown(context.Background())
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, unix.SIGINT, unix.SIGTERM)
	select {
	case <-sigc:
		// external interrupt
	case <-ctx.Done():
		// one of the services exited
	}

	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-sigc:
		// interrupted shutdown
	case <-done:
		// gracefully shutdown
	}
}
