package main

import (
	"flag"
	"os"

	"go.seankhliao.com/mono/internal/feedagg"
	"go.seankhliao.com/mono/internal/webserver"
)

func main() {
	so := feedagg.NewOptions(flag.CommandLine)
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = feedagg.New(ctx, so)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}
