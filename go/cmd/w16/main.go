package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.seankhliao.com/mono/go/internal/w16"
	"go.seankhliao.com/mono/go/webserver"
	"k8s.io/klog/v2/klogr"
)

func main() {
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()

	m, err := w16.New(ctx)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	wo.Logger = l
	wo.Handler = m

	webserver.New(ctx, wo).Run(ctx)
}
