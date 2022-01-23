package runsvr

import (
	"context"
	"runtime"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func otelGRPC(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	var dos []grpc.DialOption

	dos = append(dos, grpc.WithInsecure())

	conn, err := grpc.DialContext(ctx, endpoint, dos...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func otelMetrics(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) error {
	exporter, err := otlpmetric.New(ctx,
		otlpmetricgrpc.NewClient(otlpmetricgrpc.WithGRPCConn(conn)),
	)
	if err != nil {
		return err
	}

	// provider
	provider := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(10*time.Second),
		controller.WithResource(res),
	)
	provider.Start(ctx)

	// global
	global.SetMeterProvider(provider)

	return nil
}

func otelResources(ctx context.Context) (*resource.Resource, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return resource.Default(), nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(bi.Path),
			semconv.ServiceVersionKey.String(bi.Main.Version),
			semconv.ProcessRuntimeVersionKey.String(runtime.Version()),
		),
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func otelTracer(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) error {
	// Exporter
	exporter, err := otlptrace.New(ctx,
		otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn)),
	)
	if err != nil {
		return err
	}

	// provider
	provider := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter),
	)

	// global
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
		jaeger.Jaeger{},
	))

	return nil
}
