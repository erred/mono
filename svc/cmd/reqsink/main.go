package main

import (
	"flag"

	"go.seankhliao.com/mono/svc/runsvr"
)

func main() {
	r := runsvr.New(flag.CommandLine)
	s := New(flag.CommandLine)
	flag.Parse()

	r.HTTP(s)
}
