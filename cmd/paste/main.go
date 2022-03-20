package main

import (
	"go.seankhliao.com/mono/cmd/paste/internal/paste"
	"go.seankhliao.com/mono/internal/httpsvc"
)

func main() {
	httpsvc.Run(&paste.Server{})
}
