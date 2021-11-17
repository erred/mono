package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.seankhliao.com/mono/internal/render"
	"go.seankhliao.com/mono/internal/webserver"
)

func main() {
	var ro render.Options
	var fn string
	flag.StringVar(&fn, "file", "index.md", "file to serve")
	flag.StringVar(&ro.Data.GTMID, "gtm", "", "Google Tag Manager ID for analytics")
	flag.StringVar(&ro.Data.URLCanonical, "canonical", "https://arch.seankhliao.com", "canonical base url")
	flag.BoolVar(&ro.Data.Compact, "compact", true, "compact header")
	flag.BoolVar(&ro.MarkdownSkip, "raw", false, "skip markdown processing")
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx, l := webserver.BaseContext()

	var err error
	wo.Handler, err = newHttp(&ro, fn)
	if err != nil {
		l.Error(err, "setup")
		os.Exit(1)
	}

	webserver.Run(ctx, wo)
}

func newHttp(ro *render.Options, fn string) (http.Handler, error) {
	fin, err := os.Open(fn)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", fn, err)
	}

	buf := &bytes.Buffer{}
	err = render.Render(ro, buf, fin)
	if err != nil {
		return nil, fmt.Errorf("render: %w", err)
	}
	b := buf.Bytes()

	tracer := otel.Tracer("singlepage")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := tracer.Start(ctx, "handle")
		defer span.End()
		l := logr.FromContextOrDiscard(ctx).WithName("singlepage")
		l = l.WithValues("method", r.Method, "path", r.URL.Path, "user-agent", r.UserAgent())

		if r.URL.Path != "/" {
			l.Info("redirected")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		w.Header().Set("eicar-bane-of-antivirus", `X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
		w.Write(b)
		l.Info("served")
	})
	return mux, nil
}
