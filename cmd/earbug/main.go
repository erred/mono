package main

import (
	"go.seankhliao.com/mono/cmd/earbug/internal/earbug"
	"go.seankhliao.com/mono/internal/svc"
)

func main() {
	svc.Run(&earbug.Server{}, nil)
}
