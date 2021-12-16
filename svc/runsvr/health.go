package runsvr

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type okHandler struct{}

func (okHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func otelHttpFilter(excludePaths ...string) otelhttp.Filter {
	exclude := make(map[string]struct{}, len(excludePaths))
	for _, p := range excludePaths {
		exclude[p] = struct{}{}
	}
	return otelhttp.Filter(func(r *http.Request) bool {
		_, ok := exclude[r.URL.Path]
		return !ok
	})
}

func (okHandler) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (okHandler) Watch(r *grpc_health_v1.HealthCheckRequest, s grpc_health_v1.Health_WatchServer) error {
	s.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
	select {}
}
