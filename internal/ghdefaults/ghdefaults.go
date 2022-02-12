package ghdefaults

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/monolith/component"
	"go.seankhliao.com/mono/monolith/o11y"
	_ "gocloud.dev/blob/fileblob"
	"golang.org/x/oauth2"
)

var (
	_ component.Component = &Component{}

	defaultConfig = map[string]github.Repository{
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
)

type Component struct {
	enabled   bool
	name      string
	hostname  string
	bucketURL string

	t o11y.Tool

	webhookSecretPath string
	webhookSecret     []byte
	githubSecretPath  string
	githubSecret      []byte
	appID             int64
}

func New(name string) *Component {
	return &Component{
		name: name,
	}
}

func (c *Component) Enabled() bool { return c.enabled }

func (c *Component) Register(flags *flag.FlagSet) {
	flags.BoolVar(&c.enabled, c.name, true, "enable component")
	flags.StringVar(&c.hostname, c.name+".host", fmt.Sprintf("%s.seankhliao.com", c.name), "hostname for serving")
	flags.StringVar(&c.webhookSecretPath, c.name+".webhook", "/var/run/secrets/github/webhook-secret", "path to file containing webhook secret")
	flags.StringVar(&c.githubSecretPath, c.name+".secret", "/var/run/secrets/github/github.pem", "path to file containing github app private key")
	flags.Int64Var(&c.appID, c.name+".appid", 0, "app id in github")
}

func (c *Component) HTTP(ctx context.Context, tp o11y.ToolProvider, mux *http.ServeMux) {
	c.t = tp.Tool(c.name)

	var err error
	c.webhookSecret, err = os.ReadFile(c.webhookSecretPath)
	if err != nil {
		//
	}
	c.webhookSecret = bytes.TrimSpace(c.webhookSecret)

	c.githubSecret, err = os.ReadFile(c.githubSecretPath)
	if err != nil {
		//
	}
	c.githubSecret = bytes.TrimSpace(c.githubSecret)

	mux.HandleFunc(c.hostname+"/", c.handleIndex)
	mux.HandleFunc(c.hostname+"/webhook", c.handleWebhook)
}

func (c *Component) handleWebhook(rw http.ResponseWriter, r *http.Request) {
	ctx, span := c.t.Start(r.Context(), "webhook")
	defer span.End()

	event, eventType, err := c.getPayload(ctx, r)
	if err != nil {
		http.Error(rw, "get payload", http.StatusUnauthorized)
		return
	}

	switch event := event.(type) {
	case *github.InstallationEvent:
		c.installEvent(ctx, event)
	case *github.RepositoryEvent:
		c.repoEvent(ctx, rw, event)
	default:
		c.t.Info("ignoring event", "event", eventType)
	}
}

func (c *Component) installEvent(ctx context.Context, event *github.InstallationEvent) {
	_, span := c.t.Start(ctx, "install")
	defer span.End()

	owner := *event.Installation.Account.Login
	installID := *event.Installation.ID
	switch *event.Action {
	case "created":
		if _, ok := defaultConfig[owner]; !ok {
			return
		}

		go func() {
			ctx := context.TODO()                  // start new root to avoid getting cancelled
			ctx = trace.ContextWithSpan(ctx, span) // use span from outside
			ctx, span := c.t.Start(ctx, "install-repos", trace.WithAttributes(
				attribute.String("owner", owner),
			))

			defer span.End()
			for _, repo := range event.Repositories {
				err := c.setDefaults(ctx, installID, owner, *repo.Name)
				if err != nil {
					continue
				}
			}
		}()
	default:
	}
}

func (c *Component) repoEvent(ctx context.Context, rw http.ResponseWriter, event *github.RepositoryEvent) {
	ctx, span := c.t.Start(ctx, "repository_event")
	defer span.End()

	installID := *event.Installation.ID
	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name
	switch *event.Action {
	case "created", "transferred":
		if _, ok := defaultConfig[owner]; !ok {
			return
		}
		err := c.setDefaults(ctx, installID, owner, repo)
		if err != nil {
			http.Error(rw, "set defaults", http.StatusInternalServerError)
			return
		}
	default:
	}
}

func (c *Component) getPayload(ctx context.Context, r *http.Request) (interface{}, string, error) {
	_, span := c.t.Start(ctx, "get payload")
	defer span.End()

	payload, err := github.ValidatePayload(r, c.webhookSecret)
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

func (c *Component) handleIndex(rw http.ResponseWriter, r *http.Request) {
	_, span := c.t.Start(r.Context(), "dispatch")
	defer span.End()

	http.Redirect(rw, r, "https://seankhliao.com", http.StatusFound)
}

func (c *Component) setDefaults(ctx context.Context, installID int64, owner, repo string) error {
	ctx, span := c.t.Start(ctx, "set-defaults", trace.WithAttributes(
		attribute.String("repo", owner+"/"+"repo"),
		attribute.Int64("installation-id", installID),
	))
	defer span.End()

	config := defaultConfig[owner]
	tr := http.DefaultTransport
	tr, err := ghinstallation.NewAppsTransport(tr, c.appID, c.githubSecret)
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
