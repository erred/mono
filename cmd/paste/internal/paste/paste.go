package paste

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/envconf"
	"go.seankhliao.com/mono/internal/httpsvc"
	"go.seankhliao.com/mono/internal/web/render"
)

var (
	_ httpsvc.HTTPSvc = &Server{}

	//go:embed index.md
	indexRaw []byte

	//go:embed paste.md
	pasteRaw []byte

	//go:embed favicon.ico
	favicon []byte
)

type Server struct {
	log zerolog.Logger

	dir string

	ts    time.Time
	index []byte
	paste []byte
}

func (s Server) Desc() string {
	return `pastebin`
}

func (s Server) Help() string {
	return `
PASTE_DIR
        path to storage directory
`
}

func (s *Server) Init(log zerolog.Logger) error {
	s.log = log
	s.dir = envconf.String("PASTE_DIR", "/var/lib/paste")

	s.ts = time.Now()
	var err error
	s.index, err = render.CompactBytes("", "", "", indexRaw)
	if err != nil {
		return fmt.Errorf("render index: %w", err)
	}
	s.paste, err = render.CompactBytes("", "", "", pasteRaw)
	if err != nil {
		return fmt.Errorf("render paste: %w", err)
	}
	return nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.ServeContent(rw, r, "x.html", s.ts, bytes.NewReader(s.index))
	case "/favicon.ico":
		http.ServeContent(rw, r, "favicon.ico", s.ts, bytes.NewReader(favicon))
	case "/paste":
		http.Redirect(rw, r, "/paste/", http.StatusTemporaryRedirect)
	case "/paste/":
		switch r.Method {
		case http.MethodGet:
			http.ServeContent(rw, r, "x.html", s.ts, bytes.NewReader(s.paste))
		case http.MethodPost:
			s.upload(rw, r)
		default:
			http.Error(rw, "GET / POST only", http.StatusMethodNotAllowed)
		}
	default:
		p := r.URL.Path
		if strings.HasPrefix(p, "/p/") && !strings.Contains(p[3:], "/") {
			s.lookup(rw, r)
		} else {
			http.Error(rw, "not found", http.StatusNotFound)
		}
	}
}

func (s *Server) lookup(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, filepath.Join(s.dir, r.URL.Path))
}

func (s *Server) upload(rw http.ResponseWriter, r *http.Request) {
	val := []byte(r.FormValue("paste"))
	if len(val) == 0 {
		err := r.ParseMultipartForm(1 << 22) // 4M
		if err != nil {
			s.log.Err(err).Msg("parse multipart form")
			http.Error(rw, "bad multipart form", http.StatusBadRequest)
			return
		}
		mpf, _, err := r.FormFile("upload")
		if err != nil {
			s.log.Err(err).Msg("get form file")
			http.Error(rw, "bad multipart form", http.StatusBadRequest)
			return
		}
		defer mpf.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, mpf)
		if err != nil {
			s.log.Err(err).Msg("read form file")
			http.Error(rw, "read", http.StatusInternalServerError)
			return
		}
		val = buf.Bytes()
	}

	sum := sha256.Sum256(val)
	sum2 := base64.URLEncoding.EncodeToString(sum[:])

	key := path.Join("p", sum2[:8])
	err := os.WriteFile(filepath.Join(s.dir, key), val, 0o644)
	if err != nil {
		s.log.Err(err).Str("file", key).Msg("write file")
		http.Error(rw, "write", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "https://%s/%s\n", r.Host, key)
}
