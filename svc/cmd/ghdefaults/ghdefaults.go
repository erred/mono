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
	"github.com/google/go-github/v41/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/svc/internal/o11y"
	"golang.org/x/oauth2"
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
	flags.Int64Var(&s.appID, "github.id", 0, "app id in github")
	return &s
}

type Server struct {
	WebhookSecretPath string
	PrivateKeyPath    string

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

	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/webhook", s.handleWebhook)

	s.tr = otelhttp.NewTransport(
		http.DefaultTransport,
		otelhttp.WithMeterProvider(m),
		otelhttp.WithTracerProvider(t),
	)

	return nil
}

func (s *Server) handleWebhook(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "webhook")
	defer span.End()

	event, eventType, err := s.getPayload(ctx, r)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusUnauthorized, err, "get payload")
		return
	}

	switch event := event.(type) {
	case *github.InstallationEvent:
		s.installEvent(ctx, event)
	case *github.RepositoryEvent:
		s.repoEvent(ctx, rw, event)
	default:
		l.Info("ignoring event", "event", eventType)
	}
}

func (s *Server) installEvent(ctx context.Context, event *github.InstallationEvent) {
	_, span, l := o11y.Start(s.t, s.l, ctx, "installation_event")
	defer span.End()

	owner := *event.Installation.Account.Login
	l = l.WithValues("action", *event.Action, "owner", owner, "repos", len(event.Repositories))
	installID := *event.Installation.ID
	switch *event.Action {
	case "created":
		if _, ok := defaultConfig[owner]; !ok {
			l.Info("ignoring owner")
			return
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
		l.Info("ignoring action")
	}
}

func (s *Server) repoEvent(ctx context.Context, rw http.ResponseWriter, event *github.RepositoryEvent) {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "repository_event")
	defer span.End()

	l = l.WithValues("action", *event.Action, "repo", *event.Repo.FullName)
	installID := *event.Installation.ID
	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name
	switch *event.Action {
	case "created", "transferred":
		if _, ok := defaultConfig[owner]; !ok {
			l.Info("ignoring owner")
			return
		}
		err := s.setDefaults(ctx, installID, owner, repo)
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "set defaults on repo install")
			return
		}
		l.Info("defaults set")
	default:
		l.Info("ignoring action")
	}
}

func (s *Server) getPayload(ctx context.Context, r *http.Request) (interface{}, string, error) {
	_, span, _ := o11y.Start(s.t, s.l, ctx, "get_payload")
	defer span.End()

	payload, err := github.ValidatePayload(r, s.webhookSecret)
	if err != nil {
		return nil, "", fmt.Errorf("validate: %w", err)
	}
	eventType := github.WebHookType(r)
	event, err := github.ParseWebHook(eventType, payload)
	if err != nil {
		return nil, "", fmt.Errorf("parse: %w", err)
	}
	return event, eventType, nil
}

func (s *Server) handleIndex(rw http.ResponseWriter, r *http.Request) {
	_, span, _ := o11y.Start(s.t, s.l, r.Context(), "dispatch")
	defer span.End()

	http.Redirect(rw, r, "https://github.com/seankhliao/mono/tree/main/svc/cmd/ghdefaults", http.StatusFound)
}

func (s *Server) setDefaults(ctx context.Context, installID int64, owner, repo string) error {
	ctx, span := s.t.Start(ctx, "set-defaults", trace.WithAttributes(
		attribute.String("repo", owner+"/"+"repo"),
		attribute.Int64("installation-id", installID),
	))
	defer span.End()

	config := defaultConfig[owner]
	tr, err := ghinstallation.NewAppsTransport(s.tr, s.appID, s.privateKey)
	if err != nil {
		return fmt.Errorf("create ghinstallation transport: %w", err)
	}
	client := github.NewClient(&http.Client{Transport: tr})
	installToken, _, err := client.Apps.CreateInstallationToken(ctx, installID, nil)
	if err != nil {
		return fmt.Errorf("create installation token: %w", err)
	}

	client = github.NewClient(&http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *installToken.Token}),
		},
	})

	_, _, err = client.Repositories.Edit(ctx, owner, repo, &config)
	if err != nil {
		return fmt.Errorf("update repo settings: %w", err)
	}
	return nil
}
