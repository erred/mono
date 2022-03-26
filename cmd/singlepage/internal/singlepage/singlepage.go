package singlepage

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/envconf"
	"go.seankhliao.com/mono/internal/httpsvc"
	"go.seankhliao.com/mono/internal/web/render"
)

var (
	//go:embed favicon.ico
	favicon []byte

	_ httpsvc.HTTPSvc = &Server{}
)

type Server struct {
	log zerolog.Logger

	ts       time.Time
	rendered []byte
}

func (s Server) Desc() string {
	return `a single static page`
}

func (s Server) Help() string {
	return `
SINGLEPAGE_SOURCE
        path to markdown file to render
`
}

func (s *Server) Init(log zerolog.Logger) error {
	s.log = log
	s.ts = time.Now()

	src := envconf.String("SINGLEPAGE_SOURCE", "index.md")
	raw, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read SINGLEPAGE_SOURCE=%q: %w", src, err)
	}
	s.rendered, err = render.CompactBytes("", "", "", raw)
	if err != nil {
		return fmt.Errorf("prerender page: %w", err)
	}

	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/favicon.ico":
		http.ServeContent(rw, r, "favicon.ico", s.ts, bytes.NewReader(favicon))
	case "/":
		http.ServeContent(rw, r, "index.html", s.ts, bytes.NewReader(s.rendered))
	default:
		http.Redirect(rw, r, "/", http.StatusFound)
	}
}
