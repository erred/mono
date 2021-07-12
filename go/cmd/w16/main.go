package main

import (
	"flag"
	"os"

	"go.seankhliao.com/mono/go/internal/w16"
	"go.seankhliao.com/mono/go/webserver"
)

func main() {
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = w16.New(ctx)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}
