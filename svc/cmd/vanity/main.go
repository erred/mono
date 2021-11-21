package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"go.seankhliao.com/mono/content"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/httpsvr"
	"go.seankhliao.com/mono/svc/o11y"
)

const (
	redirectURL = "https://seankhliao.com/"
)

func main() {
	oo := o11y.NewOptions(flag.CommandLine)
	ho := httpsvr.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := oo.New()
	ho.BaseContext = ctx

	ho.Handler = New(ctx)

	ho.Run()
}

func New(ctx context.Context) http.Handler {
	l := logr.FromContextOrDiscard(ctx).WithName("vanity")

	indexF, _ := content.Content.Open("go.seankhliao.com/index.md")
	defer indexF.Close()
	var buf bytes.Buffer
	err := render.Render(&render.Options{
		Data: render.PageData{
			URLCanonical: "https://go.seankhliao.com/",
			Compact:      true,
			Title:        "go.seankhliao.com",
		},
	}, &buf, indexF)
	if err != nil {
		l.Error(err, "render index")
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeContent(rw, r, "index.html", time.Time{}, bytes.NewReader(buf.Bytes()))
			return
		}

		repo := strings.Split(r.URL.Path, "/")[1]
		err := render.Render(
			&render.Options{
				MarkdownSkip: true,
				Data: render.PageData{
					URLCanonical: "https://go.seankhliao.com/" + repo,
					Compact:      true,
					Title:        "go.seankhliao.com/" + repo,
					Head: fmt.Sprintf(`
<meta
  name="go-import"
  content="go.seankhliao.com/%[1]s git https://github.com/seankhliao/%[1]s">
<meta
  name="go-source"
  content="go.seankhliao.com/%[1]s
    https://github.com/seankhliao/%[1]s
    https://github.com/seankhliao/%[1]s/tree/master{/dir}
    https://github.com/seankhliao/%[1]s/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="5;url=https://pkg.go.dev/go.seankhliao.com/%[1]s" />`, repo),
				},
			},
			rw,
			strings.NewReader(fmt.Sprintf(`
<h3><em>go.seankhliao.com/</em>%[1]s</h3>

<p><em>source:</em>
  <a href="https://github.com/seankhliao/%[1]s">github</a>
<p><em>docs:</em>
  <a href="https://pkg.go.dev/go.seankhliao.com/%[1]s">pkg.go.dev</a>
`, repo)),
		)
		if err != nil {
			l.Error(err, "render")
		}
	})
}
