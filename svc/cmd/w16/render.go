package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strings"

	"go.seankhliao.com/mono/internal/web/render"
)

type pageInfo struct {
	path  string
	title string
	date  string // only for blog posts
}

var dateRe = regexp.MustCompile(`\d{5}-\d{2}-\d{2}`)

func (o *Options) renderAndRegister(mux *http.ServeMux, fsys fs.FS) error {
	var pis []pageInfo

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.Type() != 0 {
			return nil
		}

		if path.Ext(p) != ".md" {
			o.handleFsysFile(mux, fsys, p)
			return nil
		}

		b, pi, err := o.renderFile(fsys, p)
		if err != nil {
			return fmt.Errorf("render %s: %w", p, err)
		}

		pis = append(pis, pi)
		o.handleBytes(mux, pi.path, "index.html", b)

		if strings.HasPrefix(pi.path, "/blog/") {
			// register a redirect for old urls
			// old: /blog/2020-xxx -> /new: blog/12020-xxx
			op := "/blog/" + pi.path[7:]
			mux.Handle(op, http.RedirectHandler(pi.path, http.StatusMovedPermanently))
		} else if pi.path == "/404" {
			o.notFoundBody = b
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("walk render: %w", err)
	}

	sort.Slice(pis, func(i, j int) bool {
		return pis[i].path > pis[j].path
	})

	pi, err := o.blogIndex(mux, pis)
	if err != nil {
		if !errors.Is(err, errSkip) {
			return fmt.Errorf("render blog index: %w", err)
		}
	} else {
		pis = append(pis, pi)
	}

	o.sitemap(mux, pis)
	return nil
}

func (o *Options) renderFile(fsys fs.FS, p string) ([]byte, pageInfo, error) {
	var buf bytes.Buffer
	f, _ := fsys.Open(p)
	defer f.Close()

	cp := canonicalPath(p)
	cu := fmt.Sprintf("https://%s%s", o.Hostname, cp)

	var h1, date string
	if strings.HasPrefix(cp, "/blog/") {
		h1 = `<a href="/blog/">b<em>log</em></a>`
		date = dateRe.FindString(cp)
	}

	ro := render.Options{
		Data: render.PageData{
			URLCanonical: cu,
			GTMID:        o.GTMID,
			Compact:      o.Compact,
			H1:           h1,
			H2:           date,
		},
	}

	err := render.Render(&ro, &buf, f)
	if err != nil {
		return nil, pageInfo{}, err
	}

	pi := pageInfo{
		path:  cp,
		title: ro.Data.Title,
		date:  date,
	}

	return buf.Bytes(), pi, nil
}

var errSkip = errors.New("skip for no entries")

// blogIndex renders a blog index page
func (o *Options) blogIndex(mux *http.ServeMux, pis []pageInfo) (pageInfo, error) {
	var blogEntries []pageInfo
	for _, pi := range pis {
		if strings.HasPrefix(pi.path, "/blog/") {
			blogEntries = append(blogEntries, pi)
		}
	}
	if len(blogEntries) == 0 {
		return pageInfo{}, errSkip
	}

	var body bytes.Buffer
	body.WriteString(`
<h3><em>B</em>log</h3>
<p>we<em>b log</em> of things that never made sense,
maybe someone will find this useful</p>
<ul>
`)
	for _, be := range blogEntries {
		fmt.Fprintf(
			&body,
			`<li><time datetime="%s">%s</time> | <a href="%s">%s</a></li>`+"\n",
			be.date[1:],
			be.date,
			be.path,
			be.title,
		)
	}
	body.WriteString(`</ul>`)

	ro := render.Options{
		MarkdownSkip: true,
		Data: render.PageData{
			URLCanonical: fmt.Sprintf("https://%s/blog/", o.Hostname),
			GTMID:        o.GTMID,
			Title:        "blog | seankhliao",
			Description:  "list of things i wrote",
			H1:           `<a href="/blog/">b<em>log</em></a>`,
			H2: `Artisanal, <em>hand-crafted</em> blog posts
imbued with delayed <em>regrets</em>`,
			Style: `
ul li {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
}`,
		},
	}

	var buf bytes.Buffer
	err := render.Render(&ro, &buf, &body)
	if err != nil {
		return pageInfo{}, fmt.Errorf("render blog index: %w", err)
	}

	o.handleBytes(mux, "/blog/", "blog.html", buf.Bytes())
	return pageInfo{
		path:  "/blog/",
		title: ro.Data.Title,
	}, nil
}

// sitemap renders a plaintext list of valid urls
func (o *Options) sitemap(mux *http.ServeMux, pis []pageInfo) {
	var buf bytes.Buffer
	for _, pi := range pis {
		fmt.Fprintf(&buf, "https://%s%s\n", o.Hostname, pi.path)
	}

	o.handleBytes(mux, "/sitemap.txt", "map.txt", buf.Bytes())
}
