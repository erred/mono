package filesystem

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var nameRe = regexp.MustCompile(`^[a-z][a-z0-9-]{1,62}[a-z0-9]$`)

type Repository struct {
	Name string
}

type PackageVersion struct {
	Name string
	Data []byte
}

type Store struct {
	root string
}

func New(ctx context.Context, root string) (*Store, error) {
	err := os.MkdirAll(root, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, fmt.Errorf("ensure %s: %w", root, err)
	}
	return &Store{
		root: root,
	}, nil
}

func (s Store) ListRepositories(ctx context.Context) ([]Repository, error) {
	p := filepath.Join(s.root, "repos")
	fis, err := os.ReadDir(p)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", p, err)
	}

	var repos []Repository
	for _, fi := range fis {
		if fi.IsDir() {
			repos = append(repos, Repository{
				Name: fi.Name(),
			})
		}
	}

	return repos, nil
}

func (s Store) UpdateRepository(ctx context.Context, repo Repository) error {
	if !nameRe.MatchString(repo.Name) {
		return fmt.Errorf("invalid name: %s", repo.Name)
	}

	p := filepath.Join(s.root, "repos", repo.Name, "os", "x86_64")
	err := os.MkdirAll(p, 0o755)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("ensure %s: %w", p, err)
	}

	return nil
}

func (s Store) DeleteRepository(ctx context.Context, repo Repository) error {
	if !nameRe.MatchString(repo.Name) {
		return fmt.Errorf("invalid name: %s", repo.Name)
	}

	p := filepath.Join(s.root, "repos", repo.Name)
	err := os.RemoveAll(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", s, err)
	}

	return nil
}

func (s Store) UpdatePackageVersion(ctx context.Context, repo Repository, pkg PackageVersion) error {
	p := filepath.Join(s.root, "repos", repo.Name, "os", "x86_64", pkg.Name)
	err := os.Remove(p)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", p, err)
	}
	err = os.WriteFile(p, pkg.Data, 0o644)
	if err != nil {
		return fmt.Errorf("write %s: %w", p, err)
	}

	dbp := filepath.Join(s.root, "repos", repo.Name, "os", "x86_64", repo.Name+".db.tar.zst")
	cmd := exec.CommandContext(ctx, "repo-add", dbp, p)
	err = cmd.Run()
	if err != nil {
		if eerr := (&exec.ExitError{}); errors.As(err, &eerr) {
			return fmt.Errorf("repo-add %s: %v: %w", p, eerr.Stderr, err)
		}
		return fmt.Errorf("repo-add %s: %w", p, err)
	}

	return nil
}

func (s Store) DeletePackageVersion(ctx context.Context, repo Repository, pkg PackageVersion) error {
	p := filepath.Join(s.root, "repos", repo.Name, "os", "x86_64", pkg.Name)
	dbp := filepath.Join(s.root, "repos", repo.Name, "os", "x86_64", repo.Name+".db.tar.zst")
	cmd := exec.CommandContext(ctx, "repo-remove", dbp, p)
	err := cmd.Run()
	if err != nil {
		if eerr := (&exec.ExitError{}); errors.As(err, &eerr) {
			return fmt.Errorf("repo-remove %s: %v: %w", p, eerr.Stderr, err)
		}
		return fmt.Errorf("repo-remove %s: %w", p, err)
	}

	err = os.Remove(p)
	if err != nil {
		return fmt.Errorf("remove %s: %w", p, err)
	}

	return nil
}
