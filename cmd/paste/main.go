package main

import (
	"go.seankhliao.com/mono/cmd/paste/internal/paste"
	"go.seankhliao.com/mono/internal/svc"
)

func main() {
	svc.Run(&paste.Server{}, nil)
}
