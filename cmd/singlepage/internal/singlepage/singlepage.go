package singlepage

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/httpsvc"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/internal/webstatic"
)

var _ httpsvc.HTTPSvc = &Server{}

type Server struct {
	log zerolog.Logger

	mux *http.ServeMux
}

func (s *Server) Init(init *httpsvc.Init) error {
	s.log = init.Log
	s.mux = http.NewServeMux()
	webstatic.Register(s.mux)

	var src string
	init.Flags.StringVar(&src, "singlepage.source", "index.md", "source file in markdown")
	init.FlagsAfter = func() error {
		ts := time.Now()
		raw, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read src %s: %w", src, err)
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

	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}
