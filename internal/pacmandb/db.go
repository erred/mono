package pacmandb

import (
	"archive/tar"
	"errors"
	"io"
	"path/filepath"
	"sort"

	"github.com/klauspost/compress/zstd"
)

// DB represents a pacman database file,
// it could be one of $repo.db or $repo.files
type DB struct {
	Pkgs map[string]Package
}

// DecodeZstd reads a zstd compressed database
func (db *DB) DecodeZstd(r io.Reader) error {
	zdec, err := zstd.NewReader(r)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zdec)
	pkgs, err := readTar(tr)
	if err != nil {
		return err
	}
	db.Pkgs = pkgs
	return nil
}

// EncodeZstd writes out the database into a zstd compressed file
func (db *DB) EncodeZstd(w io.Writer) error {
	zenc, err := zstd.NewWriter(w)
	if err != nil {
		return err
	}
	defer zenc.Close()
	tw := tar.NewWriter(zenc)
	err = writeTar(db.Pkgs, tw)
	if err != nil {
		return err
	}
	return nil
}

type entry struct {
	hdr  *tar.Header
	data []byte
}

// readTar loops through a tar archive of a database
// and decodes it into Packages, keyed by the name
func readTar(r *tar.Reader) (map[string]Package, error) {
	pkgs := make(map[string]Package)
	for {
		hdr, err := r.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if hdr.Typeflag == tar.TypeDir {
			continue
		} else if hdr.Typeflag != tar.TypeReg {
			continue
		}
		b, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}

		pkgName := filepath.Dir(hdr.Name)
		pkg := pkgs[pkgName]
		switch filepath.Base(hdr.Name) {
		case "desc":
			var desc Desc
			err = desc.UnmarshalText(b)
			if err != nil {
				return nil, err
			}
			pkg.Desc = desc
		case "files":
			var files Files
			err = files.UnmarshalText(b)
			if err != nil {
				return nil, err
			}
			pkg.Files = &files
		default:
			// ignore unknown files
		}
		pkgs[pkgName] = pkg

	}
	return pkgs, nil
}

// writeTar writes the Packages into the tar archive
// in sorted lexical order.
// The files entry is written if it is not nil
func writeTar(pkgs map[string]Package, w *tar.Writer) error {
	// stable tar entries order
	pkgNames := make([]string, 0, len(pkgs))
	for pkgName := range pkgs {
		pkgNames = append(pkgNames, pkgName)
	}
	sort.Strings(pkgNames)

	for _, pkgName := range pkgNames {
		pkg := pkgs[pkgName]
		err := w.WriteHeader(&tar.Header{
			Name: pkgName + "/",
			Mode: 0o755,
		})
		if err != nil {
			return err
		}

		desc, _ := pkg.Desc.MarshalText()
		err = w.WriteHeader(&tar.Header{
			Name: filepath.Join(pkgName, "desc"),
			Mode: 0o644,
			Size: int64(len(desc)),
		})
		if err != nil {
			return err
		}
		_, err = w.Write(desc)
		if err != nil {
			return err
		}

		if pkg.Files == nil {
			continue
		}
		files, _ := pkg.Files.MarshalText()
		err = w.WriteHeader(&tar.Header{
			Name: filepath.Join(pkgName, "files"),
			Mode: 0o644,
			Size: int64(len(files)),
		})
		if err != nil {
			return err
		}
		_, err = w.Write(files)
		if err != nil {
			return err
		}
	}
	err := w.Close()
	if err != nil {
		return err
	}
	return nil
}
