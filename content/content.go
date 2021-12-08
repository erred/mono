package content

import (
	"embed"
	"io/fs"
)

var (
	//go:embed medea/*
	fsMedea  embed.FS
	Medea, _ = fs.Sub(fsMedea, "medea")
	//go:embed paste/*
	fsPaste  embed.FS
	Paste, _ = fs.Sub(fsPaste, "paste")
	//go:embed stylesheet/*
	fsStylesheet  embed.FS
	Stylesheet, _ = fs.Sub(fsStylesheet, "stylesheet")
	//go:embed vanity/*
	fsVanity  embed.FS
	Vanity, _ = fs.Sub(fsVanity, "vanity")
	//go:embed w16/*
	fsW16  embed.FS
	W16, _ = fs.Sub(fsW16, "w16")
)
