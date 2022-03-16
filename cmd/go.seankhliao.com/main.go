package main

import (
	"log"

	goseankhliaocom "go.seankhliao.com/mono/internal/go.seankhliao.com"
	"go.seankhliao.com/mono/internal/runhttp"
)

func main() {
	s, err := goseankhliaocom.New()
	if err != nil {
		log.Fatalln("setup go.seankhliao.com", err)
	}
	runhttp.Run(s)
}
