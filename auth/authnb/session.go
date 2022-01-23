package authnb

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"
	authnbv1 "go.seankhliao.com/mono/apis/authnb/v1"
	"go.seankhliao.com/mono/internal/o11y"
)

const (
	sessionKeyFmt = "authnb/sessions/%s"
)

func (s *Server) GetSession(ctx context.Context, r *authnbv1.GetSessionRequest) (*authnbv1.GetSessionResponse, error) {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "GetSession")
	defer span.End()

	key := fmt.Sprintf(sessionKeyFmt, r.SessionToken)
	result, err := s.store.Get(ctx, key)
	if err != nil {
		o11y.Error(l, span, err, "store get session")
		return nil, err
	}

	var response authnbv1.GetSessionResponse
	switch len(result.Kvs) {
	case 0: // not found
	case 1: // found
		response.UserId = string(result.Kvs[0].Value)
	default: // ???
		o11y.Error(l, span, err, "got extra sessions", attribute.Int("results", len(result.Kvs)))
		return &response, nil
	}

	o11y.OK(l, span, "got session")
	return &response, nil
}

func (s *Server) CreateSession(ctx context.Context, r *authnbv1.CreateSessionRequest) (*authnbv1.CreateSessionResponse, error) {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "CreateSession")
	defer span.End()

	buf := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		o11y.Error(l, span, err, "read from rand")
		return nil, err
	}
	sessionToken := base64.URLEncoding.EncodeToString(buf)

	leaseRes, err := s.store.Grant(ctx, r.Ttl)
	if err != nil {
		o11y.Error(l, span, err, "store grant lease")
	}

	key := fmt.Sprintf(sessionKeyFmt, sessionToken)
	_, err = s.store.Put(ctx, key, r.UserId, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		o11y.Error(l, span, err, "store put session")
		return nil, err
	}

	o11y.OK(l, span, "created session")
	return &authnbv1.CreateSessionResponse{
		SessionToken: sessionToken,
	}, nil
}

func (s *Server) DeleteSession(ctx context.Context, r *authnbv1.DeleteSessionRequest) (*authnbv1.DeleteSessionResponse, error) {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "DeleteSession")
	defer span.End()

	key := fmt.Sprintf(sessionKeyFmt, r.SessionToken)
	result, err := s.store.Delete(ctx, key)
	if err != nil {
		o11y.Error(l, span, err, "store delete session")
		return nil, err
	}

	o11y.OK(l, span, "deleted session")
	return &authnbv1.DeleteSessionResponse{
		Found: result.Deleted > 0,
	}, nil
}
