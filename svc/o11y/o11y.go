// Package o11y provices a standardized initial setup
// for observability (o11y) instrumentation.
package o11y

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/go-logr/stdr"
	"github.com/iand/logfmtr"
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
	"google.golang.org/grpc/credentials"

	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var ErrSetup = errors.New("o11y: error setting up observability instrumentation")

type Options struct {
	BaseContext context.Context

	LogFormat         string
	GRPCEndpoint      string
	GRPCInsecure      bool
	GRPCClientCrtFile string
	GRPCClientKeyFile string
	GRPCServerCAFile  string

	l logr.Logger
}

// NewOptions registers flags for th default options
func NewOptions(fs *flag.FlagSet) *Options {
	var o Options
	fs.StringVar(&o.LogFormat, "log.format", "json", "log output format: text|json|logfmt")
	fs.StringVar(&o.GRPCEndpoint, "otel.grpc.endpoint", "", "otel collector grpc host port")
	fs.BoolVar(&o.GRPCInsecure, "otel.grpc.insecure", false, "disable TLS for otel grpc")
	fs.StringVar(&o.GRPCClientKeyFile, "otel.grpc.clientkey", "", "path to otel grpc client tls key file")
	fs.StringVar(&o.GRPCClientCrtFile, "otel.grpc.clientcrt", "", "path to otel grpc client tls cert file")
	fs.StringVar(&o.GRPCServerCAFile, "otel.grpc.serverca", "", "path to otel grpc server ca cert file")
	return &o
}

// New returns an initiallized base context with a logger
// and sets up global trace/metrics providers.
// calls os.Exit(1) on failure
func (o *Options) New() context.Context {
	// logger
	var l logr.Logger
	switch o.LogFormat {
	case "text":
		l = stdr.New(log.New(os.Stderr, "", log.LstdFlags))
	case "logfmt":
		l = logfmtr.NewWithOptions(logfmtr.Options{
			Writer:          os.Stderr,
			TimestampFormat: time.RFC3339,
		})
	case "json":
		fallthrough
	default:
		l = funcr.NewJSON(func(obj string) {
			fmt.Fprintln(os.Stderr, obj)
		}, funcr.Options{
			LogTimestamp:    true,
			TimestampFormat: time.RFC3339,
		})
	}

	o.l = l.WithName("o11y")

	// context
	ctx := o.BaseContext
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = logr.NewContext(ctx, l)

	if o.GRPCEndpoint == "" {
		l.Info("no otel grpc endpoint set, skipping")
		return ctx
	}

	// global otel setup
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		o.l.Error(err, "opentelemetry error")
	}))
	res, err := o.otelResources(ctx)
	if err != nil {
		os.Exit(1)
	}
	conn, err := o.otelGRPC(ctx)
	if err != nil {
		os.Exit(1)
	}

	// tracer
	err = o.otelTracer(ctx, res, conn)
	if err != nil {
		os.Exit(1)
	}

	// metrics
	err = o.otelMetrics(ctx, res, conn)
	if err != nil {
		os.Exit(1)
	}

	return ctx
}

// otelTracer sets up a global tracer pushing to the given grpc endpoint
func (o *Options) otelTracer(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) error {
	// Exporter
	exporter, err := otlptrace.New(ctx,
		otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn)),
	)
	if err != nil {
		o.l.Error(err, "init trace exporter")
		return ErrSetup
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

// otelMetrics sets up a global metric provider pushing to the given grpc endpoint
func (o *Options) otelMetrics(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) error {
	exporter, err := otlpmetric.New(ctx,
		otlpmetricgrpc.NewClient(otlpmetricgrpc.WithGRPCConn(conn)),
	)
	if err != nil {
		o.l.Error(err, "init metric exporter")
		return ErrSetup
	}

	// provider
	provider := controller.New(
		processor.NewFactory(
			selector.NewWithExactDistribution(),
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

// otelResources returns a set of resources to be associated with the trace/metrics sdks
func (o *Options) otelResources(ctx context.Context) (*resource.Resource, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		o.l.Info("no no build info found")
		return resource.Default(), nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceVersionKey.String(bi.Main.Version),
			semconv.ProcessRuntimeVersionKey.String(runtime.Version()),
		),
	)
	if err != nil {
		o.l.Error(err, "merge resources")
		return nil, ErrSetup
	}
	return res, nil
}

// otelGRPC sets up a shared grpc connection for trace/metrics
func (o *Options) otelGRPC(ctx context.Context) (*grpc.ClientConn, error) {
	var dos []grpc.DialOption

	if o.GRPCInsecure {
		dos = append(dos, grpc.WithInsecure())
	} else {
		var tlscfg tls.Config

		// Server CA
		if o.GRPCServerCAFile != "" {
			b, err := os.ReadFile(o.GRPCServerCAFile)
			if err != nil {
				o.l.Error(err, "read server ca file", "ca_file", o.GRPCServerCAFile)
				return nil, ErrSetup
			}
			tlscfg.RootCAs = x509.NewCertPool()
			if !tlscfg.RootCAs.AppendCertsFromPEM(b) {
				o.l.Error(err, "append server ca file", "ca_file", o.GRPCServerCAFile)
				return nil, ErrSetup
			}
		} else {
			var err error
			tlscfg.RootCAs, err = x509.SystemCertPool()
			if err != nil {
				o.l.Error(err, "get system cert pool")
				return nil, ErrSetup
			}
		}

		// Client Certs
		if o.GRPCClientKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(o.GRPCClientCrtFile, o.GRPCClientKeyFile)
			if err != nil {
				o.l.Error(err, "load client certs", "client_cert", o.GRPCClientCrtFile, "client_key", o.GRPCClientKeyFile)
				return nil, ErrSetup
			}
			tlscfg.Certificates = []tls.Certificate{cert}
		}

		dos = append(dos, grpc.WithTransportCredentials(credentials.NewTLS(&tlscfg)))
	}

	conn, err := grpc.DialContext(ctx, o.GRPCEndpoint, dos...)
	if err != nil {
		o.l.Error(err, "dial otel", "endpoint, o.GRPCEndpoint")
		return nil, ErrSetup
	}
	return conn, nil
}
