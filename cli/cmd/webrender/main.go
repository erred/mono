package main

import (
	"flag"
	"os"

	"github.com/go-logr/stdr"
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

	log := stdr.New(nil)

	fi, err := os.Open(src)
	if err != nil {
		log.Error(err, "open src", "file", src)
		os.Exit(1)
	}
	defer fi.Close()
	fo, err := os.Create(dst)
	if err != nil {
		log.Error(err, "open dst", "file", dst)
		os.Exit(1)
	}
	defer fo.Close()

	err = render.Render(&o, fo, fi)
	if err != nil {
		log.Error(err, "render")
		os.Exit(1)
	}
}
