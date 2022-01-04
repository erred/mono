package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/content"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/runsvr"
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
	r := runsvr.New(flag.CommandLine)
	runsvr.Desc(flag.CommandLine, docgo)
	s := New(flag.CommandLine)
	flag.Parse()

	r.HTTP(s)
}

type Server struct {
	l logr.Logger
	t trace.Tracer

	repoCtr metric.Int64Counter

	indexBytes []byte
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	return &s
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l
	s.t = t.Tracer("vanity")

	var err error
	s.repoCtr, err = m.Meter("vanity").NewInt64Counter("vanity")
	if err != nil {
		return fmt.Errorf("register metric: %w", err)
	}

	indexRaw, err := fs.ReadFile(content.Vanity, "index.md")
	if err != nil {
		return fmt.Errorf("read index.md: %w", err)
	}

	s.indexBytes, err = render.CompactBytes(
		"go.seankhliao.com",
		"Go custom import path server",
		"https://go.seankhliao.com",
		indexRaw,
	)
	if err != nil {
		return fmt.Errorf("render index.md: %w", err)
	}

	mux.Handle("/", s)

	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, span := s.t.Start(r.Context(), "vanity")
	defer span.End()

	if r.URL.Path == "/" {
		http.ServeContent(rw, r, "index.html", time.Time{}, bytes.NewReader(s.indexBytes))
		return
	}

	repo := strings.Split(r.URL.Path, "/")[1]
	span.SetAttributes(attribute.String("repo", repo))
	err := render.Render(&render.Options{
		Data: render.PageData{
			URLCanonical: "https://go.seankhliao.com/" + repo,
			Compact:      true,
			Title:        "go.seankhliao.com/" + repo,
			Head:         fmt.Sprintf(repoPageHeader, repo),
		},
	}, rw, strings.NewReader(fmt.Sprintf(repoPageBody, repo)))
	if err != nil {
		s.l.Error(err, "render repo page")
	}
}
