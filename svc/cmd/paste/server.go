package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-logr/logr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/content"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/internal/o11y"
)

type Server struct {
	host string

	storeURL    string
	storePrefix string
	store       *clientv3.Client

	l        logr.Logger
	t        trace.Tracer
	sizeHist metric.Int64Histogram

	startTime time.Time
	pastePage []byte
	indexPage []byte
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	flags.StringVar(&s.storeURL, "store.url", "http://etcd-0.etcd:2379", "etcd url")
	flags.StringVar(&s.storePrefix, "store.prefix", "paste", "key prefix in etcd")
	flag.StringVar(&s.host, "paste.host", "paste.seankhliao.com", "canonical hostname")

	return &s
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("paste")
	s.t = t.Tracer("paste")

	var err error
	s.sizeHist, err = m.Meter("paste").NewInt64Histogram("paste.size")
	if err != nil {
		return err
	}

	s.store, err = clientv3.NewFromURL(s.storeURL)
	if err != nil {
		return err
	}

	pasteRaw, err := fs.ReadFile(content.Paste, "paste.md")
	if err != nil {
		return err
	}
	s.pastePage, err = render.CompactBytes("paste", "upload", fmt.Sprintf("https://%s/paste/", s.host), pasteRaw)
	if err != nil {
		return err
	}

	indexRaw, err := fs.ReadFile(content.Paste, "index.md")
	if err != nil {
		return err
	}
	s.indexPage, err = render.CompactBytes("paste", "simple paste host", fmt.Sprintf("https://%s/", s.host), indexRaw)
	if err != nil {
		return err
	}

	s.startTime = time.Now()

	mux.HandleFunc("/p/", s.lookupHandler)
	mux.HandleFunc("/paste/", s.pasteHandler)
	mux.HandleFunc("/", s.indexHandler)

	return nil
}

func (s *Server) indexHandler(rw http.ResponseWriter, r *http.Request) {
	_, span, _ := o11y.Start(s.t, s.l, r.Context(), "index")
	defer span.End()

	if r.URL.Path != "/" {
		http.Redirect(rw, r, "/", http.StatusFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(rw, "GET only", http.StatusMethodNotAllowed)
		return
	}

	http.ServeContent(rw, r, "index.html", s.startTime, bytes.NewReader(s.indexPage))
}

func (s *Server) pasteHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "paste")
	defer span.End()

	if r.URL.Path != "/paste/" {
		http.Redirect(rw, r, "/paste/", http.StatusFound)
		return
	}

	if r.Method == http.MethodGet {
		http.ServeContent(rw, r, "index.html", s.startTime, bytes.NewReader(s.pastePage))
		return
	} else if r.Method != http.MethodPost {
		http.Error(rw, "not implemented", http.StatusMethodNotAllowed)
		return
	}

	val := r.FormValue("paste")
	if val == "" {
		err := r.ParseMultipartForm(1 << 22) // 4M
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusBadRequest, err, "bad multipart form")
			return
		}
		mpf, _, err := r.FormFile("upload")
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusBadRequest, err, "bad form file")
			return
		}
		defer mpf.Close()
		var buf strings.Builder
		n, err := io.Copy(&buf, mpf)
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "copy file")
			return
		}
		val = buf.String()

		s.sizeHist.Record(ctx, n)
	}

	sum := sha256.Sum256([]byte(val))
	sum2 := base64.URLEncoding.EncodeToString(sum[:])

	basekey := path.Join(s.storePrefix, "p")
	key := ""

	for le := 7; le < 21; le++ {
		key = path.Join(basekey, sum2[:le])
		res, err := s.store.Get(ctx, key)
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "precheck key")
			return
		}
		if len(res.Kvs) == 0 {
			break
		}
		sum3 := base64.URLEncoding.EncodeToString(res.Kvs[0].Value)
		if sum2 == sum3 {
			fmt.Fprintf(rw, "https://paste.seankhliao.com%s", strings.TrimPrefix(key, s.storePrefix))
			return
		}

		if le == 20 {
			o11y.HttpError(rw, l, span, http.StatusInsufficientStorage, errors.New("no unique key"), "max key len reached", attribute.String("key", key))
			return
		}
	}

	_, err := s.store.Put(ctx, key, val)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInsufficientStorage, err, "store", attribute.String("key", key))
		return
	}

	s.l.Info("stored", "key", key, "size", len(val))

	fmt.Fprintf(rw, "https://paste.seankhliao.com%s", strings.TrimPrefix(key, s.storePrefix))
}

func (s *Server) lookupHandler(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "lookup")
	defer span.End()

	if r.Method != http.MethodGet {
		o11y.HttpError(rw, l, span, http.StatusMethodNotAllowed, nil, "use GET")
		return
	}
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) != 3 { // /p/$id
		o11y.HttpError(rw, l, span, http.StatusBadRequest, nil, "unknown path format")
		return
	}

	key := path.Join(s.storePrefix, r.URL.Path)
	s.l.Info("get", "key", key)
	res, err := s.store.Get(ctx, key)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "get from store")
		return
	}

	if len(res.Kvs) != 1 {
		var err error
		if len(res.Kvs) > 1 {
			err = errors.New("extra pastes")
		}
		http.Error(rw, "not found", http.StatusNotFound)
		o11y.HttpError(rw, l, span, http.StatusNotFound, err, "not found", attribute.Int("pastes_found", len(res.Kvs)))
		return
	}

	http.ServeContent(rw, r, "", time.Unix(res.Kvs[0].Version, 0), bytes.NewReader(res.Kvs[0].Value))
}
