package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

func (s *Server) checkAllowlist(host, path string, header map[string]string) bool {
	for _, re := range s.allow[host] {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func (s *Server) checkTokens(host string, headers map[string]string) string {
	token := headers["authorization"]
	token = strings.TrimPrefix(token, "Bearer ")

	tokens, ok := s.tokens[host]
	if !ok {
		return ""
	}
	return tokens[token]
}

func (s *Server) checkBasic(headers map[string]string) string {
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

func (s *Server) compareHtpasswd(user, pass []byte) error {
	hashed, ok := s.passwds[string(user)]
	if !ok {
		return errNotRegistered
	}
	return bcrypt.CompareHashAndPassword(hashed, pass)
}

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
