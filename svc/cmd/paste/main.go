package main

import (
	"flag"
	"os"

	"go.seankhliao.com/mono/internal/paste"
	"go.seankhliao.com/mono/internal/webserver"
)

func main() {
	wo := webserver.NewOptions(flag.CommandLine)
	so := paste.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = paste.New(ctx, so)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}