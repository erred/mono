package main

import (
	"go.seankhliao.com/mono/cmd/vanity/internal/vanity"
	"go.seankhliao.com/mono/internal/svc"
)

func main() {
	svc.Run(&vanity.Server{}, nil)
}
