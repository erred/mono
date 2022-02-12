package run

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Runner interface {
	Run(context.Context) error
}

func Run(runners ...Runner) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	for _, runner := range runners {
		go func(runner Runner) {
			defer wg.Done()
			defer cancel()
			runner.Run(ctx)
		}(runner)
	}

	select {
	case <-sigc:
		// external interrupt
	case <-ctx.Done():
		// one of the services exited
	}

	cancel()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-sigc:
		// interrupted shutdown
	case <-done:
		// gracefully shutdown
	}
}
