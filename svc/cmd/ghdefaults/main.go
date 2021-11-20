package main

import (
	"flag"
	"os"

	"go.seankhliao.com/mono/svc/webserver"
	ghdefaults "go.seankhliao.com/mono/svc/cmd/ghdefaults/internal"
)

func main() {
	so := ghdefaults.NewOptions(flag.CommandLine)
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = ghdefaults.New(ctx, so)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}
