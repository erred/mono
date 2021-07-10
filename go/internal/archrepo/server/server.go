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
)

const envBearerToken = "AR_TOKENS" // comma separated list of allowed bearer tokens

type Options struct {
	Dir string
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	o.InitFlags(fs)
	return &o
}

func (o *Options) InitFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Dir, "data", "/data", "path to store data")
}

type server struct {
	root          string
	allowedTokens map[string]struct{}
}

func New(ctx context.Context, opt *Options) (http.Handler, error) {
	s := &server{
		root:          opt.Dir,
		allowedTokens: make(map[string]struct{}),
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
}

func (s *server) auth(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := logr.FromContext(ctx).WithName("auth")
		ah := r.Header.Get("authorization")
		if _, ok := s.allowedTokens[strings.TrimPrefix(ah, "Bearer ")]; !strings.HasPrefix(ah, "Bearer ") || !ok {
			l.Info("no auth", "method", r.Method, "path", r.URL.Path, "user-agent", r.UserAgent())
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *server) addPackageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := logr.FromContext(ctx).WithName("add")
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
	defer rc.Close()
	p := filepath.Join(s.root, "repos", repo, arch)
	os.MkdirAll(p, 0o755) // ignore errors
	p = filepath.Join(p, pkg)
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	if err != nil {
		return err
	}
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
	l := logr.FromContext(ctx).WithName("delete")
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
	p := filepath.Join(s.root, "repos", repo)
	err := os.RemoveAll(p)
	return err
}

func (s *server) deleteArch(ctx context.Context, repo, arch string) error {
	p := filepath.Join(s.root, "repos", repo, arch)
	err := os.RemoveAll(p)
	return err
}

func (s *server) deletePackage(ctx context.Context, repo, arch, pkg string) error {
	p := filepath.Join(s.root, "repos", repo, arch, pkg)
	err := os.RemoveAll(p)
	if err != nil {
		return err
	}

	dbp := filepath.Join(filepath.Dir(p), repo+".db.tar.zst")
	cmd := exec.CommandContext(ctx, "repo-remove", dbp, p)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w: %s", cmd, err, string(out))
	}
	return nil
}
