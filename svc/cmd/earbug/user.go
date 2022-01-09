package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"text/template"

	"go.opentelemetry.io/otel/attribute"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/internal/o11y"
)

// handleUser shows the user data page
func (s *Server) handleIndex(rw http.ResponseWriter, r *http.Request) {
	_, span, l := o11y.Start(s.t, s.l, r.Context(), "index")
	defer span.End()

	data := indexData{
		ID:       r.Header.Get("auth-id"),
		LoginURL: "https://authn.seankhliao.com/?" + url.Values{"redirect": {s.CanonicalURL + "/"}}.Encode(),
	}

	if data.ID == "anonymous" {
		data.ID = ""
	}

	if data.ID != "" {
		s.pollWorkerMu.Lock()
		_, data.Poller = s.pollWorkerMap[data.ID]
		s.pollWorkerMu.Unlock()

		if !data.Poller {
			u, err := url.Parse(s.Auth.AuthURL(data.ID))
			if err != nil {
				o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "generate spotify auth url")
				return
			}
			data.Query = u.Query()
			u.RawQuery = ""
			data.AuthPath = u.String()
		}

	}

	var buf bytes.Buffer
	err := indexTmpl.Execute(&buf, data)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "prerender page")
		return
	}

	l.Info("serving user", "id", data.ID)

	err = render.Compact(
		rw,
		"index",
		"spotify history log",
		s.CanonicalURL+"/",
		buf.Bytes(),
	)
	if err != nil {
		l.Error(err, "render")
	}
}

// handleAuthCallback handles the spotify authorization callback
// by storing the token and starting a poll worker
func (s *Server) handleAuthCallback(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "spotify_callback")
	defer span.End()

	user := r.FormValue("state")
	l = l.WithValues("user", user)

	token, err := s.Auth.Token(ctx, user, r)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "spotify token exchange")
		return
	}

	b, err := json.Marshal(token)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "marshal token")
		return
	}

	p := path.Join(s.StorePrefix, "token", user)
	_, err = s.Store.Put(ctx, p, string(b))
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "store token", attribute.String("store_path", p))
		return
	}

	l.Info("authorized")

	s.addPollWorker(ctx, user, token)
	http.Redirect(rw, r, "/user", http.StatusFound)
}

type indexData struct {
	ID string

	LoginURL string

	Poller   bool
	AuthPath string
	Query    url.Values
}

var indexTmpl = template.Must(template.New("").Parse(`
### _earbug_

A simple spotify listening history logger.

#### _user_ info
{{ if not .ID }}

[login to earbug]({{ .LoginURL }}) to see your info

{{ else }}

Hello {{ .ID }}

[view](/user/history) _history_

{{ if .Poller }}
Background worker is _running_
{{ else }}
Background worker is _not_ running

<form action="{{ .AuthPath }}" method="get">
{{ range $k,$v := .Query }}
<input type="hidden" name="{{ $k }}" value="{{ index $v 0 }}">
{{ end }}
<input type="submit" value="Authorize earbug with Spotify">
</form>
{{ end }}

{{ end }}
`))
