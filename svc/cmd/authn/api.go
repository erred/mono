package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.seankhliao.com/mono/svc/internal/o11y"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleApiLogin(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "login")
	defer span.End()

	if r.Method != http.MethodPost {
		o11y.HttpError(rw, l, span, http.StatusMethodNotAllowed, nil, "bad login request")
		return
	}

	id, pass := r.FormValue("email"), r.FormValue("password")
	span.SetAttributes(attribute.String("id", id))

	err := s.checkCredentials(ctx, id, pass)
	if err != nil {
		l.Info("invalid credentials", "id", id, "err", err)
		http.Redirect(rw, r, "/?"+url.Values{"msg": {"invalid credentials"}, "id": {id}}.Encode(), http.StatusFound)
		return
	}

	sessionToken, err := s.generateSession(ctx, id)
	if err != nil {
		l.Error(err, "no session token generated")
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

func (s *Server) handleApiLogout(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "logout")
	defer span.End()

	sessionTokenCookie, err := r.Cookie(s.cookieName)
	if err == nil {
		err = s.store.DeleteSession(ctx, sessionTokenCookie.Value)
		if err != nil {
			l.Error(err, "delete session")
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
	ctx, span, _ := o11y.Start(s.t, s.l, ctx, "create_session")
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
	_, span, _ := o11y.Start(s.t, s.l, ctx, "compare_password")
	defer span.End()

	stored, err := s.store.GetUserPassword(ctx, id)
	if err != nil {
		return fmt.Errorf("get stored password: %w", err)
	}

	return bcrypt.CompareHashAndPassword(stored, []byte(pass))
}
