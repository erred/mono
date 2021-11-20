package main

import (
	"flag"
	"os"

	"go.seankhliao.com/mono/svc/webserver"
	"go.seankhliao.com/mono/svc/cmd/archrepod/internal/server"
)

func main() {
	so := server.NewOptions(flag.CommandLine)
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = server.New(ctx, so)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}
