package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	subs chan []chan []byte

	l logr.Logger
	t trace.Tracer
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	return &s
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("reqsink")
	s.t = t.Tracer("reqsink")

	s.subs = make(chan []chan []byte, 1)
	s.subs <- nil

	mux.HandleFunc("/watch", s.watch)
	mux.HandleFunc("/", s.slurp)

	return nil
}

func (s *Server) watch(rw http.ResponseWriter, r *http.Request) {
	s.l.Info("new watcher")

	flusher, ok := rw.(http.Flusher)

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Connection", "keep-alive")

	fmt.Fprintf(rw, "event: %s\n\n", "watching requests")
	if ok {
		flusher.Flush()
	}

	c := make(chan []byte, 1)
	subs := <-s.subs
	subs = append(subs, c)
	s.subs <- subs
	defer func() {
		subs := <-s.subs
		for i, sub := range subs {
			if sub == c {
				subs = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		s.subs <- subs
	}()

watchloop:
	for {
		select {
		case <-r.Context().Done():
			break watchloop
		case reqb := <-c:
			fmt.Fprintf(rw, "event: --------------------\n\n")
			rw.Write(reqb)
			rw.Write([]byte("\n"))
			if ok {
				flusher.Flush()
			}
		}
	}
}

func (s *Server) slurp(rw http.ResponseWriter, r *http.Request) {
	s.l.Info("slurping request", "path", r.URL.Path)

	reqb, err := httputil.DumpRequest(r, true)
	if err != nil {
		s.l.Error(err, "slurp request")
		http.Error(rw, "dump request", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	sc := bufio.NewScanner(bytes.NewReader(reqb))
	for sc.Scan() {
		buf.WriteString("data: ")
		buf.Write(sc.Bytes())
		buf.WriteRune('\n')
	}

	subs := <-s.subs
	for _, sub := range subs {
		sub <- buf.Bytes()
	}
	s.subs <- subs
}
