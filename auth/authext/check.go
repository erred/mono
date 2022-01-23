package authext

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"go.opentelemetry.io/otel/attribute"
	authnbv1 "go.seankhliao.com/mono/apis/authnb/v1"
	"go.seankhliao.com/mono/internal/o11y"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
)

// Check implements the envoy extensions.filters.http.ext_authz.v3.ExtAuthz API
func (s *Server) Check(ctx context.Context, r *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "check")
	defer span.End()

	h := r.GetAttributes().GetRequest().GetHttp()
	headers, host, path := h.GetHeaders(), h.GetHost(), h.GetPath()
	status, identity, check := "denied", "anonymous", "all"

	l = l.WithValues("host", host, "path", path)
	span.SetAttributes(attribute.String("host", host), attribute.String("path", path))

	defer func() {
		l.Info(status, "check", check)
		span.SetAttributes(
			attribute.String("authorization", status),
			attribute.String("check", check),
			attribute.String("identity", identity),
		)
	}()

	checkRes := make(chan checkStatus, 4)
	// go func() {
	// 	id := s.checkAllowlist(ctx, host, path, headers)
	// 	checkRes <- checkStatus{"allowlist", id}
	// }()
	// go func() {
	// 	id := s.checkTokens(ctx, host, headers)
	// 	checkRes <- checkStatus{"token", id}
	// }()
	// go func() {
	// 	id := s.checkBasic(ctx, headers)
	// 	checkRes <- checkStatus{"basic", id}
	// }()
	go func() {
		id := s.checkSession(ctx, headers)
		checkRes <- checkStatus{"session", id}
	}()

	for i := 0; i < 4; i++ {
		res := <-checkRes
		if res.id != "" {
			status, check, identity = "allowed", res.check, res.id
			return s.okResponse(res.id, headers)
		}
	}

	return s.deniedResponse(h.Scheme + "://" + h.Host + h.Path)
}

type checkStatus struct {
	check string
	id    string
}

func (s *Server) checkSession(ctx context.Context, headers map[string]string) string {
	ctx, span, l := o11y.Start(s.t, s.l, ctx, "checkSession")
	defer span.End()

	rawRequest := fmt.Sprintf("GET / HTTP/1.1\r\nCookie: %s\r\n\r\n", headers["cookie"])
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewBufferString(rawRequest)))
	if err != nil {
		o11y.Error(l, span, err, "parse http")
		return ""
	}
	c, err := req.Cookie(s.cookieName)
	if err != nil {
		o11y.Error(l, span, err, "parse cookie")
		return ""
	}

	res, err := s.authnb.GetSession(ctx, &authnbv1.GetSessionRequest{
		SessionToken: c.Value,
	})
	if err != nil {
		o11y.Error(l, span, err, "store get session")
		return ""
	}

	return res.UserId
}

// okResponse constructs an response allowing the request through,
// setting the `auth-user` header to the resolved identity
func (s *Server) okResponse(id string, headers map[string]string) (*envoy_service_auth_v3.CheckResponse, error) {
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
							Key:   s.headerID,
							Value: id,
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
func (s *Server) deniedResponse(original string) (*envoy_service_auth_v3.CheckResponse, error) {
	return &envoy_service_auth_v3.CheckResponse{
		Status: &status.Status{Code: int32(codes.PermissionDenied)},
		HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{Code: envoy_type_v3.StatusCode_Found},
				Headers: []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "Location",
							Value: s.redirectURL + "?" + url.Values{"redirect": {original}}.Encode(),
						},
					},
				},
			},
		},
	}, nil
}
