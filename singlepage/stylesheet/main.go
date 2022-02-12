package main

import (
	_ "embed"
	"log"

	"go.seankhliao.com/mono/internal/runhttp"
	"go.seankhliao.com/mono/internal/singlepage"
)

//go:embed index.md
var content []byte

func main() {
	s, err := singlepage.New("stylesheet", content)
	if err != nil {
		log.Fatalln(err)
	}

	runhttp.Run(s)
}
