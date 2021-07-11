package w16

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/go/static"
)

type server struct {
	mux *http.ServeMux
	fs  fs.FS
	t   trace.Tracer
}

func New(ctx context.Context) (http.Handler, error) {
	s := &server{
		mux: http.NewServeMux(),
		fs:  static.S,
		t:   otel.Tracer("w16"),
	}

	ctx, span := s.t.Start(ctx, "walk-dir")
	defer span.End()
	err := fs.WalkDir(s.fs, ".", func(op string, d fs.DirEntry, err error) error {
		_, span := s.t.Start(ctx, "walk-content", trace.WithAttributes(
			attribute.String("path", op),
			attribute.Bool("dir", d.IsDir()),
		))
		defer span.End()

		if d.IsDir() {
			return nil
		}
		// canonical path
		p := canonicalPath(op)

		// l.Info("registering", "path", p)
		s.mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, span := s.t.Start(ctx, "serve")
			defer span.End()
			l := logr.FromContextOrDiscard(ctx).WithName("servefile")
			l = l.WithValues("method", r.Method, "path", r.URL.Path, "user-agent", r.UserAgent())

			// handle unknown paths
			if r.URL.Path != p {
				l.Error(errors.New("path mismatch"), "not found", "expected", p, "got", r.URL.Path)
				s.notFoundHandler(w, r)
				return
			}

			setHeaders(w, op)
			s.serveFile(ctx, w, op)
			l.Info("served")
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk embedded fs: %w", err)
	}

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	corsAllowAll(s.mux).ServeHTTP(w, r)
}

func (s *server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "not-found")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx).WithName("notfound")
	l = l.WithValues("path", r.URL.Path)

	if regexp.MustCompile(`/blog/\d{4}-\d{2}-\d{2}-.*/`).MatchString(r.URL.Path) {
		l.Info("redirected")
		http.Redirect(w, r, "/blog/1"+r.URL.Path[6:], http.StatusMovedPermanently)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	s.serveFile(ctx, w, "404.html")
	l.Info("not found")
}

func (s *server) serveFile(ctx context.Context, w http.ResponseWriter, p string) {
	ctx, span := s.t.Start(ctx, "serve-file")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx)

	file, err := s.fs.Open(p)
	if err != nil {
		l.Error(err, "open", "file", p)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		l.Error(err, "copy", "file", p)
	}
}

func canonicalPath(p string) string {
	if strings.HasSuffix(p, "index.html") {
		p = p[:len(p)-10]
	} else if strings.HasSuffix(p, ".html") {
		p = p[:len(p)-5] + "/"
	}
	return "/" + p
}

func setHeaders(w http.ResponseWriter, op string) {
	var ct, cc string
	switch path.Ext(op) {
	case ".css":
		ct = "text/css"
		cc = "max-age=2592000"
	case ".js":
		ct = "application/javascript"
		cc = "max-age=2592000"
	case ".jpg", ".jpeg":
		ct = "image/jpeg"
		cc = "max-age=2592000"
	case ".png":
		ct = "image/png"
		cc = "max-age=2592000"
	case ".webp":
		ct = "image/webp"
		cc = "max-age=2592000"
	case ".svg":
		ct = "image/svg+xml"
		cc = "max-age=2592000"
	case ".json":
		ct = "application/json"
	case ".otf", ".ttf", ".woff", ".woff2":
		ct = "font/" + path.Ext(op)
	case ".html":
		ct = "text/html"
	}

	if ct != "" {
		w.Header().Set("content-type", ct)
	}
	if cc != "" {
		w.Header().Set("cache-control", cc)
	}
}

func corsAllowAll(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("svr", "go.seankhliao.com/w/v16")
		w.Header().Set("easter-egg", "🐇*(🍆-🪴)=🐇🥚")

		switch r.Method {
		case http.MethodOptions:
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return
		case http.MethodGet, http.MethodPost:
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
