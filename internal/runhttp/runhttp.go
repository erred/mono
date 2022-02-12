package runhttp

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.seankhliao.com/mono/internal/envconf"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Run(h http.Handler) {
	var h2 http2.Server
	s := http.Server{
		Addr:              ":" + envconf.String("PORT", "8080"),
		Handler:           h2c.NewHandler(h, &h2),
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
		log.Println("starting on", s.Addr)
		s.ListenAndServe()
	}()
	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
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
