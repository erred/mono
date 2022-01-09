package main

import (
	"flag"

	"go.seankhliao.com/mono/svc/runsvr"
)

func main() {
	s := New(flag.CommandLine)
	r := runsvr.New(flag.CommandLine)
	runsvr.Desc(flag.CommandLine, docgo)
	flag.Parse()

	r.GRPC(s)
}
