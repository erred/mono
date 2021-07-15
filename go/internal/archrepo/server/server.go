package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/go/pacmandb"
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
	err = repoAdd(ctx, dbp, p)
	if err != nil {
		return fmt.Errorf("add %s to db: %w", pkg, err)
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
	err = repoRemove(ctx, dbp, pkg)
	if err != nil {
		return fmt.Errorf("remove %s from db: %w", pkg, err)
	}
	return nil
}

func repoRemove(ctx context.Context, repoPath, pkgName string) error {
	db, files, err := readRepos(repoPath)
	if err != nil {
		return fmt.Errorf("read repo %s: %w", repoPath, err)
	}

	delete(db.Pkgs, pkgName)
	delete(files.Pkgs, pkgName)

	err = writeRepos(repoPath, db, files)
	if err != nil {
		return fmt.Errorf("write repos %s: %w", repoPath, err)
	}
	return nil
}

func repoAdd(ctx context.Context, repoPath, pkgPath string) error {
	pkgF, err := os.Open(pkgPath)
	if err != nil {
		return fmt.Errorf("open package %s: %w", pkgPath, err)
	}
	pkgName := filepath.Base(pkgPath)
	pkg, err := pacmandb.ParsePackage(pkgName, pkgF)
	if err != nil {
		return fmt.Errorf("parse package %s: %w", pkgPath, err)
	}

	db, files, err := readRepos(repoPath)
	if err != nil {
		return fmt.Errorf("read repo %s: %w", repoPath, err)
	}

	i := strings.LastIndex(pkgName, "-")
	pkgName = pkgName[:i] // trim -x86_64.pkg.tar.zst

	pkg2 := *pkg
	pkg2.Files = nil
	db.Pkgs[pkgName] = pkg2
	files.Pkgs[pkgName] = *pkg

	err = writeRepos(repoPath, db, files)
	if err != nil {
		return fmt.Errorf("write repos %s: %w", repoPath, err)
	}
	return nil
}

func repoPaths(repoPath string) (db, files string) {
	dir := filepath.Dir(repoPath)
	repo := strings.Split(filepath.Base(repoPath), ".")[0]
	db = filepath.Join(dir, repo+".db.tar.zst")
	files = filepath.Join(dir, repo+".files.tar.zst")
	return db, files
}

func readRepos(repoPath string) (db, files *pacmandb.DB, err error) {
	dbPath, filesPath := repoPaths(repoPath)

	db, err = readRepo(dbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read db: %w", err)
	}
	files, err = readRepo(filesPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read files: %w", err)
	}
	return db, files, nil
}

func readRepo(repoPath string) (*pacmandb.DB, error) {
	f, err := os.Open(repoPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return &pacmandb.DB{Pkgs: make(map[string]pacmandb.Package)}, nil
		}
		return nil, fmt.Errorf("open %s: %w", repoPath, err)
	}
	defer f.Close()
	var db pacmandb.DB
	err = db.DecodeZstd(f)
	if err != nil {
		return nil, fmt.Errorf("decode database %s: %w", repoPath, err)
	}
	return &db, nil
}

func writeRepos(repoPath string, db, files *pacmandb.DB) error {
	dbPath, filesPath := repoPaths(repoPath)

	err := writeRepo(dbPath, db)
	if err != nil {
		return fmt.Errorf("write db: %w", err)
	}
	err = writeRepo(filesPath, files)
	if err != nil {
		return fmt.Errorf("write files: %w", err)
	}
	return nil
}

func writeRepo(repoPath string, repo *pacmandb.DB) error {
	tmpPath := repoPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", tmpPath, err)
	}
	defer f.Close()
	err = repo.EncodeZstd(f)
	if err != nil {
		return fmt.Errorf("encode db: %w", err)
	}
	err = os.Rename(tmpPath, repoPath)
	if err != nil {
		return fmt.Errorf("rename db: %w", err)
	}

	symlink := strings.TrimSuffix(repoPath, ".tar.zst")
	os.Remove(symlink)
	err = os.Symlink(repoPath, symlink)
	if err != nil {
		return fmt.Errorf("symlink db: %w", err)
	}
	return nil
}
