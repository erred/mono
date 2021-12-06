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
	paste := New(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	var err error
	ho.Handler, err = paste.Handler()
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "setup paste")
		os.Exit(1)
	}

	ho.Run()
}
