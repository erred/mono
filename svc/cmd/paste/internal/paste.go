package paste

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base32"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/internal/web/render"
)

var (
	//go:embed body.html
	body string
	//go:embed style.css
	style string
)

type Options struct {
	Dir string
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.Dir, "data", "/data", "path to data directory")
	return &o
}

type server struct {
	dir string
	t   trace.Tracer
}

func New(ctx context.Context, o *Options) (http.Handler, error) {
	tracer := otel.Tracer("paste")
	ctx, span := tracer.Start(ctx, "new")
	defer span.End()

	_, span = tracer.Start(ctx, "ensure-directory")
	err := os.MkdirAll(o.Dir, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		span.End()
		return nil, fmt.Errorf("ensure directory %s: %w", o.Dir, err)
	}
	span.End()

	_, span = tracer.Start(ctx, "render-index")
	defer span.End()
	ro := &render.Options{
		MarkdownSkip: true,
		Data: render.PageData{
			Compact:      true,
			URLCanonical: "https://p.seankhliao.com",
			Title:        `paste`,
			Description:  `simple paste service`,
			H1:           `paste`,
			H2:           `upload`,
			Style:        style,
		},
	}

	b := bytes.NewReader([]byte(body))

	fn := filepath.Join(o.Dir, "index.html")
	fout, err := os.Create(fn)
	if err != nil {
		return nil, fmt.Errorf("create index file %s: %w", fn, err)
	}

	err = render.Render(ro, fout, b)
	if err != nil {
		return nil, fmt.Errorf("render index file: %w", err)
	}

	return &server{
		dir: o.Dir,
		t:   tracer,
	}, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "handle")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx).WithName("dispatch")
	l = l.WithValues("method", r.Method, "path", r.URL.Path, "user-agent", r.UserAgent())
	defer l.Info("access")

	r = r.WithContext(ctx)
	if r.Method == http.MethodPost && r.URL.Path == "/api/v0/form" {
		s.formHandler(w, r)
		return
	}
	http.FileServer(http.Dir(s.dir)).ServeHTTP(w, r)
}

func (s *server) formHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "form-handler")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx).WithName("formhandler")

	ctx, span = s.t.Start(ctx, "process-form")
	err := r.ParseMultipartForm(1 << 25) // 32M
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		span.End()
		return
	}
	mpf, mph, err := r.FormFile("upload")
	span.End()
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		l.Error(err, "process form file")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} else if err == nil {
		defer mpf.Close()

		ctx, span := s.t.Start(ctx, "upload-file")
		defer span.End()
		l = l.WithValues("type", "upload")

		_, span = s.t.Start(ctx, "uniq-dir")
		ext := filepath.Ext(mph.Filename)
		f, err := uniqDir(s.dir, ext)
		l = l.WithValues("file", f.Name())
		span.End()
		if err != nil {
			l.Error(err, "find uniq dir")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		_, span = s.t.Start(ctx, "write-file")
		_, err = io.Copy(f, mpf)
		span.End()
		if err != nil {
			l.Error(err, "write file")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/"+filepath.Base(f.Name()), http.StatusFound)
		l.Info("handled upload")
		return
	}

	_, span = s.t.Start(ctx, "paste-content")
	defer span.End()

	l = l.WithValues("type", "paste")
	p := r.FormValue("paste")
	if p == "" {
		l.Error(errors.New("no content"), "get paste")
		http.Error(w, "one of upload or paste is required", http.StatusBadRequest)
		return
	}

	_, span = s.t.Start(ctx, "uniq-dir")
	f, err := uniqDir(s.dir, ".txt")
	span.End()
	if err != nil {
		l.Error(err, "find uniq dir")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, span = s.t.Start(ctx, "write-file")
	_, err = f.WriteString(p)
	span.End()
	if err != nil {
		l.Error(err, "write file")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+filepath.Base(f.Name()), http.StatusSeeOther)
	l.Info("handled")
}

func uniqDir(dir, ext string) (*os.File, error) {
	totalRetries := 3
	for retries := totalRetries; retries > 0; retries-- {
		b := make([]byte, 5)
		rand.Read(b)
		fn := base32.StdEncoding.EncodeToString(b) + ext
		f, err := os.OpenFile(filepath.Join(dir, fn), os.O_RDWR|os.O_CREATE, 0o644)
		switch {
		case err == nil:
			return f, nil
		case errors.Is(err, os.ErrExist):
			continue
		default:
			return nil, fmt.Errorf("create file=%s: %w", fn, err)
		}
	}
	return nil, fmt.Errorf("exceeded retries=%d", totalRetries)
}
