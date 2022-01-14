package main

import (
	"flag"

	"go.seankhliao.com/mono/auth/authext"
	"go.seankhliao.com/mono/svc/runsvr"
)

func main() {
	r := runsvr.New(flag.CommandLine)
	runsvr.Desc(flag.CommandLine, docgo)
	s := authext.New(flag.CommandLine)
	flag.Parse()

	r.GRPC(s)
}
