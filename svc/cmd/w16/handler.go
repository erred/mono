package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.seankhliao.com/mono/content"
	"go.seankhliao.com/mono/static"
)

type Options struct {
	Hostname string
	GTMID    string
	Compact  bool

	notFoundBody []byte
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.Hostname, "hostname", "seankhliao.com", "hostname for generating canonical paths")
	fs.StringVar(&o.GTMID, "gtm", "", "google tag manager id")
	fs.BoolVar(&o.Compact, "compact", false, "render with compact header")
	return &o
}

func (o *Options) Handler(ctx context.Context) (http.Handler, error) {
	staticFsys, err := fs.Sub(static.Static, o.Hostname)
	if err != nil {
		return nil, fmt.Errorf("generating static sub fs for %s: %w", o.Hostname, err)
	}

	var contentFsys fs.FS
	switch o.Hostname {
	case "go.seankhliao.com":
		contentFsys = content.Vanity
	case "medea.seankhliao.com":
		contentFsys = content.Medea
	case "seankhliao.com":
		contentFsys = content.W16
	case "stylesheet.seankhliao.com":
		contentFsys = content.Stylesheet
	default:
		return nil, fmt.Errorf("no matching embedded fs: %s", o.Hostname)
	}

	mux := http.NewServeMux()
	o.registerStatic(mux, staticFsys)
	o.renderAndRegister(mux, contentFsys)

	return defaultHandler(mux), nil
}

func defaultHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := logr.FromContextOrDiscard(r.Context())

		l.Info("received request",
			"url", r.URL.String(),
			"referer", r.Referer(),
			"user_agent", r.UserAgent(),
		)

		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodGet:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Max-Age", "86400")
			h.ServeHTTP(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// registerStatic registers handlers for all file paths in the fsys
func (o *Options) registerStatic(mux *http.ServeMux, fsys fs.FS) error {
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.Type() != 0 {
			return nil
		}

		o.handleFsysFile(mux, fsys, p)
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk static: %w", err)
	}
	return nil
}

// handleFsys serves a path directly from an embedded fs
func (o *Options) handleFsysFile(mux *http.ServeMux, fsys fs.FS, p string) {
	t := time.Now()
	cp := canonicalPath(p)
	mux.HandleFunc(cp, func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != cp {
			o.notfound(rw, r)
			return
		}

		f, _ := fsys.Open(p)
		defer f.Close()
		http.ServeContent(rw, r, p, t, f.(io.ReadSeeker))
	})
}

// handleBytes serves a path directly from bytes
func (o *Options) handleBytes(mux *http.ServeMux, p, name string, b []byte) {
	t := time.Now()
	br := bytes.NewReader(b)
	mux.HandleFunc(p, func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path != p {
			o.notfound(rw, r)
			return
		}

		http.ServeContent(rw, r, name, t, br)
	})
}

// notfound is a simple not found handler
func (o *Options) notfound(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusNotFound)
	http.ServeContent(rw, r, "404.html", time.Time{}, bytes.NewReader(o.notFoundBody))
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
