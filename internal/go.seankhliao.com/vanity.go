package goseankhliaocom

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"go.seankhliao.com/mono/internal/envconf"
	"go.seankhliao.com/mono/internal/web/render"
)

var (
	//go:embed index.md
	indexRaw []byte

	//go:embed repo.md.tpl
	repoRaw []byte
	repoTpl = template.Must(template.New("").Parse(string(repoRaw)))
)

type Server struct {
	host  string
	index []byte
}

func New() (*Server, error) {
	s := &Server{
		host: envconf.String("VANITY_HOST", "go.seankhliao.com"),
	}

	var err error
	s.index, err = render.CompactBytes(
		s.host,
		"custom go import paths",
		"https://"+s.host+"/",
		indexRaw,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeContent(rw, r, "index.html", time.Time{}, bytes.NewReader(s.index))
		return
	}

	repo := strings.Split(r.URL.Path, "/")[1]
	var buf bytes.Buffer
	err := repoTpl.Execute(&buf, map[string]string{"Repo": repo})
	if err != nil {
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
		http.Error(rw, "render err", http.StatusInternalServerError)
		return
	}
}
