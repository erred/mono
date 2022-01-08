package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/svc/cmd/authn/store"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	l logr.Logger
	t trace.Tracer

	cookieDomain      string
	cookieName        string
	canonicalHostname string
	hashedPasswordsMu sync.RWMutex
	hashedPasswords   map[string][]byte

	etcdUrl string
	store   *store.Store
}

func New(flags *flag.FlagSet) *Server {
	s := Server{
		hashedPasswords: make(map[string][]byte),
	}

	flags.StringVar(&s.cookieName, "cookie", "__authn_session", "name of cookie")
	flags.StringVar(&s.cookieDomain, "domain", "seankhliao.com", "cookie domain")
	flags.StringVar(&s.canonicalHostname, "hostname", "authn.seankhliao.com", "canonical hostname")
	flags.StringVar(&s.etcdUrl, "etcd", "http://etcd-0.etcd:2379", "etcd session store")

	return &s
}

func (s *Server) RegisterHTTP(ctx context.Context, mux *http.ServeMux, l logr.Logger, m metric.MeterProvider, t trace.TracerProvider, shutdown func()) error {
	s.l = l.WithName("authn")
	s.t = t.Tracer("authn")

	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/login", s.handleLogin)
	mux.HandleFunc("/logout", s.handleLogout)

	var err error
	s.store, err = store.New(s.etcdUrl, t)
	if err != nil {
		return err
	}

	return nil
}

//
// Web handlers
//

func (s *Server) handleIndex(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := s.t.Start(ctx, "/")
	defer span.End()

	status := "unauthenticated"
	id := "anonymous"

	defer func() {
		s.l.Info("index", "status", status, "id", id)
		span.SetAttributes(
			attribute.String("status", status),
			attribute.String("id", id),
		)
	}()

	var err error
	id, err = s.checkAuth(ctx, r)
	if err != nil {
		msg := r.FormValue("msg")
		id = r.FormValue("id")
		redirect := r.FormValue("redirect")
		err = render.Compact(
			rw,
			"login",
			"login to seankhliao.com",
			"https://"+s.canonicalHostname+"/",
			[]byte(fmt.Sprintf(loginStr, msg, id, redirect)),
		)
		if err != nil {
			s.l.Error(err, "render", "page", "login")
		}
		return
	}

	status = "authenticated"
	err = render.Compact(
		rw,
		"user",
		"user info",
		"https://"+s.canonicalHostname+"/",
		[]byte(fmt.Sprintf(userStr, id)),
	)
	if err != nil {
		s.l.Error(err, "render", "page", "info")
	}
}

//
// API Handlers
//

func (s *Server) handleLogin(rw http.ResponseWriter, r *http.Request) {
	ctx, span := s.t.Start(r.Context(), "/login")
	defer span.End()

	if r.Method != http.MethodPost {
		s.l.Info("bad login request", "method", r.Method)
		http.Error(rw, "POST only", http.StatusMethodNotAllowed)
		return
	}

	id, pass := r.FormValue("email"), r.FormValue("password")
	span.SetAttributes(attribute.String("id", id))

	err := s.checkCredentials(ctx, id, pass)
	if err != nil {
		s.l.Info("invalid credentials", "id", id, "err", err)
		http.Redirect(rw, r, "/?"+url.Values{"msg": {"invalid credentials"}, "id": {id}}.Encode(), http.StatusFound)
		return
	}

	sessionToken, err := s.generateSession(ctx, id)
	if err != nil {
		s.l.Error(err, "no session token generated")
		http.Redirect(rw, r, "/?"+url.Values{"msg": {"internal error"}, "id": {id}}.Encode(), http.StatusFound)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Domain:   s.cookieDomain,
		HttpOnly: true,
		Name:     s.cookieName,
		Path:     "/",
		Value:    sessionToken,
	})

	redirect := r.FormValue("redirect")
	if redirect == "" {
		redirect = "/"
	}

	http.Redirect(rw, r, redirect, http.StatusFound)
}

func (s *Server) handleLogout(rw http.ResponseWriter, r *http.Request) {
	ctx, span := s.t.Start(r.Context(), "/logout")
	defer span.End()

	sessionTokenCookie, err := r.Cookie(s.cookieName)
	if err == nil {
		err = s.store.DeleteSession(ctx, sessionTokenCookie.Value)
		if err != nil {
			s.l.Error(err, "delete session")
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

	http.Redirect(rw, r, "/", http.StatusOK)
}

func (s *Server) generateSession(ctx context.Context, userID string) (string, error) {
	ctx, span := s.t.Start(ctx, "CreateSessionToken")
	defer span.End()

	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("read rand: %w", err)
	}
	sessionID := base64.URLEncoding.EncodeToString(buf)

	err = s.store.AddSession(ctx, sessionID, userID)
	if err != nil {
		return "", fmt.Errorf("store session: %w", err)
	}
	return sessionID, nil
}

func (s *Server) checkCredentials(ctx context.Context, id, pass string) error {
	_, span := s.t.Start(ctx, "ComparePassword")
	defer span.End()

	stored, err := s.store.GetUserPassword(ctx, id)
	if err != nil {
		return fmt.Errorf("get stored password: %w", err)
	}

	return bcrypt.CompareHashAndPassword(stored, []byte(pass))
}

func (s *Server) checkAuth(ctx context.Context, r *http.Request) (string, error) {
	ctx, span := s.t.Start(ctx, "check_auth")
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
### _Login_

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

<form action="/logout" method="post">
  <input type="submit" value="Logout">
</form>
`
)
