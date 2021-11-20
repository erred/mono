package pacmandb

import (
	"archive/tar"
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/klauspost/compress/zstd"
)

// Package is the representation of a package within a database
type Package struct {
	Desc  Desc
	Files *Files
}

// ParsePackage parses a pacman package for the info (from .PKGINFO)
// to populate a database entry.
// The filename arg is used to populate Desc.
// PGPSIG is never populated.
func ParsePackage(filename string, r io.Reader) (*Package, error) {
	md5hash := md5.New()
	sha256hash := sha256.New()
	var sizeCtr sizeCounter
	r = io.TeeReader(r, io.MultiWriter(md5hash, sha256hash, &sizeCtr))

	zdec, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(zdec)

	var filenames []string
	var desc Desc
	for {
		hdr, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if !strings.HasPrefix(hdr.Name, ".") {
			filenames = append(filenames, hdr.Name)
		}

		if path.Base(hdr.Name) == ".PKGINFO" {
			desc, err = parsePKGINFO(tr)
			if err != nil {
				return nil, err
			}
		}
	}

	sort.Strings(filenames)
	desc.Filename = filename
	desc.MD5Sum = hex.EncodeToString(md5hash.Sum(nil))
	desc.SHA256Sum = hex.EncodeToString(sha256hash.Sum(nil))
	desc.CSize = strconv.Itoa(sizeCtr.n)

	return &Package{
		Desc: desc,
		Files: &Files{
			Filenames: filenames,
		},
	}, nil
}

// parses .PKGINFO to get a partial Desc
func parsePKGINFO(r io.Reader) (Desc, error) {
	// based the first half of db_write_entry() in repo-add
	var desc Desc
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if bytes.HasPrefix(sc.Bytes(), []byte("#")) {
			continue
		}
		bb := bytes.SplitN(sc.Bytes(), []byte("="), 2)
		if len(bb) != 2 {
			return Desc{}, fmt.Errorf("invalid line: %q", sc.Text())
		}
		val := string(bytes.TrimSpace(bb[1]))
		switch string(bytes.TrimSpace(bb[0])) {
		case "pkgname":
			desc.Name = val
		case "pkgbase":
			desc.Base = val
		case "pkgver":
			desc.Version = val
		case "pkgdesc":
			desc.Desc = val
		case "size":
			desc.ISize = val
		case "url":
			desc.URL = val
		case "arch":
			desc.Arch = val
		case "builddate":
			desc.BuildDate = val
		case "packager":
			desc.Packager = val
		case "group":
			desc.Groups = append(desc.Groups, val)
		case "license":
			desc.License = append(desc.License, val)
		case "replaces":
			desc.Replaces = append(desc.Replaces, val)
		case "depend":
			desc.Depends = append(desc.Depends, val)
		case "conflict":
			desc.Conflicts = append(desc.Conflicts, val)
		case "provides":
			desc.Provides = append(desc.Provides, val)
		case "optdepend":
			desc.OptDepends = append(desc.OptDepends, val)
		case "makedepend":
			desc.MakeDepends = append(desc.MakeDepends, val)
		case "checkdepend":
			desc.CheckDepends = append(desc.CheckDepends, val)
		default:
			// ignore unknown fields
		}
	}
	return desc, nil
}

// Desc is the desc file in databases
type Desc struct {
	Filename string
	Name     string
	Base     string
	Version  string
	Desc     string
	Groups   []string
	CSize    string // computed over input
	ISize    string

	MD5Sum    string // computed over input
	SHA256Sum string // computed over input

	PGPSig string // unset

	URL       string
	License   []string
	Arch      string
	BuildDate string
	Packager  string
	Replaces  []string
	Conflicts []string
	Provides  []string

	Depends      []string
	OptDepends   []string
	MakeDepends  []string
	CheckDepends []string
}

// MarshalText formats the data into the desc format for databases,
// unset fields are skipped
func (d Desc) MarshalText() ([]byte, error) {
	// based the second half of db_write_entry() in repo-add
	var b bytes.Buffer
	formatEntry(&b, "FILENAME", d.Filename)
	formatEntry(&b, "NAME", d.Name)
	formatEntry(&b, "BASE", d.Base)
	formatEntry(&b, "VERSION", d.Version)
	formatEntry(&b, "DESC", d.Desc)
	formatEntry(&b, "GROUPS", d.Groups...)
	formatEntry(&b, "CSIZE", d.CSize)
	formatEntry(&b, "ISIZE", d.ISize)
	formatEntry(&b, "MD5SUM", d.MD5Sum)
	formatEntry(&b, "SHA256SUM", d.SHA256Sum)
	formatEntry(&b, "PGPSIG", d.PGPSig)
	formatEntry(&b, "URL", d.URL)
	formatEntry(&b, "LICENSE", d.License...)
	formatEntry(&b, "ARCH", d.Arch)
	formatEntry(&b, "BUILDDATE", d.BuildDate)
	formatEntry(&b, "PACKAGER", d.Packager)
	formatEntry(&b, "REPLACES", d.Replaces...)
	formatEntry(&b, "CONFLICTS", d.Conflicts...)
	formatEntry(&b, "PROVIDES", d.Provides...)
	formatEntry(&b, "DEPENDS", d.Depends...)
	formatEntry(&b, "OPTDEPENDS", d.OptDepends...)
	formatEntry(&b, "MAKEDEPENDS", d.MakeDepends...)
	formatEntry(&b, "CHECKDEPENDS", d.CheckDepends...)
	return b.Bytes(), nil
}

// UnmarshalText extracts the data from the desc format in databases.
// All fields are overridden.
func (d *Desc) UnmarshalText(b []byte) error {
	m := scanFile(b)
	d.Filename = slice2string(m["FILENAME"])
	d.Name = slice2string(m["NAME"])
	d.Base = slice2string(m["BASE"])
	d.Version = slice2string(m["VERSION"])
	d.Desc = slice2string(m["DESC"])
	d.Groups = m["GROUPS"]
	d.CSize = slice2string(m["CSIZE"])
	d.ISize = slice2string(m["ISIZE"])
	d.MD5Sum = slice2string(m["MD5SUM"])
	d.SHA256Sum = slice2string(m["SHA256SUM"])
	d.PGPSig = slice2string(m["PGPSIG"])
	d.URL = slice2string(m["URL"])
	d.License = m["LICENSE"]
	d.Arch = slice2string(m["ARCH"])
	d.BuildDate = slice2string(m["BUILDDATE"])
	d.Packager = slice2string(m["PACKAGER"])
	d.Replaces = m["REPLACES"]
	d.Conflicts = m["CONFLICTS"]
	d.Provides = m["PROVIDES"]
	d.Depends = m["DEPENDS"]
	d.OptDepends = m["OPTDEPENDS"]
	d.MakeDepends = m["MAKEDEPENDS"]
	d.CheckDepends = m["CHECKDEPENDS"]
	return nil
}

// Files is the files entry in databases
type Files struct {
	Filenames []string
}

// MarshalText formats the data for databases
func (f Files) MarshalText() ([]byte, error) {
	var b bytes.Buffer
	formatEntry(&b, "FILES", f.Filenames...)
	return b.Bytes(), nil
}

// UnmarshalText parses the data from database format
func (f *Files) UnmarshalText(b []byte) error {
	m := scanFile(b)
	f.Filenames = m["FILES"]
	return nil
}

// scanFile parses the input in desc/files format
// returning a map of entry names and values
func scanFile(b []byte) map[string][]string {
	m := make(map[string][]string)
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || s[0] != '%' {
			continue
		}
		key := s[1 : len(s)-1]
		var vals []string
		for sc.Scan() {
			s := strings.TrimSpace(sc.Text())
			if s == "" {
				break
			}
			vals = append(vals, s)
		}
		m[key] = vals
	}
	return m
}

// formatEntry formats a single entry in database format
func formatEntry(b *bytes.Buffer, name string, vals ...string) {
	if (len(vals) == 0) || (len(vals) == 1 && vals[0] == "") {
		return
	}
	fmt.Fprintf(b, "%%%s%%\n", name)
	for _, val := range vals {
		b.WriteString(val)
		b.WriteRune('\n')
	}
	b.WriteRune('\n')
}

// sizeCounter counts the bytes written to it
type sizeCounter struct {
	n int
}

func (s *sizeCounter) Write(p []byte) (int, error) {
	s.n += len(p)
	return len(p), nil
}

// slice2string returns the first element or the empty string if the slice is empty
func slice2string(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return ss[0]
}
