package server

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const envBearerToken = "AR_TOKENS" // comma separated list of allowed bearer tokens

type Options struct {
	Dir string
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.Dir, "data", "/data", "path to store data")
	return &o
}

type server struct {
	root          string
	allowedTokens map[string]struct{}
	t             trace.Tracer
}

func New(ctx context.Context, opt *Options) (http.Handler, error) {
	s := &server{
		root:          opt.Dir,
		allowedTokens: make(map[string]struct{}),
		t:             otel.Tracer("archrepo/server"),
	}
	for _, t := range strings.Split(os.Getenv(envBearerToken), ",") {
		s.allowedTokens[t] = struct{}{}
	}
	return s, nil
}

type repositoryKey struct {
	repo string
	arch string
}

// /repos/$repo/$arch/$pkg
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "dispatch")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx).WithName("dispatch")
	l = l.WithValues("method", r.Method, "path", r.URL.Path, "user-agent", r.UserAgent())

	fields := strings.Split(r.URL.Path, "/")
	switch {
	case len(fields) >= 2 && fields[1] == "repos":
		switch r.Method {
		case http.MethodHead, http.MethodGet:
			http.FileServer(http.FS(os.DirFS(s.root))).ServeHTTP(w, r)
		case http.MethodPost:
			s.auth(s.addPackageHandler).ServeHTTP(w, r)
		case http.MethodDelete:
			s.auth(s.deleteHandler).ServeHTTP(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	defer l.Info("handled")
}

func (s *server) auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_, span := s.t.Start(ctx, "auth")
		l := logr.FromContextOrDiscard(ctx).WithName("auth")
		ah := r.Header.Get("authorization")
		if _, ok := s.allowedTokens[strings.TrimPrefix(ah, "Bearer ")]; !strings.HasPrefix(ah, "Bearer ") || !ok {
			l.Info("no auth", "method", r.Method, "path", r.URL.Path, "user-agent", r.UserAgent())
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			span.End()
			return
		}
		span.End()
		next.ServeHTTP(w, r)
	})
}

func (s *server) addPackageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := logr.FromContextOrDiscard(ctx).WithName("add")
	fields := strings.Split(r.URL.Path, "/")
	if len(fields) != 5 {
		l.Info("bad request", "path", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	l = l.WithValues("repo", fields[2], "pkg", fields[4])
	err := s.addPackage(ctx, fields[2], fields[3], fields[4], r.Body)
	if err != nil {
		l.Error(err, "add package")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	l.Info("package added")
	w.WriteHeader(http.StatusOK)
}

func (s *server) addPackage(ctx context.Context, repo, arch, pkg string, rc io.ReadCloser) error {
	ctx, span := s.t.Start(ctx, "add-package", trace.WithAttributes(
		attribute.String("repo", repo),
		attribute.String("arch", arch),
		attribute.String("package", pkg),
	))
	defer span.End()
	defer rc.Close()

	_, span = s.t.Start(ctx, "ensure-directory")
	p := filepath.Join(s.root, "repos", repo, arch)
	os.MkdirAll(p, 0o755) // ignore errors
	span.End()

	_, span = s.t.Start(ctx, "write-pkg")
	p = filepath.Join(p, pkg)
	f, err := os.Create(p)
	if err != nil {
		span.End()
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	if err != nil {
		span.End()
		return err
	}
	span.End()

	ctx, span = s.t.Start(ctx, "repo-add")
	defer span.End()
	dbp := filepath.Join(filepath.Dir(p), repo+".db.tar.zst")
	cmd := exec.CommandContext(ctx, "repo-add", dbp, p)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", cmd, err, string(out))
	}
	return nil
}

func (s *server) deleteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := logr.FromContextOrDiscard(ctx).WithName("delete")
	l = l.WithValues("path", r.URL.Path)
	fields := strings.Split(r.URL.Path, "/")
	var err error
	switch len(fields) {
	case 3:
		err = s.deleteRepo(ctx, fields[2])
	case 4:
		err = s.deleteArch(ctx, fields[2], fields[3])
	case 5:
		err = s.deletePackage(ctx, fields[2], fields[3], fields[4])
	default:
		l.Info("bad request")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err != nil {
		l.Error(err, "delete")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	l.Info("deleted")
	w.WriteHeader(http.StatusOK)
}

func (s *server) deleteRepo(ctx context.Context, repo string) error {
	_, span := s.t.Start(ctx, "delete-repo", trace.WithAttributes(
		attribute.String("repo", repo),
	))
	defer span.End()

	p := filepath.Join(s.root, "repos", repo)
	err := os.RemoveAll(p)
	return err
}

func (s *server) deleteArch(ctx context.Context, repo, arch string) error {
	_, span := s.t.Start(ctx, "delete-arch", trace.WithAttributes(
		attribute.String("repo", repo),
		attribute.String("arch", arch),
	))
	defer span.End()

	p := filepath.Join(s.root, "repos", repo, arch)
	err := os.RemoveAll(p)
	return err
}

func (s *server) deletePackage(ctx context.Context, repo, arch, pkg string) error {
	ctx, span := s.t.Start(ctx, "delete-pkg", trace.WithAttributes(
		attribute.String("repo", repo),
		attribute.String("arch", arch),
		attribute.String("pkg", pkg),
	))
	defer span.End()

	p := filepath.Join(s.root, "repos", repo, arch, pkg)
	err := os.RemoveAll(p)
	if err != nil {
		return err
	}

	ctx, span = s.t.Start(ctx, "repo-remove")
	defer span.End()
	dbp := filepath.Join(filepath.Dir(p), repo+".db.tar.zst")
	cmd := exec.CommandContext(ctx, "repo-remove", dbp, p)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", cmd, err, string(out))
	}
	return nil
}
