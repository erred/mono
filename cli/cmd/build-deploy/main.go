package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

type options struct {
	gocache   string
	gopath    string
	workspace string
	ref       string
}

func newOptins() (*options, error) {
	workspace, err := os.MkdirTemp("", "mono-build-*")
	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}
	return &options{
		workspace: workspace,
		ref:       "main",
	}, nil
}

func run() error {
	o, err := newOptins()
	if err != nil {
		return fmt.Errorf("configure: %w", err)
	}

	defer cleanup(o)

	err = prepare(o)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}

	err = build(o)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	err = deploy(o)
	if err != nil {
		return fmt.Errorf("deploy: %w", err)
	}

	return nil
}

func prepare(o *options) error {
	eg, _ := errgroup.WithContext(context.TODO())
	eg.Go(func() error {
		_, err := git.PlainClone(filepath.Join(o.workspace, "mono"), false, &git.CloneOptions{
			URL:           "https://github.com/seankhliao/mono",
			ReferenceName: plumbing.ReferenceName(o.ref),
			SingleBranch:  true,
			Depth:         1,
		})
		if err != nil {
			return fmt.Errorf("clone mono: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		res, err := http.Get("https://go.dev/dl/?mode=json")
		if err != nil {
			return fmt.Errorf("get go versions: %w", err)
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("get go versions: %s", res.Status)
		}
		defer res.Body.Close()
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("read go versions: %w", err)
		}
		type goVersions []struct {
			Version string
			Files   []struct {
				OS       string
				Arch     string
				Filename string
				SHA256   string
			}
		}
		var versions goVersions
		err = json.Unmarshal(b, &versions)
		if err != nil {
			return fmt.Errorf("unmarshal go versions")
		}
		var fn string
		for _, file := range versions[0].Files {
			if file.OS == runtime.GOOS && file.Arch == runtime.GOARCH {
				fn = file.Filename
				break
			}
		}
		if fn == "" {
			return fmt.Errorf("no suitable go release was found")
		}

		res, err = http.Get("https://go.dev/dl/" + fn)
		if err != nil {
			return fmt.Errorf("get go release: %w", err)
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("get go release: %s", res.Status)
		}
		defer res.Body.Close()
		gr, err := gzip.NewReader(res.Body)
		if err != nil {
			return fmt.Errorf("create go release gzip reader: %w", err)
		}
		tr := tar.NewReader(gr)
		for th, err := tr.Next(); err == nil; th, err = tr.Next() {
			func() error {
				if th.Typeflag == tar.TypeDir {
					fp := filepath.Join(o.workspace, th.Name)
					err := os.MkdirAll(fp, 0o755)
					if err != nil {
						return fmt.Errorf("mkdir %s: %w", fp, err)
					}
				} else if th.Typeflag == tar.TypeReg {
					fn := filepath.Join(o.workspace, th.Name)
					f, err := os.Create(fn)
					if err != nil {
						return fmt.Errorf("create %s: %w", fn, err)
					}
					defer f.Close()
					_, err = io.Copy(f, tr)
					if err != nil {
						return fmt.Errorf("write %s: %w", fn, err)
					}
					os.Chmod(fn, th.FileInfo().Mode())
				} else {
					return fmt.Errorf("unknown file type: %v", th.Typeflag)
				}
				return nil
			}()
		}
		return nil
	})

	return eg.Wait()
}

func build(o *options) error {
}

func deploy(o *options) error {
}

func cleanup(o *options) {
	err := os.RemoveAll(o.workspace)
	if err != nil {
		log.Println("cleanup:", err)
	}
}

func runCmd(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Env
}
