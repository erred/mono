package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"go.seankhliao.com/mono/svc/internal/testhttp"
)

func TestRenderFs(t *testing.T) {
	fsys := fstest.MapFS{
		"blog/12021-01-02-xxx.md": &fstest.MapFile{
			Data: []byte(`---
title: xxx
description: hello world
---

### hello world

foo bar
`),
		},
		"blog/12021-05-01-yyy.md": &fstest.MapFile{
			Data: []byte(`---
title: yyy yy
description: fizz buzz
---

### _blog_ post 1

example blog post
`),
		},
		"hello.md": &fstest.MapFile{
			Data: []byte(`---
title: hello
description: not a blog page
---

### fizz buzz
`),
		},
	}

	o := Options{
		Hostname: "example.com",
	}
	mux := http.NewServeMux()
	err := o.renderAndRegister(mux, fsys)
	if err != nil {
		t.Errorf("renderAndRegister err=%v", err)
		return
	}

	svr := httptest.NewServer(mux)
	defer svr.Close()

	testhttp.Expect(t, svr.Client(), "blogpage",
		http.MethodGet, svr.URL+"/blog/12021-05-01-yyy/", http.StatusOK,
		map[string]string{"content-type": "text/html; charset=utf-8"},
		[]string{
			`<meta name=description content="fizz buzz">`,
			`<h3 id=-blog--post-1><em>blog</em> post 1</h3>`,
			`example blog post`,
		},
	)

	testhttp.Expect(t, svr.Client(), "otherpage",
		http.MethodGet, svr.URL+"/hello/", http.StatusOK,
		map[string]string{"content-type": "text/html; charset=utf-8"},
		[]string{
			`<h3 id=fizz-buzz>fizz buzz</h3>`,
		},
	)

	testhttp.Expect(t, svr.Client(), "notfound",
		http.MethodGet, svr.URL+"/notfound/", http.StatusNotFound,
		nil,
		nil,
	)

	testhttp.Expect(t, svr.Client(), "blogindex",
		http.MethodGet, svr.URL+"/blog/", http.StatusOK,
		map[string]string{"content-type": "text/html; charset=utf-8"},
		[]string{
			`<title>blog | seankhliao</title>`,
			`<li><time datetime=2021-01-02>12021-01-02</time> | <a href=/blog/12021-01-02-xxx/>xxx</a>`,
		},
	)

	testhttp.Expect(t, svr.Client(), "sitemap",
		http.MethodGet, svr.URL+"/sitemap.txt", http.StatusOK,
		map[string]string{"content-type": "text/plain; charset=utf-8"},
		[]string{
			`https://example.com/blog/12021-01-02-xxx/`,
			`https://example.com/hello/`,
		},
	)
}
