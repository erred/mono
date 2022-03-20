package main

import (
	"go.seankhliao.com/mono/cmd/blog/internal/blog"
	"go.seankhliao.com/mono/internal/httpsvc"
)

func main() {
	httpsvc.Run(&blog.Server{})
}
