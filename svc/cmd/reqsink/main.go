package main

import (
	"flag"

	"go.seankhliao.com/mono/svc/runsvr"
)

func main() {
	r := runsvr.New(flag.CommandLine)
	runsvr.Desc(flag.CommandLine, docgo)
	s := New(flag.CommandLine)
	flag.Parse()

	r.HTTP(s)
}
