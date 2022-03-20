package ghdefaults

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v41/github"
	"github.com/rs/zerolog"
	"go.seankhliao.com/mono/internal/envconf"
	"go.seankhliao.com/mono/internal/httpsvc"
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

var _ httpsvc.HTTPSvc = &Server{}

type Server struct {
	log zerolog.Logger

	appID      int64
	appKey     []byte
	hookSecret []byte
}

func (s *Server) Init(log zerolog.Logger) error {
	s.log = log

	var err error
	appIDStr := envconf.String("GH_APP_ID", "0")
	s.appID, err = strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("parse GH_APP_ID=%q: %w", appIDStr, err)
	}
	appKeyFile := envconf.String("GH_APP_KEY_FILE", "")
	s.appKey, err = os.ReadFile(appKeyFile)
	if err != nil {
		return fmt.Errorf("read GH_APP_KEY_FILE=%q: %w", appKeyFile, err)
	}
	s.appKey = bytes.TrimSpace(s.appKey)
	hookSecretFile := envconf.String("GH_WEBHOOK_SECRET_FILE", "")
	s.hookSecret, err = os.ReadFile(hookSecretFile)
	if err != nil {
		return fmt.Errorf("read GH_WEBHOOK_SECRET_FILE=%q: %w", hookSecretFile, err)
	}
	s.hookSecret = bytes.TrimSpace(s.hookSecret)

	return nil
}

func (s Server) Desc() string {
	return "sets defaults on newly created github repos"
}

func (s Server) Help() string {
	return `
GH_APP_ID
        github app id
GH_APP_KEY_FILE
        path to github secret key (rsa pem)
GH_WEBHOOK_SECRET_FILE
        path to github shared webhook secret
`
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	event, eventType, err := s.getPayload(ctx, r)
	if err != nil {
		s.log.Info().Err(err).Msg("invalid payload")
		http.Error(rw, "invalid payload", http.StatusBadRequest)
		return
	}

	switch event := event.(type) {
	case *github.InstallationEvent:
		s.installEvent(ctx, event)
	case *github.RepositoryEvent:
		s.repoEvent(ctx, rw, event)
	default:
		s.log.Debug().Str("event", eventType).Msg("ignoring event")
	}
}

func (s *Server) getPayload(ctx context.Context, r *http.Request) (any, string, error) {
	payload, err := github.ValidatePayload(r, s.hookSecret)
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

func (s *Server) installEvent(ctx context.Context, event *github.InstallationEvent) {
	owner := *event.Installation.Account.Login
	installID := *event.Installation.ID
	switch *event.Action {
	case "created":
		if _, ok := defaultConfig[owner]; !ok {
			return
		}

		s.log.Info().Int("repos", len(event.Repositories)).Msg("setting defaults on repos")
		go func() {
			for _, repo := range event.Repositories {
				err := s.setDefaults(ctx, installID, owner, *repo.Name)
				if err != nil {
					s.log.Err(err).Str("repo", owner+"/"+*repo.Name).Msg("setting defaults")
					continue
				}
				s.log.Info().Str("repo", owner+"/"+*repo.Name).Msg("set defaults")
			}
		}()
	default:
		s.log.Debug().Str("action", *event.Action).Msg("ignoring installation action")
	}
}

func (s *Server) repoEvent(ctx context.Context, rw http.ResponseWriter, event *github.RepositoryEvent) {
	installID := *event.Installation.ID
	owner := *event.Repo.Owner.Login
	repo := *event.Repo.Name
	switch *event.Action {
	case "created", "transferred":
		if _, ok := defaultConfig[owner]; !ok {
			return
		}
		err := s.setDefaults(ctx, installID, owner, repo)
		if err != nil {
			s.log.Err(err).Str("repo", owner+"/"+repo).Msg("setting defaults")
			http.Error(rw, "set defaults", http.StatusInternalServerError)
			return
		}
		s.log.Info().Str("repo", owner+"/"+repo).Msg("set defaults")
	default:
		s.log.Debug().Str("action", *event.Action).Msg("ignoring repo event")
	}
}

func (s *Server) setDefaults(ctx context.Context, installID int64, owner, repo string) error {
	config := defaultConfig[owner]
	tr := http.DefaultTransport
	tr, err := ghinstallation.NewAppsTransport(tr, s.appID, s.appKey)
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
