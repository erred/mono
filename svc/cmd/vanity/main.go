package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.seankhliao.com/mono/content"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

const (
	repoPageHeader = `
<meta
  name="go-import"
  content="go.seankhliao.com/%[1]s git https://github.com/seankhliao/%[1]s">
<meta
  name="go-source"
  content="go.seankhliao.com/%[1]s
    https://github.com/seankhliao/%[1]s
    https://github.com/seankhliao/%[1]s/tree/master{/dir}
    https://github.com/seankhliao/%[1]s/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="5;url=https://pkg.go.dev/go.seankhliao.com/%[1]s" />`
	repoPageBody = `---
title: %[1]s
description: module go.seankhliao.com/%[1]s
---

### _go.seankhliao.com_ / %[1]s

_source:_ [github](https://github.com/seankhliao.com/%[1]s)

_docs:_ [pkg.go.dev](https://pkg.go.dev/go.seankhliao.com/%[1]s)
`
)

func main() {
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	ho.Handler = New(ctx)

	ho.Run()
}

func New(ctx context.Context) http.Handler {
	l := logr.FromContextOrDiscard(ctx).WithName("vanity")

	indexRaw, err := content.Content.ReadFile("go.seankhliao.com/index.md")
	if err != nil {
		l.Error(err, "read index")
	}

	indexBytes, err := render.CompactBytes(
		"go.seankhliao.com",
		"Go custom import path server",
		"https://go.seankhliao.com",
		indexRaw,
	)
	if err != nil {
		l.Error(err, "render index")
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeContent(rw, r, "index.html", time.Time{}, bytes.NewReader(indexBytes))
			return
		}

		repo := strings.Split(r.URL.Path, "/")[1]
		err := render.Render(&render.Options{
			Data: render.PageData{
				URLCanonical: "https://go.seankhliao.com/" + repo,
				Compact:      true,
				Title:        "go.seankhliao.com/" + repo,
				Head:         fmt.Sprintf(repoPageHeader, repo),
			},
		}, rw, strings.NewReader(fmt.Sprintf(repoPageBody, repo)))
		if err != nil {
			l.Error(err, "render repo page")
		}
	})
}
