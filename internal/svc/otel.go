package svc

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type otelclient struct {
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *controller.Controller
	done           chan error
}

func newOtelClient(otlpURL string, log zerolog.Logger) (*otelclient, error) {
	ctx := context.TODO()

	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Err(err).Msg("otel error")
	}))

	u, err := url.Parse(otlpURL)
	if err != nil {
		return nil, fmt.Errorf("parse otlp endpoint: %w", err)
	}
	_, short, version := serviceName()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(short),
			semconv.ServiceVersionKey.String(version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create otel resource: %w", err)
	}

	log.Debug().Str("otlp_endpoint", u.Host).Msg("connecting to otlp")
	conn, err := grpc.DialContext(ctx, u.Host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUserAgent(userAgent()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("setup otlp grpc conn to %s: %w", u.Host, err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	// traceExporter, err := stdouttrace.New()
	if err != nil {
		return nil, fmt.Errorf("setup trace exporter: %w", err)
	}

	err = traceExporter.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("start trace exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSyncer(traceExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
			b3.New(),
		),
	)

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("setup metric exporter: %w", err)
	}

	meterProvider := controller.New(
		processor.NewFactory(
			simple.NewWithHistogramDistribution(),
			metricExporter,
		),
		controller.WithExporter(metricExporter),
		controller.WithCollectPeriod(2*time.Second),
	)
	err = meterProvider.Start(context.Background())
	if err != nil {
		return nil, fmt.Errorf("start otel metric pusher: %w", err)
	}

	global.SetMeterProvider(meterProvider)

	err = host.Start()
	if err != nil {
		return nil, fmt.Errorf("start otel host instrumentation: %w", err)
	}
	err = runtime.Start()
	if err != nil {
		return nil, fmt.Errorf("start otel runtime instrumentation: %w", err)
	}

	return &otelclient{
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		done:           make(chan error, 2),
	}, nil
}

func (o *otelclient) start() error {
	return <-o.done
}

func (o *otelclient) stop() error {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		o.done <- o.meterProvider.Stop(context.Background())
	}()
	go func() {
		defer wg.Done()
		o.done <- o.tracerProvider.Shutdown(context.Background())
	}()

	wg.Wait()
	return nil
}

func globalOtel(otlpEndpoint string, grpcOpts []grpc.DialOption) error {
	return nil
}
