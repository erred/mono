package store

import (
	"context"
	"errors"
	"fmt"
	"path"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/trace"
)

const (
	sessionPrefix = "authn/sessions"
	userPrefix    = "authn/users"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	t trace.Tracer
	c *clientv3.Client
}

func New(u string, t trace.TracerProvider) (*Store, error) {
	c, err := clientv3.NewFromURL(u)
	if err != nil {
		return nil, err
	}
	return &Store{
		t.Tracer("authnstore"),
		c,
	}, nil
}

//
// Session Storage
//

// GetSession gets the user id associated with a session
func (s *Store) GetSession(ctx context.Context, sessionID string) (string, error) {
	ctx, span := s.t.Start(ctx, "GetSession")
	defer span.End()

	res, err := s.c.Get(ctx, path.Join(sessionPrefix, sessionID))
	if err != nil {
		return "", fmt.Errorf("get session: %w", err)
	}
	if len(res.Kvs) != 1 {
		return "", ErrNotFound
	}

	return string(res.Kvs[0].Value), nil
}

// AddSession adds a session for a user
func (s *Store) AddSession(ctx context.Context, sessionID, userID string) error {
	ctx, span := s.t.Start(ctx, "AddSession")
	defer span.End()

	_, err := s.c.Put(ctx, path.Join(sessionPrefix, sessionID), userID)
	if err != nil {
		return fmt.Errorf("put session: %w", err)
	}
	return nil
}

// DeleteSession removes a session
func (s *Store) DeleteSession(ctx context.Context, sessionID string) error {
	ctx, span := s.t.Start(ctx, "DeleteSession")
	defer span.End()

	_, err := s.c.Delete(ctx, path.Join(sessionPrefix, sessionID))
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

//
// User Storage
//

// GetUserPassword retrieves the hashed user password
func (s *Store) GetUserPassword(ctx context.Context, userID string) ([]byte, error) {
	ctx, span := s.t.Start(ctx, "GetUserPassword")
	defer span.End()

	p := path.Join(userPrefix, userID, "bcrypt")
	res, err := s.c.Get(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("get bcrypt hash: %w", err)
	}
	if len(res.Kvs) != 1 {
		return nil, ErrNotFound
	}

	return res.Kvs[0].Value, nil
}
