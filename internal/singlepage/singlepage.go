package singlepage

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"go.seankhliao.com/mono/internal/envconf"
	"go.seankhliao.com/mono/internal/web/render"
)

type Server struct {
	host     string
	rendered []byte
}

func New(name string, content []byte) (*Server, error) {
	s := &Server{
		host: envconf.String("SINGLEPAGE_HOST", ""),
	}

	var err error
	s.rendered, err = render.CompactBytes(
		name,
		"welcome to "+name,
		fmt.Sprintf("https://%s/", s.host),
		content,
	)
	if err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}

	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(rw, r, "/", http.StatusFound)
		return
	}
	http.ServeContent(rw, r, "x.html", time.Time{}, bytes.NewReader(s.rendered))
}
