package main

import (
	"go.seankhliao.com/mono/cmd/singlepage/internal/singlepage"
	"go.seankhliao.com/mono/internal/svc"
)

func main() {
	svc.Run(&singlepage.Server{}, nil)
}
