package main

import (
	"go.seankhliao.com/mono/cmd/vanity/internal/vanity"
	"go.seankhliao.com/mono/internal/httpsvc"
)

func main() {
	httpsvc.Run(&vanity.Server{})
}
