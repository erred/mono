package monolith

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.seankhliao.com/mono/monolith/component"
	"go.seankhliao.com/mono/monolith/o11y"
	"go.seankhliao.com/mono/monolith/run"
	"go.seankhliao.com/mono/monolith/run/rungrpc"
	"go.seankhliao.com/mono/monolith/run/runhttp"
)

func Run(components ...component.Component) {
	// servers
	oe := o11y.New()
	hs := runhttp.New()
	gs := rungrpc.New()
	components = append(components, oe, hs, gs)

	// flags
	flags := flag.CommandLine
	for _, c := range components {
		c.Register(flags)
	}
	flags.Parse(os.Args[1:])

	op := oe.Provider()

	runners := []run.Runner{oe}
	func() {
		// setup context
		ctx := context.Background()
		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		for _, c := range components {
			if !c.Enabled() {
				continue
			}
			if h, ok := c.(component.HTTP); ok {
				h.HTTP(ctx, op, hs.M)
			}
			if g, ok := c.(component.GRPC); ok {
				g.GRPC(ctx, op, gs.S)
			}
			if b, ok := c.(component.Background); ok {
				b.Background(ctx, op)
				runners = append(runners, b)
			}
		}
	}()

	run.Run(runners...)
}
