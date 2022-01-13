package main

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/internal/o11y"
)

// handleIndex serves either the login page or user info page (if logged in)
func (s *Server) handleIndex(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "index")
	defer span.End()

	id, status := "anonymous", "unauthenticated"

	defer func() {
		l.Info("index", "status", status, "id")
		span.SetAttributes(
			attribute.String("status", status),
			attribute.String("id", id),
		)
	}()

	redirect := r.FormValue("redirect")

	var err error
	id, err = s.checkAuth(ctx, r)
	if err != nil {
		id = r.FormValue("id")
		msg := r.FormValue("msg")
		err = render.Compact(
			rw,
			"login",
			"login to seankhliao.com",
			"https://"+s.canonicalHostname+"/",
			[]byte(fmt.Sprintf(loginStr, msg, id, redirect)),
		)
		if err != nil {
			l.Error(err, "render", "page", "login")
		}
		return
	}

	status = "authenticated"

	if redirect != "" {
		http.Redirect(rw, r, redirect, http.StatusFound)
		return
	}
	err = render.Compact(
		rw,
		"user",
		"user info",
		"https://"+s.canonicalHostname+"/",
		[]byte(fmt.Sprintf(userStr, id)),
	)
	if err != nil {
		l.Error(err, "render", "page", "info")
	}
}

func (s *Server) checkAuth(ctx context.Context, r *http.Request) (string, error) {
	ctx, span, _ := o11y.Start(s.t, s.l, ctx, "check_auth")
	defer span.End()

	sessionTokenCookie, err := r.Cookie(s.cookieName)
	if err != nil {
		return "", fmt.Errorf("no session token found: %w", err)
	}

	id, err := s.store.GetSession(ctx, sessionTokenCookie.Value)
	if err != nil {
		return "", fmt.Errorf("get stored session: %w", err)
	}
	return id, nil
}

var (
	loginStr = `
### _authn_

Single sign on for seankhliao.com

#### _login_

<form action="/login" method="post">
  %s
  <label>Email: <input name="email" type="email" value="%s"></label>
  <label>Password: <input name="password" type="password"></label>
  <input name="redirect" type="hidden" value="%s">
  <input type="submit" value="login">
  </form>
`

	userStr = `
### _Hello_ %s

Long time no see!

<form action="/logout" method="post">
  <input type="submit" value="Logout">
</form>
`
)
