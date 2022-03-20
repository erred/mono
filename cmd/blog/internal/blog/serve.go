package blog

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

// registerStatic registers handlers for all file paths in the fsys
func (s *Server) registerStatic(mux *http.ServeMux, fsys fs.FS) error {
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.Type() != 0 {
			return nil
		}

		s.handleFsysFile(mux, fsys, p)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk static: %w", err)
	}
	return nil
}

func (s *Server) handleBytes(mux *http.ServeMux, p, name string, b []byte) {
	t := time.Now()
	br := bytes.NewReader(b)
	mux.HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != p {
			s.notfound(rw, r)
			return
		}

		http.ServeContent(rw, r, name, t, br)
	})
}

// handleFsys serves a path directly from an embedded fs
func (s *Server) handleFsysFile(mux *http.ServeMux, fsys fs.FS, p string) {
	cp := canonicalPath(p)
	mux.HandleFunc(cp, func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != cp {
			s.notfound(rw, r)
			return
		}

		f, _ := fsys.Open(p)
		defer f.Close()
		http.ServeContent(rw, r, p, time.Time{}, f.(io.ReadSeeker))
	})
}

func (s *Server) notfound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	rw.Write(s.notFound)
}

func canonicalPath(p string) string {
	var end string
	if strings.HasSuffix(p, ".html") {
		p = p[:len(p)-5]
		end = "/"
	} else if strings.HasSuffix(p, ".md") {
		p = p[:len(p)-3]
		end = "/"
	}
	ps := strings.Split(p, "/")
	if ps[len(ps)-1] == "index" {
		ps = ps[:len(ps)-1]
	}
	p = path.Join(ps...) + end
	if p[0] != '/' {
		p = "/" + p
	}
	return p
}
