package render

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"sync"
	"text/template"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.seankhliao.com/mono/internal/web/picture"
	"sigs.k8s.io/yaml"
)

var (
	//go:embed layout.tpl
	layoutTpl string
	//go:embed base.css
	baseCssTpl string

	T = template.Must(
		template.Must(
			template.New("basecss").Parse(baseCssTpl),
		).New("").Parse(layoutTpl),
	)
)

var defaultMarkdown = goldmark.New(
	goldmark.WithExtensions(
		extension.Strikethrough,
		extension.Table,
		extension.TaskList,
		picture.Picture,
	),
	goldmark.WithParserOptions(
		parser.WithHeadingAttribute(), // {#some-id}
		parser.WithAutoHeadingID(),    // based on heading
	),
	goldmark.WithRendererOptions(html.WithUnsafe()),
)

var (
	defaultMinifyOnce sync.Once
	defaultMinifyM    *minify.M
	defaultMinify     = func() *minify.M {
		defaultMinifyOnce.Do(func() {
			m := minify.New()
			m.AddFunc("text/html", mhtml.Minify)
			m.AddFunc("text/css", css.Minify)
			m.AddFunc("image/svg+xml", svg.Minify)
			m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
			defaultMinifyM = m
		})
		return defaultMinifyM
	}
)

type Options struct {
	Markdown goldmark.Markdown

	Minify     *minify.M
	MinifySkip bool
	Template   *template.Template

	Data PageData
}

func (o *Options) init() {
	if o.Markdown == nil {
		o.Markdown = defaultMarkdown
	}
	if o.Minify == nil {
		o.Minify = defaultMinify()
	}
	if o.Template == nil {
		o.Template = T
	}
}

type PageData struct {
	Compact      bool
	Date         string // for blog posts
	GTMID        string // for analytics
	URLCanonical string

	// inserted early in html head
	Head string
	// literal HTML body, replaced with rendered markdown content
	// if Render body is non empty
	Main string

	// Extracted from body
	Title       string
	Description string
	Style       string
	H1          string // Only for non compact
	H2          string // Only for non compact
}

// fromHeader tries to fill in values from a yaml header delimited by "---\n",
// returning the remaining bytes
func (d *PageData) fromHeader(b []byte) ([]byte, error) {
	b = bytes.TrimSpace(b)
	if !bytes.HasPrefix(b, []byte("---\n")) { // no header marker
		return b, nil
	}
	b = b[4:]
	i := bytes.Index(b, []byte("---\n"))
	if i == -1 { // no header ending
		return b, nil
	}
	meta := make(map[string]string)
	err := yaml.Unmarshal(b[:i], &meta)
	if err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}

	d.Head = first(d.Head, meta["head"])
	d.Title = first(d.Title, meta["title"])
	d.Description = first(d.Description, meta["description"])
	d.Style = first(d.Style, meta["style"])
	d.H1 = first(d.H1, meta["h1"])
	d.H2 = first(d.H2, meta["h2"])

	return b[i+4:], nil
}

// first returns the first non empty string
func first(args ...string) string {
	for _, s := range args {
		if s != "" {
			return s
		}
	}
	return ""
}
