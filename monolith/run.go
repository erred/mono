package monolith

import (
	"context"
	"net"
	"net/http"
	"sync"
)

type runner interface {
	run(ctx context.Context) error
}

func run(ctx context.Context, runners ...runner) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var runwg sync.WaitGroup
	for _, r := range runners {
		runwg.Add(1)
		go func(r runner) {
			defer runwg.Done()
			defer cancel()
			_ = r.run(ctx)
		}(r)
	}

	<-ctx.Done()
	cancel()
	runwg.Wait()
}

type runHttp struct {
	s *http.Server
}

func (r *runHttp) run(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	defer lis.Close()
	go func() {
		<-ctx.Done()
		r.s.Shutdown(ctx)
	}()

	err = r.s.Serve(lis)
	if err != nil {
		// TODO: filter
		return err
	}
	return nil
}
