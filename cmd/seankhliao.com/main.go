package main

import (
	"log"

	"go.seankhliao.com/mono/internal/runhttp"
	seankhliaocom "go.seankhliao.com/mono/internal/seankhliao.com"
)

func main() {
	s, err := seankhliaocom.New()
	if err != nil {
		log.Fatalln("setup seankhliao.com", err)
	}
	runhttp.Run(s)
}
