package authnf

import (
	"bytes"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"go.seankhliao.com/mono/auth/authnbpb"
	"go.seankhliao.com/mono/internal/o11y"
	"go.seankhliao.com/mono/internal/web/render"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleIndex(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "index")
	defer span.End()

	data := indexData{
		ID:       r.FormValue("id"),
		Message:  r.FormValue("msg"),
		Redirect: r.FormValue("redirect"),
	}

	sessionToken := extractSession(r, s.cookieName)
	if sessionToken != "" {
		sRes, err := s.authnb.GetSession(ctx, &authnbpb.GetSessionRequest{
			SessionToken: sessionToken,
		})
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "get user auth")
			return
		}
		data.ID = sRes.UserId
		data.LoggedIn = data.ID != "" && data.ID != "anonymous"
	}

	if data.LoggedIn && data.Redirect != "" {
		http.Redirect(rw, r, data.Redirect, http.StatusFound)
		return
	}

	var buf bytes.Buffer
	err := loginTmpl.Execute(&buf, data)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "prerender")
		return
	}
	err = render.Compact(
		rw,
		`authnf`,
		`sso for `+s.cookieDomain,
		`https://authnf.`+s.cookieDomain+`/`,
		buf.Bytes(),
	)
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "render")
		return
	}
}

type indexData struct {
	LoggedIn bool
	ID       string
	Message  string
	Redirect string
}

var loginTmpl = template.Must(template.New("").Parse(`
### _authn_

Single sign on for seankhliao.com

{{ if .LoggedIn}}

#### _hello_ {{ .ID }}

Long time no see!

<form action="/logout" method="post">
  <input type="submit" value="Logout">
</form>

{{ else }}

#### _login_

{{ if .Message }}
_message:_ {{ .Message }}
{{ end }}

<form action="/login" method="post">
  <label>Email: <input name="email" type="email" value="{{ .ID }}"></label>
  <label>Password: <input name="password" type="password"></label>
  <input name="redirect" type="hidden" value="{{ .Redirect }}">
  <input type="submit" value="login">
</form>

{{ end }}
        `))

func (s *Server) handleLogin(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "logout")
	defer span.End()

	id, pass := r.FormValue("email"), r.FormValue("password")
	uaRes, err := s.authnb.GetUserAuth(ctx, &authnbpb.GetUserAuthRequest{
		UserId: id,
	})
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "get user auth")
		return
	}

	err = bcrypt.CompareHashAndPassword(uaRes.Bcrypt, []byte(pass))
	if err != nil {
		l.Info("invalid credentials", "id", id, "err", err)
		http.Redirect(rw, r, "/?"+url.Values{"msg": {"invalid credentials"}, "id": {id}}.Encode(), http.StatusFound)
		return
	}

	sRes, err := s.authnb.CreateSession(ctx, &authnbpb.CreateSessionRequest{
		UserId: id,
		Ttl:    int64(s.cookieTTL.Seconds()),
	})
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "create session")
	}

	http.SetCookie(rw, &http.Cookie{
		Domain:   s.cookieDomain,
		Expires:  time.Now().Add(s.cookieTTL),
		HttpOnly: true,
		Name:     s.cookieName,
		Path:     "/",
		Value:    sRes.SessionToken,
	})

	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/"
	}

	http.Redirect(rw, r, redirect, http.StatusFound)
}

func (s *Server) handleLogout(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "logout")
	defer span.End()

	sessionToken := extractSession(r, s.cookieName)
	if sessionToken != "" {
		_, err := s.authnb.DeleteSession(ctx, &authnbpb.DeleteSessionRequest{
			SessionToken: sessionToken,
		})
		if err != nil {
			l.Error(err, "backend session delete")
		}
	}

	http.SetCookie(rw, &http.Cookie{
		Domain:   s.cookieDomain,
		HttpOnly: true,
		Name:     s.cookieName,
		Path:     "/",
		Value:    "",
		Expires:  time.Unix(0, 0),
	})

	l.Info("logout")
	http.Redirect(rw, r, "/", http.StatusOK)
}

func extractSession(r *http.Request, cookieName string) string {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return c.Value
}
