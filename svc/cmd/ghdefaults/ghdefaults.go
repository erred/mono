package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-logr/logr"
	"github.com/google/go-github/v38/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var defaultConfig = map[string]github.Repository{
	"erred": {
		AllowMergeCommit:    github.Bool(false),
		AllowSquashMerge:    github.Bool(true),
		AllowRebaseMerge:    github.Bool(false),
		DeleteBranchOnMerge: github.Bool(true),
		HasIssues:           github.Bool(false),
		HasWiki:             github.Bool(false),
		HasPages:            github.Bool(false),
		HasProjects:         github.Bool(false),
		IsTemplate:          github.Bool(false),
		Archived:            github.Bool(true),
	},
	"seankhliao": {
		AllowMergeCommit:    github.Bool(false),
		AllowSquashMerge:    github.Bool(true),
		AllowRebaseMerge:    github.Bool(false),
		DeleteBranchOnMerge: github.Bool(true),
		HasIssues:           github.Bool(false),
		HasWiki:             github.Bool(false),
		HasPages:            github.Bool(false),
		HasProjects:         github.Bool(false),
		IsTemplate:          github.Bool(false),
	},
}

func New(flags *flag.FlagSet) *Server {
	var s Server
	flags.StringVar(&s.WebhookSecretPath, "github.secret", "/var/run/secrets/github/webhook-secret", "path to file containing webhook secret")
	flags.StringVar(&s.PrivateKeyPath, "github.key", "/var/run/secrets/github/github.pem", "path to file containing github app private key")
	flags.Int64Var(&s.AppID, "github.id", 0, "app id in github")
	return &s
}

type Server struct {
	WebhookSecretPath string
	PrivateKeyPath    string
	AppID             int64

	webhookSecret []byte
	privateKey    []byte
	appID         int64

	tr http.RoundTripper
	t  trace.Tracer
	l  logr.Logger
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("ghdefaults")
	s.t = t.Tracer("ghdefaults")

	var err error
	webhookSecret, err := os.ReadFile(s.WebhookSecretPath)
	if err != nil {
		return fmt.Errorf("webhook-secret %s: %w", s.WebhookSecretPath, err)
	}
	s.webhookSecret = bytes.TrimSpace(webhookSecret)

	s.privateKey, err = os.ReadFile(s.PrivateKeyPath)
	if err != nil {
		return fmt.Errorf("private-key %s: %w", s.PrivateKeyPath, err)
	}

	mux.Handle("/", s)

	s.tr = otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithMeterProvider(m),
		otelhttp.WithTracerProvider(t),
	)

	return nil
}

// ServeHTTP is the main entrypoint and dispatch for different event types
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "dispatch")
	defer span.End()

	switch r.URL.Path {
	case "/webhook":
		_, span = s.t.Start(ctx, "validate-parse")
		payload, err := github.ValidatePayload(r, s.webhookSecret)
		if err != nil {
			span.End()
			s.l.Error(err, "validate webhook")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		eventType := github.WebHookType(r)
		event, err := github.ParseWebHook(eventType, payload)
		if err != nil {
			span.End()
			s.l.Error(err, "parse payload")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		span.End()

		ctx, span = s.t.Start(ctx, "process")
		defer span.End()
		l := s.l.WithValues("event", eventType)
		msg := "processed"
		switch event := event.(type) {
		case *github.InstallationEvent:
			owner := *event.Installation.Account.Login
			l = l.WithValues("action", *event.Action, "owner", owner, "repos", len(event.Repositories))
			installID := *event.Installation.ID
			switch *event.Action {
			case "created":
				if _, ok := defaultConfig[owner]; !ok {
					msg = "ignoring owner"
					break
				}

				go func() {
					ctx := context.TODO()                  // start new root to avoid getting cancelled
					ctx = trace.ContextWithSpan(ctx, span) // use span from outside
					ctx, span := s.t.Start(ctx, "install-repos", trace.WithAttributes(
						attribute.String("owner", owner),
					))

					defer span.End()
					for _, repo := range event.Repositories {
						l := l.WithValues("repo", *repo.Name)
						err := s.setDefaults(ctx, installID, owner, *repo.Name)
						if err != nil {
							l.Error(err, "set defaults on app install")
							continue
						}
						l.Info("processed")
					}
				}()
			default:
				msg = "ignoring action"
			}
		case *github.RepositoryEvent:
			l = l.WithValues("action", *event.Action, "repo", *event.Repo.FullName)
			installID := *event.Installation.ID
			owner := *event.Repo.Owner.Login
			repo := *event.Repo.Name
			switch *event.Action {
			case "created", "transferred":
				if _, ok := defaultConfig[owner]; !ok {
					msg = "ignoring owner"
					break
				}
				err = s.setDefaults(ctx, installID, owner, repo)
				if err != nil {
					l.Error(err, "set defaults on repo install")
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
			default:
				msg = "ignoring action"
			}
		default:
			msg = "ignoring event"
		}
		l.Info(msg)
		w.WriteHeader(http.StatusOK)
	default:
		http.Redirect(w, r, "https://github.com/seankhliao/mono/tree/main/go/cmd/ghdefaults", http.StatusFound)
	}
}

func (s *Server) setDefaults(ctx context.Context, installID int64, owner, repo string) error {
	ctx, span := s.t.Start(ctx, "set-defaults", trace.WithAttributes(
		attribute.String("repo", owner+"/"+"repo"),
		attribute.Int64("installation-id", installID),
	))
	defer span.End()

	config := defaultConfig[owner]
	tr, err := ghinstallation.New(s.tr, s.appID, installID, s.privateKey)
	if err != nil {
		return fmt.Errorf("create ghinstallation transport: %w", err)
	}
	client := github.NewClient(&http.Client{Transport: tr})
	_, _, err = client.Repositories.Edit(ctx, owner, repo, &config)
	if err != nil {
		return fmt.Errorf("update repo settings: %w", err)
	}
	return nil
}
