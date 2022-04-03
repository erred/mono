package vanity

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/svc"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/internal/webstatic"
)

var (
	//go:embed index.md
	indexRaw []byte

	//go:embed repo.md.tpl
	repoRaw []byte
	repoTpl = template.Must(template.New("").Parse(string(repoRaw)))

	_ svc.SHTTP = &Server{}
)

type Server struct {
	log  zerolog.Logger
	host string

	ts    time.Time
	mux   *http.ServeMux
	index []byte
}

func (s *Server) Register(r svc.Register) error {
	return nil
}

func (s *Server) Init(init svc.Init) error {
	s.log = init.Logger
	s.ts = time.Now()

	s.mux = http.NewServeMux()
	webstatic.Register(s.mux)
	s.mux.HandleFunc("/", s.handler)

	var err error
	s.index, err = render.CompactBytes("", "", "", indexRaw)
	if err != nil {
		return fmt.Errorf("prerender index: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(rw, r)
}

func (s *Server) handler(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeContent(rw, r, "index.html", s.ts, bytes.NewReader(s.index))
		return
	}

	repo := strings.Split(r.URL.Path, "/")[1]
	var buf bytes.Buffer
	err := repoTpl.Execute(&buf, map[string]string{"Repo": repo})
	if err != nil {
		s.log.Err(err).Msg("render template")
		http.Error(rw, "render err", http.StatusInternalServerError)
		return
	}
	err = render.Render(&render.Options{
		Data: render.PageData{
			URLCanonical: fmt.Sprintf("https://%s/%s", s.host, repo),
			Compact:      true,
			Title:        "go.seankhliao.com/" + repo,
		},
	}, rw, bytes.NewReader(buf.Bytes()))
	if err != nil {
		s.log.Err(err).Msg("render page")
		http.Error(rw, "render err", http.StatusInternalServerError)
		return
	}
}
