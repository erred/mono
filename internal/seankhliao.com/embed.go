package seankhliaocom

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static
	staticFS    embed.FS
	StaticFS, _ = fs.Sub(staticFS, "static")

	//go:embed content
	contentFS    embed.FS
	ContentFS, _ = fs.Sub(contentFS, "content")
)
