package main

import (
	"go.seankhliao.com/mono/cmd/ghdefaults/internal/ghdefaults"
	"go.seankhliao.com/mono/internal/svc"
)

func main() {
	svc.Run(&ghdefaults.Server{}, nil)
}
