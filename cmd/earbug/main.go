package main

import (
	"go.seankhliao.com/mono/cmd/earbug/internal/earbug"
	"go.seankhliao.com/mono/internal/httpsvc"
)

func main() {
	httpsvc.Run(&earbug.Server{})
}
