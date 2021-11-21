package main

import (
	"flag"
	"os"

	"github.com/go-logr/logr"
	"go.seankhliao.com/mono/svc/cmd/archrepod/internal/server"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

func main() {
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	so := server.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = server.New(ctx, so)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "setup archrepod")
		os.Exit(1)
	}

	ho.Run()
}
