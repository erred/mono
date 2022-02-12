package paste

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"go.opentelemetry.io/otel/metric"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/monolith/component"
	"go.seankhliao.com/mono/monolith/o11y"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
)

var (
	_ component.Component = &Component{}

	//go:embed index.md
	indexRaw []byte

	//go:embed paste.md
	pasteRaw []byte
)

type Component struct {
	enabled   bool
	name      string
	hostname  string
	bucketURL string

	t             o11y.Tool
	uploadCounter metric.Int64Counter
	accessCounter metric.Int64Counter

	bkt *blob.Bucket

	index []byte
	paste []byte
}

func New(name string) *Component {
	return &Component{
		name: name,
	}
}

func (c *Component) Enabled() bool { return c.enabled }

func (c *Component) Register(flags *flag.FlagSet) {
	flags.BoolVar(&c.enabled, c.name, true, "enable component")
	flags.StringVar(&c.hostname, fmt.Sprintf("%s.host", c.name), fmt.Sprintf("%s.seankhliao.com", c.name), "hostname for serving")
	flags.StringVar(&c.bucketURL, fmt.Sprintf("%s.bucket", c.name), fmt.Sprintf("file://%s", c.name), "bucket url")
}

func (c *Component) HTTP(ctx context.Context, tp o11y.ToolProvider, mux *http.ServeMux) {
	c.t = tp.Tool(c.name)
	c.uploadCounter = c.t.NewInt64Counter("paste.upload.count")

	var err error
	c.index, err = render.CompactBytes(
		c.hostname,
		"pastebin",
		fmt.Sprintf("https://%s/", c.hostname),
		indexRaw,
	)
	if err != nil {
		c.t.Error(err, "render content")
		return
	}
	c.paste, err = render.CompactBytes(
		"paste",
		"upload",
		fmt.Sprintf("https://%s/paste/", c.hostname),
		pasteRaw,
	)
	if err != nil {
		c.t.Error(err, "render content")
		return
	}

	c.bkt, err = blob.OpenBucket(ctx, c.bucketURL)
	if err != nil {
		c.t.Error(err, "open bucket", "bucket", c.bucketURL)
		return
	}

	mux.Handle(c.hostname+"/", c)
}

func (c *Component) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	_, span := c.t.Start(r.Context(), "paste dispatch")
	defer span.End()

	segs := strings.Split(r.URL.Path, "/")
	switch len(segs) {
	case 1:
		http.ServeContent(rw, r, "x.html", time.Time{}, bytes.NewReader(c.index))
	case 2:
		if segs[1] != "p" {
			http.Error(rw, "not found", http.StatusNotFound)
			return
		}
		switch r.Method {
		case http.MethodGet:
			http.ServeContent(rw, r, "x.html", time.Time{}, bytes.NewReader(c.paste))
		case http.MethodPost:
			c.upload(rw, r)
		default:
			http.Error(rw, "not implemented", http.StatusMethodNotAllowed)
		}
	case 3:
		c.lookup(rw, r)
	default:
		http.Error(rw, "not found", http.StatusNotFound)
	}
}

func (c *Component) lookup(rw http.ResponseWriter, r *http.Request) {
	ctx, span := c.t.Start(r.Context(), "paste lookup")
	defer span.End()

	relpath := r.URL.Path[1:]
	obj, err := c.bkt.NewReader(ctx, relpath, nil)
	if err != nil {
		http.Error(rw, "read", http.StatusNotFound)
		return
	}
	defer obj.Close()
	rw.Header().Set("content-type", obj.ContentType())
	_, err = io.Copy(rw, obj)
	if err != nil {
		http.Error(rw, "copy", http.StatusInternalServerError)
		return
	}
}

func (c *Component) upload(rw http.ResponseWriter, r *http.Request) {
	ctx, span := c.t.Start(r.Context(), "paste upload")
	defer span.End()

	val := []byte(r.FormValue("paste"))
	if len(val) == 0 {
		err := r.ParseMultipartForm(1 << 22) // 4M
		if err != nil {
			http.Error(rw, "bad multipart form", http.StatusBadRequest)
			return
		}
		mpf, _, err := r.FormFile("upload")
		if err != nil {
			http.Error(rw, "bad multipart form", http.StatusBadRequest)
			return
		}
		defer mpf.Close()
		var buf bytes.Buffer
		_, err = io.Copy(&buf, mpf)
		if err != nil {
			http.Error(rw, "read", http.StatusInternalServerError)
			return
		}
		val = buf.Bytes()
	}

	sum := sha256.Sum256(val)
	sum2 := base64.URLEncoding.EncodeToString(sum[:])

	key := path.Join("p", sum2[:8])
	err := c.bkt.WriteAll(ctx, key, val, nil)
	if err != nil {
		http.Error(rw, "write", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "https://%s/%s\n", c.hostname, key)
}
