package main

import (
	"go.seankhliao.com/mono/cmd/singlepage/internal/singlepage"
	"go.seankhliao.com/mono/internal/httpsvc"
)

func main() {
	httpsvc.Run(&singlepage.Server{})
}
