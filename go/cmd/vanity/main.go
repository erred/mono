package main

import (
	"context"
	_ "embed"
	"flag"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"text/template"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/go/webserver"
	"k8s.io/klog/v2/klogr"
)

const (
	name        = "go.seankhliao.com/vanity"
	redirectURL = "https://seankhliao.com/"
)

func main() {
	wo := webserver.NewOptions(flag.CommandLine)
	flag.Parse()

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	l := klogr.New()

	s := New()

	wo.Logger = l
	wo.Handler = s

	webserver.New(ctx, wo).Run(ctx)
}

//go:embed index.gohtml
var tmplStr string

type server struct {
	// config
	tmpl *template.Template
	t    trace.Tracer
}

func New() *server {
	return &server{
		tmpl: template.Must(template.New("page").Parse(tmplStr)),
		t:    otel.Tracer("vanity"),
	}
}

func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "dispatch")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx)
	// filter paths
	if r.URL.Path == "/" {
		l.Info("redirected")
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	repo := strings.Split(r.URL.Path, "/")[1]
	err := s.tmpl.Execute(w, map[string]string{"Repo": repo})
	if err != nil {
		l.Error(err, "exec template", "path", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	l.Info("served")
}
