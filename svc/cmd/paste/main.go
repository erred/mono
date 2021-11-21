package main

import (
	"flag"
	"os"

	"github.com/go-logr/logr"
	paste "go.seankhliao.com/mono/svc/cmd/paste/internal"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

func main() {
	so := paste.NewOptions(flag.CommandLine)
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = paste.New(ctx, so)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "setup paste")
		os.Exit(1)
	}

	ho.Run()
}
