package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// Check implements the envoy extensions.filters.http.ext_authz.v3.ExtAuthz API
func (s *Server) Check(ctx context.Context, r *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	ctx, span := s.t.Start(ctx, "check")
	defer span.End()

	h := r.GetAttributes().GetRequest().GetHttp()
	headers, host, path := h.GetHeaders(), h.GetHost(), h.GetPath()
	status, identity, check := "denied", "anonymous", "all"

	l := s.l.WithValues("host", host, "path", path)
	span.SetAttributes(attribute.String("host", host), attribute.String("path", path))
	defer func() {
		l.Info(status, "check", check)
		span.SetAttributes(
			attribute.String("authorization", status),
			attribute.String("check", check),
			attribute.String("identity", identity),
		)
	}()

	ok := s.checkAllowlist(ctx, host, path, headers)
	if ok {
		status, check = "allowed", "allowlist"
		return okResponse("anonymous", headers)
	}

	id := s.checkTokens(ctx, host, headers)
	if id != "" {
		status, identity, check = "allowed", id, "token"
		return okResponse(id, headers)
	}

	user := s.checkBasic(ctx, headers)
	if user != "" {
		status, identity, check = "allowed", user, "basic"
		return okResponse(user, headers)
	}

	return deniedResponse(s.realm)
}

// checkAllowlist checks if a request should be unconditionally allowed based on host/path
func (s *Server) checkAllowlist(ctx context.Context, host, path string, header map[string]string) bool {
	_, span := s.t.Start(ctx, "allowlist")
	defer span.End()

	for _, re := range s.allow[host] {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

// checkTokens checks if a request should be allowed based on a Bearer token
func (s *Server) checkTokens(ctx context.Context, host string, headers map[string]string) string {
	_, span := s.t.Start(ctx, "tokens")
	defer span.End()

	token := headers["authorization"]
	token = strings.TrimPrefix(token, "Bearer ")

	tokens, ok := s.tokens[host]
	if !ok {
		return ""
	}
	return tokens[token]
}

// checkBasic checks if a request should be allowed based on HTTP Basic Auth
func (s *Server) checkBasic(ctx context.Context, headers map[string]string) string {
	_, span := s.t.Start(ctx, "basic")
	defer span.End()

	user, pass := getBasicAuth(headers)
	if len(user) == 0 {
		return ""
	}
	err := s.compareHtpasswd(user, pass)
	if err != nil {
		return ""
	}
	return string(user)
}

// compareHtpasswd checks the password for a given user
func (s *Server) compareHtpasswd(user, pass []byte) error {
	hashed, ok := s.passwds[string(user)]
	if !ok {
		return errNotRegistered
	}
	return bcrypt.CompareHashAndPassword(hashed, pass)
}

// getBasicAuth extracts the user/pass from a header
func getBasicAuth(header map[string]string) (user, pass []byte) {
	v, ok := header["authorization"]
	if !ok {
		return
	}
	prefix := "Basic "
	if !strings.HasPrefix(v, prefix) {
		return
	}
	b, err := base64.StdEncoding.DecodeString(v[(len(prefix)):])
	if err != nil {
		return
	}
	i := bytes.Index(b, []byte{':'})
	if i < 0 {
		return
	}
	return b[:i], b[i+1:]
}

// okResponse constructs an response allowing the request through,
// setting the `auth-user` header to the resolved identity
func okResponse(user string, headers map[string]string) (*envoy_service_auth_v3.CheckResponse, error) {
	var toRemove []string
	for _, h := range []string{"auth-user"} {
		if _, ok := headers[h]; ok {
			toRemove = append(toRemove, h)
		}
	}

	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
			OkResponse: &envoy_service_auth_v3.OkHttpResponse{
				Headers: []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "auth-user",
							Value: user,
						},
					},
				},
				HeadersToRemove: toRemove,
			},
		},
	}, nil
}

// deniedResponse constructs a response denying the request,
// and asks for HTTP Basic Auth
func deniedResponse(realm string) (*envoy_service_auth_v3.CheckResponse, error) {
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{Code: int32(codes.PermissionDenied)},
		HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{Code: envoy_type_v3.StatusCode_Unauthorized},
				Headers: []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "www-authenticate",
							Value: fmt.Sprintf(`Basic realm="%s", charset="UTF-8"`, realm),
						},
					},
				},
			},
		},
	}, nil
}
