package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/feeds"
	"go.seankhliao.com/mono/go/render"
)

type App struct {
	l     logr.Logger
	feeds map[string]struct{}
	store Storer
}

func NewApp(c Config, dataDir string) (*App, error) {
	store, err := NewSQLite(context.Background(), dataDir, c.Feeds)
	if err != nil {
		return nil, err
	}

	a := &App{
		feeds: make(map[string]struct{}),
		store: store,
	}

	type UP struct {
		r time.Duration
		u string
	}
	// url: refresh
	upstreams := make(map[string]UP)

	for feedID, feedConf := range c.Feeds {
		a.feeds[feedID] = struct{}{}
		for id, u := range feedConf.URLs {
			up, ok := upstreams[id]
			if ok {
				if up.r.Seconds() < feedConf.refresh.Seconds() {
					continue
				}
			}
			upstreams[id] = UP{feedConf.refresh, u}
		}
	}

	for id, up := range upstreams {
		RunFeedUpdater(id, up.u, up.r, a.store)
	}

	return a, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.URL.Path == "/" {
		ro := &render.Options{
			MarkdownSkip: true,
			Data: render.PageData{
				Title:       "feed-agg",
				Description: "index of aggregated feeds",
				H1:          "feed-agg",
				H2:          "index",
				Style:       style,
				Compact:     true,
			},
		}
		err := render.Render(ro, w, renderFeeds(a.feeds))
		if err != nil {
			a.l.Error(err, "render index")
		}
		return
	}

	segs := strings.Split(r.URL.Path[1:], ".")
	if len(segs) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	feedID, format := segs[0], segs[1]

	_, ok := a.feeds[feedID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	feed, err := a.store.GetFeed(ctx, feedID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch format {
	case "rss":
		feed.WriteAtom(w)
	case "atom":
		feed.WriteAtom(w)
	case "json":
		feed.WriteJSON(w)
	case "html":
		ro := &render.Options{
			MarkdownSkip: true,
			Data: render.PageData{
				Compact:      true,
				URLCanonical: "https://feed-agg.seankhliao.com/" + feedID + ".html",
				Title:        feedID,
				Description:  "aggregated " + feedID,
				H1:           "feed-agg",
				H2:           feedID,
			},
		}
		err := render.Render(ro, w, renderFeed(feed))
		if err != nil {
			a.l.Error(err, "render feed", "feed", feedID)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func renderFeed(feed *feeds.Feed) io.Reader {
	buf := &bytes.Buffer{}
	buf.WriteString(`
  <table>
    <thead>
      <tr>
        <td><em>title</em></td>
        <td>source</td>
        <td>updated</td>
        <td>author</td>
        <td>description</td>
      </tr>
    </thead>
    <tbody>
`)
	for _, it := range feed.Items {
		fmt.Fprintf(buf, `
      <tr>
        <td><a href="%s">%s</a></td>
        <td>%s</td>
        <td><time>%s</time</td>
        <td>%s</td>
        <td><p class="content">%s</p></td>
      </tr>
                `, it.Link.Href, it.Title, it.Source.Href, it.Updated, it.Author.Name, it.Description)
	}
	buf.WriteString(`
    </tbody>
  </table>
`)
	return buf
}

func renderFeeds(feeds map[string]struct{}) io.Reader {
	buf := &bytes.Buffer{}
	buf.WriteString(`
  <table>
    <thead>
      <tr>
        <td><em>feed</em></td>
        <td>atom</td>
        <td>rss</td>
        <td>json</td>
      </tr>
    </thead>
    <tbody>
`)
	var fs []string
	for f := range feeds {
		fs = append(fs, f)
	}
	sort.Strings(fs)
	for _, f := range fs {
		fmt.Fprintf(buf, `
      <tr>
        <td><a href="%[1]s.html">%[1]s</a></td>
        <td><a href="%[1]s.rss">rss</a></td>
        <td><a href="%[1]s.atom">atom</a></td>
        <td><a href="%[1]s.json">json</a></td>
      </tr>
`, f)
	}
	buf.WriteString(`
    </tbody>
  </table>
`)
	return buf
}

var style = `
table {
  margin: 10vh 0;
}
`
