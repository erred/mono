package internal

import (
	"embed"
	"io/fs"
)

var (
	//go:embed all:static
	staticFS    embed.FS
	StaticFS, _ = fs.Sub(staticFS, "static")

	//go:embed all:content
	contentFS    embed.FS
	ContentFS, _ = fs.Sub(contentFS, "content")
)
