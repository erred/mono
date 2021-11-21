package main

import (
	"flag"
	"os"

	"github.com/go-logr/logr"
	feedagg "go.seankhliao.com/mono/svc/cmd/feedagg/internal"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

func main() {
	so := feedagg.NewOptions(flag.CommandLine)
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = feedagg.New(ctx, so)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "setup feedagg")
		os.Exit(1)
	}

	ho.Run()
}
