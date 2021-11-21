package main

import (
	"flag"
	"os"

	"github.com/go-logr/logr"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

func main() {
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	wo := NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = wo.Handler(ctx)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "setup w16")
		os.Exit(1)
	}

	ho.Run()
}
