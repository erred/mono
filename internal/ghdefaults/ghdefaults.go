package ghdefaults

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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

type Options struct {
	WebhookSecretPath string
	PrivateKeyPath    string
	AppID             int64
}

func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.WebhookSecretPath, "webhook-secret-path", "/var/run/secrets/github/webhook-secret", "path to file containing webhook secret")
	fs.StringVar(&o.PrivateKeyPath, "private-key-path", "/var/run/secrets/github/private-key", "path to file containing github app private key")
	fs.Int64Var(&o.AppID, "app-id", 0, "app id in github")
	return &o
}

type Server struct {
	webhookSecret []byte
	privateKey    []byte
	appID         int64

	tr http.RoundTripper
	t  trace.Tracer
}

func New(ctx context.Context, o *Options) (*Server, error) {
	tracer := otel.Tracer("ghdefaults")
	ctx, span := tracer.Start(ctx, "new")
	defer span.End()

	_, span = tracer.Start(ctx, "webhook-secret")
	webhookSecret, err := os.ReadFile(o.WebhookSecretPath)
	span.End()
	if err != nil {
		return nil, fmt.Errorf("webhook-secret %s: %w", o.WebhookSecretPath, err)
	}
	webhookSecret = bytes.TrimSpace(webhookSecret)

	_, span = tracer.Start(ctx, "private-key")
	privateKey, err := os.ReadFile(o.PrivateKeyPath)
	span.End()
	if err != nil {
		return nil, fmt.Errorf("private-key %s: %w", o.PrivateKeyPath, err)
	}

	if o.AppID == 0 {
		return nil, fmt.Errorf("app-id must be set")
	}

	return &Server{
		webhookSecret: webhookSecret,
		privateKey:    privateKey,
		appID:         o.AppID,
		tr:            otelhttp.NewTransport(http.DefaultTransport),
		t:             tracer,
	}, nil
}

// ServeHTTP is the main entrypoint and dispatch for different event types
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "dispatch")
	defer span.End()
	l := logr.FromContextOrDiscard(ctx).WithName("dispatch")

	switch r.URL.Path {
	case "/webhook":
		_, span = s.t.Start(ctx, "validate-parse")
		payload, err := github.ValidatePayload(r, s.webhookSecret)
		if err != nil {
			span.End()
			l.Error(err, "validate webhook")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		eventType := github.WebHookType(r)
		event, err := github.ParseWebHook(eventType, payload)
		if err != nil {
			span.End()
			l.Error(err, "parse payload")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		span.End()

		ctx, span = s.t.Start(ctx, "process")
		defer span.End()
		l = l.WithValues("event", eventType)
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