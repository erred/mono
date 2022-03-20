package main

import (
	"flag"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/web/render"
)

func main() {
	var o render.Options
	var src, dst string
	flag.BoolVar(&o.Data.Compact, "compact", true, "compact header")
	flag.StringVar(&src, "src", "index.md", "source file")
	flag.StringVar(&dst, "dst", "index.html", "destination file")
	flag.StringVar(
		&o.Data.URLCanonical,
		"url",
		"https://seankhliao.com",
		"base url for canonicalization",
	)
	flag.Parse()

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	fi, err := os.Open(src)
	if err != nil {
		log.Err(err).Str("src", src).Msg("open file")
		os.Exit(1)
	}
	defer fi.Close()
	fo, err := os.Create(dst)
	if err != nil {
		log.Err(err).Str("dst", dst).Msg("create file")
		os.Exit(1)
	}
	defer fo.Close()

	err = render.Render(&o, fo, fi)
	if err != nil {
		log.Err(err).Msg("render")
		os.Exit(1)
	}
}
