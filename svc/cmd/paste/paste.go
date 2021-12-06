package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/go-logr/logr"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.seankhliao.com/mono/content"
	"go.seankhliao.com/mono/internal/web/render"
)

type Server struct {
	storeURL    string
	storePrefix string
	store       *clientv3.Client

	startTime time.Time
	pastePage []byte
	indexPage []byte
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	flags.StringVar(&s.storeURL, "store.url", "http://etcd-0.etcd:2379", "etcd url")
	flags.StringVar(&s.storePrefix, "store.prefix", "paste", "key prefix in etcd")

	return &s
}

func (s *Server) Handler() (http.Handler, error) {
	var err error
	s.store, err = clientv3.NewFromURL(s.storeURL)
	if err != nil {
		return nil, err
	}

	pasteRaw, err := content.Content.Open("paste.seankhliao.com/paste.html")
	if err != nil {
		return nil, err
	}
	defer pasteRaw.Close()
	styleRaw, err := content.Content.ReadFile("paste.seankhliao.com/style.css")
	if err != nil {
		return nil, err
	}
	pasteRo := &render.Options{
		MarkdownSkip: true,
		Data: render.PageData{
			Compact:      true,
			URLCanonical: "https://paste.seankhliao.com/paste",
			Title:        `paste`,
			Description:  `upload`,
			Style:        string(styleRaw),
		},
	}
	var pasteBuf bytes.Buffer
	err = render.Render(pasteRo, &pasteBuf, pasteRaw)
	if err != nil {
		return nil, err
	}
	s.pastePage = pasteBuf.Bytes()

	indexRaw, err := content.Content.Open("paste.seankhliao.com/index.md")
	if err != nil {
		return nil, err
	}
	indexRo := &render.Options{
		Data: render.PageData{
			Compact:      true,
			URLCanonical: "https://paste.seankhliao.com/",
		},
	}
	var indexBuf bytes.Buffer
	err = render.Render(indexRo, &indexBuf, indexRaw)
	if err != nil {
		return nil, err
	}
	s.indexPage = indexBuf.Bytes()

	s.startTime = time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/p/", s.lookupHandler)
	mux.HandleFunc("/paste/", s.pasteHandler)
	mux.HandleFunc("/", s.indexHandler)

	return mux, nil
}

func (s *Server) indexHandler(rw http.ResponseWriter, r *http.Request) {
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
	ctx := r.Context()
	l := logr.FromContextOrDiscard(ctx).WithName("paste")

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
			l.Error(err, "parse multipart form")
			http.Error(rw, "bad parse", http.StatusBadRequest)
			return
		}
		mpf, _, err := r.FormFile("upload")
		if err != nil {
			http.Error(rw, "bad file", http.StatusBadRequest)
			return
		}
		defer mpf.Close()
		var buf strings.Builder
		_, err = io.Copy(&buf, mpf)
		if err != nil {
			l.Error(err, "copy file")
			http.Error(rw, "copy err", http.StatusInternalServerError)
			return
		}
		val = buf.String()
	}

	sum := base64.URLEncoding.EncodeToString([]byte(val))

	basekey := path.Join(s.storePrefix, "p")
	key := ""

	for le := 7; le < 21; le++ {
		key = path.Join(basekey, sum[:le])
		res, err := s.store.Get(ctx, key)
		if err != nil {
			l.Error(err, "precheck get", "key", key)
			http.Error(rw, "etcd get", http.StatusInternalServerError)
			return
		}
		if len(res.Kvs) == 0 {
			break
		}
		sum2 := base64.URLEncoding.EncodeToString(res.Kvs[0].Value)
		if sum2 == sum {
			fmt.Fprintf(rw, "https://paste.seankhliao.com%s", strings.TrimPrefix(key, s.storePrefix))
			return
		}

		if le == 20 {
			l.Error(errors.New("can't find unique key"), "max len reached", "key", key, "sum", sum[:])
			http.Error(rw, "no unique key", http.StatusInsufficientStorage)
			return
		}
	}

	_, err := s.store.Put(ctx, key, val)
	if err != nil {
		l.Error(err, "store", "key", key)
		http.Error(rw, "put", http.StatusInsufficientStorage)
		return
	}

	l.Info("stored", "key", key, "size", len(val))

	fmt.Fprintf(rw, "https://paste.seankhliao.com%s", strings.TrimPrefix(key, s.storePrefix))
}

func (s *Server) lookupHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := logr.FromContextOrDiscard(ctx).WithName("lookup")

	if r.Method != http.MethodGet {
		http.Error(rw, "use GET", http.StatusMethodNotAllowed)
		return
	}
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) != 3 { // /p/$id
		http.Error(rw, "unknown path format", http.StatusBadRequest)
		return
	}

	key := path.Join(s.storePrefix, r.URL.Path)
	l.Info("get", "key", key)
	res, err := s.store.Get(ctx, key)
	if err != nil {
		l.Error(err, "etcd get", "key")
		http.Error(rw, "storage err", http.StatusInternalServerError)
		return
	}

	if len(res.Kvs) != 1 {
		if len(res.Kvs) > 1 {
			l.Error(errors.New("extra pastes"), "more than expected pastes", "pastes", len(res.Kvs), "key", key)
		}
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	http.ServeContent(rw, r, "", time.Unix(res.Kvs[0].Version, 0), bytes.NewReader(res.Kvs[0].Value))
}
