package main

import (
	"flag"
	"os"

	"github.com/go-logr/logr"
	ghdefaults "go.seankhliao.com/mono/svc/cmd/ghdefaults/internal"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

func main() {
	so := ghdefaults.NewOptions(flag.CommandLine)
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = ghdefaults.New(ctx, so)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "setup ghdefaults")
		os.Exit(1)
	}

	ho.Run()
}
