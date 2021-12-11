package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"go.seankhliao.com/mono/internal/web/render"
)

// authPage shows the link a user should follow to authorize earbug
// /auth/user/<user_id>
func (s *Server) authPage(rw http.ResponseWriter, r *http.Request) {
	l := s.l.WithName("auth").WithValues("page", "user")

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[1] != "auth" || parts[2] != "user" {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}
	user := parts[3]

	err := render.Compact(
		rw,
		"earbug auth",
		"allow earbug access to your spotify account",
		s.CanonicalURL+"/auth/user/"+user,
		[]byte(fmt.Sprintf(authPageMsg, s.Auth.AuthURL(user))),
	)
	if err != nil {
		l.Error(err, "render")
		return
	}
}

const authPageMsg = `
### _authorize_ earbug

_earbug_ needs access to yur spotify account to get your listening data.
Grant authorization via[link](%s)
`

// authCallback handles the spotify authorization callback
// by storing the token and starting a poll worker
func (s *Server) authCallback(rw http.ResponseWriter, r *http.Request) {
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
	err = render.Compact(
		rw,
		"earbug authorized",
		"earbug has been successfully authorized",
		s.CanonicalURL+"/auth/callback",
		[]byte(fmt.Sprintf(authCallbackMsg, user)),
	)
	if err != nil {
		l.Error(err, "render")
		return
	}
}

const authCallbackMsg = `
### auth _success_

Welcome! _earbug_ has been authorized for _%s_
`
