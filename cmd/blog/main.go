package main

import (
	"go.seankhliao.com/mono/cmd/blog/internal/blog"
	"go.seankhliao.com/mono/internal/svc"
)

func main() {
	svc.Run(&blog.Server{}, nil)
}
