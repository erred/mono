package main

import (
	"log"

	"go.seankhliao.com/mono/internal/runhttp"
	"go.seankhliao.com/mono/internal/vanity"
)

func main() {
	s, err := vanity.New()
	if err != nil {
		log.Fatalln("setup vanity", err)
	}
	runhttp.Run(s)
}
