package static

import (
	"embed"
	"io/fs"
)

var (
	//go:embed seankhliao.com/*
	s    embed.FS
	S, _ = fs.Sub(s, "seankhliao.com")
)
