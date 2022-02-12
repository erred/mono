package main

import (
	"log"

	"go.seankhliao.com/mono/internal/earbug"
	"go.seankhliao.com/mono/internal/runhttp"
)

func main() {
	s, err := earbug.New()
	if err != nil {
		log.Fatalln("setup earbug", err)
	}
	runhttp.Run(s)
}
