package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"text/template"

	"go.seankhliao.com/mono/internal/web/render"
)

// handleUser shows the user data page
func (s *Server) handleUser(rw http.ResponseWriter, r *http.Request) {
	l := s.l.WithName("auth").WithValues("page", "/user")

	id := r.Header.Get("auth-id")
	if id == "" {
		http.Error(rw, "no id", http.StatusUnauthorized)
		return
	}

	s.pollWorkerMu.Lock()
	_, ok := s.pollWorkerMap[id]
	s.pollWorkerMu.Unlock()

	var authPath string
	var query url.Values
	if !ok {
		u, err := url.Parse(s.Auth.AuthURL(id))
		if err != nil {
			l.Error(err, "spotify auth url")
			http.Error(rw, "internal auth setup", http.StatusInternalServerError)
			return
		}
		query = u.Query()
		u.RawQuery = ""
		authPath = u.String()
	}
	ud := userPageData{
		ID:       id,
		Poller:   ok,
		AuthPath: authPath,
		Query:    query,
	}

	var buf bytes.Buffer
	err := userPageTmpl.Execute(&buf, ud)
	if err != nil {
		l.Error(err, "prerender")
		http.Error(rw, "internal render", http.StatusInternalServerError)
		return
	}

	err = render.Compact(
		rw,
		"user",
		id,
		s.CanonicalURL+"/user",
		buf.Bytes(),
	)
	if err != nil {
		l.Error(err, "render page")
	}
}

// handleAuthCallback handles the spotify authorization callback
// by storing the token and starting a poll worker
func (s *Server) handleAuthCallback(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := s.l.WithName("auth").WithValues("page", "callback")

	user := r.FormValue("state")
	l = l.WithValues("user", user)

	token, err := s.Auth.Token(ctx, user, r)
	if err != nil {
		l.Error(err, "exchange token")
		http.Error(rw, "token exchange", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		l.Error(err, "marhsal token")
		http.Error(rw, "marshal token", http.StatusInternalServerError)
		return
	}

	p := path.Join(s.StorePrefix, "token", user)
	_, err = s.Store.Put(ctx, p, string(b))
	if err != nil {
		l.Error(err, "store token", "path", p)
		http.Error(rw, "store token", http.StatusInternalServerError)
		return
	}

	l.Info("authorized")

	s.addPollWorker(ctx, user, token)
	http.Redirect(rw, r, "/user", http.StatusFound)
}

type userPageData struct {
	ID       string
	Poller   bool
	AuthPath string
	Query    url.Values
}

var userPageTmpl = template.Must(template.New("").Parse(`
### _hello_ {{ .ID }}

{{ if .Poller }}
Background Worker is _running_
{{ else }}
<form action="{{ .AuthPath }}" method="get">
{{ range $k,$v := .Query }}
<input type="hidden" name="{{ $k }}" value="{{ index $v 0 }}">
{{ end }}
<input type="submit" value="Authorize earbug">
</form>
{{ end }}

[view](/user/history) _history_
`))
