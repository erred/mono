package server

import (
	"archive/tar"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
	"go.seankhliao.com/mono/go/pacmandb"
)

type server struct {
	c *http.Client
}

type PkgVersion struct {
	GoPkg   string
	Version string

	PkgInfo pacmandb.Desc
}

func (s *server) go2repo(ctx context.Context, repoURLBase string, p *PkgVersion) error {
	bins, err := goInstall(ctx, p.GoPkg, p.Version)
	if err != nil {
		return fmt.Errorf("install: %w", err)
	}
	defer os.RemoveAll(filepath.Dir(bins[0]))
	f, err := packageBins(ctx, bins)
	if err != nil {
		return fmt.Errorf("package: %w", err)
	}
	defer f.Close()
	defer os.RemoveAll(f.Name())
	u := fmt.Sprintf("%s/%s-%s-1.pkg.tar.zst", repoURLBase, p.PkgInfo.Name, p.Version)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, f)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	res, err := s.c.Do(req)
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	return nil
}

func goInstall(ctx context.Context, pkg, version string) ([]string, error) {
	tmpDir, err := os.MkdirTemp("", "goinstall")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	cmd := exec.CommandContext(ctx, "go", "install", pkg+"@"+version)
	cmd.Env = []string{
		"GOBIN=" + tmpDir,
	}
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("install %s @ %s: %w", pkg, version, err)
	}
	des, err := os.ReadDir(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", tmpDir, err)
	}
	var ps []string
	for _, de := range des {
		ps = append(ps, de.Name())
	}
	return ps, nil
}

func packageBins(ctx context.Context, bins []string) (*os.File, error) {
	// TODO: .PKGINFO
	f, err := os.CreateTemp("", "gopkg")
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}
	zenc, err := zstd.NewWriter(f)
	if err != nil {
		return nil, fmt.Errorf("create zstd encoder: %w", err)
	}
	defer zenc.Close()
	tw := tar.NewWriter(zenc)
	for _, bin := range bins {
		err = tarFile(tw, bin)
		if err != nil {
			return nil, fmt.Errorf("tar file: %w", err)
		}
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek 0: %w", err)
	}
	return f, nil
}

func tarFile(tw *tar.Writer, fn string) error {
	err := tw.WriteHeader(&tar.Header{
		Name: filepath.Join("/usr/bin", filepath.Base(fn)),
		Mode: 0o755,
	})
	if err != nil {
		return fmt.Errorf("write header %s: %w", fn, err)
	}
	f, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("open %s: %w", fn, err)
	}
	defer f.Close()
	_, err = io.Copy(tw, f)
	if err != nil {
		return fmt.Errorf("copy %s: %w", fn, err)
	}
	return nil
}
