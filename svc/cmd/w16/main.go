package main

import (
	"flag"
	"os"

	server "go.seankhliao.com/mono/svc/cmd/w16/internal"
	"go.seankhliao.com/mono/svc/webserver"
)

func main() {
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = server.New(ctx)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}
