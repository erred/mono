package main

import (
	"log"

	"go.seankhliao.com/mono/internal/blog"
	"go.seankhliao.com/mono/internal/runhttp"
)

func main() {
	s, err := blog.New()
	if err != nil {
		log.Fatalln("setup vanity", err)
	}
	runhttp.Run(s)
}
