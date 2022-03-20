package httpsvc

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/envconf"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sys/unix"
)

type HTTPSvc interface {
	Init(zerolog.Logger) error
	Desc() string
	Help() string
	http.Handler
}

func help() string {
	return `
HOST
        ip address to listen on
PORT
        port to listen on over http/h2c
`
}

func Run(s HTTPSvc) {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s %s\n", os.Args[0], s.Desc())
		for _, str := range []string{help(), s.Help()} {
			fmt.Fprintln(os.Stderr, strings.TrimSpace(str))
		}
	}
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		log.Error().Strs("args", args).Msg("unexpected arguments")
		os.Exit(1)
	}

	s.Init(log)

	host := envconf.String("HOST", "")
	port := envconf.String("PORT", "8080")
	handler := h2c.NewHandler(s, &http2.Server{})
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
