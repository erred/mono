package feedagg

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/gorilla/feeds"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/internal/render"
)

type Options struct {
	ConfigPath string
	Dir        string
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.Dir, "data", "/data", "path to data dir")
	fs.StringVar(&o.ConfigPath, "config", "/etc/feedagg.yaml", "path to config file")
	return &o
}

type server struct {
	feeds map[string]struct{}
	store Storer
	t     trace.Tracer
}

func New(ctx context.Context, o *Options) (*server, error) {
	config, err := newConfig(o.ConfigPath)
	if err != nil {
		return nil, err
	}

	store, err := NewSQLite(ctx, o.Dir, config.Feeds)
	if err != nil {
		return nil, err
	}

	s := &server{
		feeds: make(map[string]struct{}),
		store: store,
		t:     otel.Tracer("feedagg"),
	}

	type UP struct {
		r time.Duration
		u string
	}
	// url: refresh
	upstreams := make(map[string]UP)

	for feedID, feedConf := range config.Feeds {
		s.feeds[feedID] = struct{}{}
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
		RunFeedUpdater(ctx, id, up.u, up.r, s.store)
	}

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "dispatch")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx).WithName("dispatch")
	l = l.WithValues("method", r.Method, "path", r.URL.Path, "user_agent", r.UserAgent())

	if r.URL.Path == "/" {
		ro := &render.Options{
			MarkdownSkip: true,
			Data: render.PageData{
				Title:       "feedagg",
				Description: "index of aggregated feeds",
				H1:          "feedagg",
				H2:          "index",
				Style:       style,
				Compact:     true,
			},
		}
		_, span = s.t.Start(ctx, "render")
		err := render.Render(ro, w, renderFeeds(s.feeds))
		span.End()
		if err != nil {
			l.Error(err, "render index")
		}
		return
	}

	segs := strings.Split(r.URL.Path[1:], ".")
	if len(segs) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		l.Info("expected 2 segments")
		return
	}
	feedID, format := segs[0], segs[1]

	_, ok := s.feeds[feedID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		l.Info("feed not found", "feed_id", feedID)
		return
	}

	feed, err := s.store.GetFeed(ctx, feedID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		l.Error(err, "get feed", "feed_id", feedID)
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
				URLCanonical: "https://feedagg.seankhliao.com/" + feedID + ".html",
				Title:        feedID,
				Description:  "aggregated " + feedID,
				H1:           "feedagg",
				H2:           feedID,
			},
		}
		_, span = s.t.Start(ctx, "render")
		err := render.Render(ro, w, renderFeed(feed))
		span.End()
		if err != nil {
			l.Error(err, "render feed", "feed_id", feedID)
			return
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		l.Info("unknown format", "format", format)
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
