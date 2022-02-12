package o11y

import (
	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type ToolProvider struct {
	L logr.Logger
	T trace.TracerProvider
	M metric.MeterProvider
}

func (t ToolProvider) Tool(name string) Tool {
	return Tool{
		Logger:    t.L.WithName(name),
		MeterMust: metric.Must(t.M.Meter(name)),
		Tracer:    t.T.Tracer(name),
	}
}

type Tool struct {
	logr.Logger
	metric.MeterMust
	trace.Tracer
}
