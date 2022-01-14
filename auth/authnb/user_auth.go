package authnb

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.seankhliao.com/mono/auth/authnbpb"
	"go.seankhliao.com/mono/internal/o11y"
)

const (
	userAuthKeyFmt = "authnb/auth/%s/%s"
)

func (s *Server) GetUserAuth(ctx context.Context, r *authnbpb.GetUserAuthRequest) (*authnbpb.GetUserAuthResponse, error) {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "GetUserAuth")
	defer span.End()

	key := fmt.Sprintf(userAuthKeyFmt, r.UserId, "bcrypt")
	result, err := s.store.Get(ctx, key)
	if err != nil {
		o11y.Error(l, span, err, "store get auth", attribute.String("kind", "bcrypt"))
		return nil, err
	}

	var response authnbpb.GetUserAuthResponse
	switch len(result.Kvs) {
	case 0: // not found
	case 1: // found
		response.Bcrypt = result.Kvs[0].Value
	default: // ???
		o11y.Error(l, span, err, "got extra results", attribute.Int("results", len(result.Kvs)))
		return &response, nil
	}

	o11y.OK(l, span, "got user auth")
	return &response, nil
}
