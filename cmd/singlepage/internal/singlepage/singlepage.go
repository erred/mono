package singlepage

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/svc"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/internal/webstatic"
)

var _ svc.SHTTP = &Server{}

type Server struct {
	log zerolog.Logger

	src string
	mux *http.ServeMux
}

func (s *Server) Register(r svc.Register) error {
	r.Flags.StringVar(&s.src, "singlepage.source", "index.md", "source file in markdown")
	return nil
}

func (s *Server) Init(init svc.Init) error {
	s.log = init.Logger
	s.mux = http.NewServeMux()
	webstatic.Register(s.mux)

	ts := time.Now()
	raw, err := os.ReadFile(s.src)
	if err != nil {
		return fmt.Errorf("read src %s: %w", s.src, err)
	}
	rendered, err := render.CompactBytes("", "", "", raw)
	if err != nil {
		return fmt.Errorf("prerender page: %w", err)
	}
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		http.ServeContent(w, r, "index.html", ts, bytes.NewReader(rendered))
	})

	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}
