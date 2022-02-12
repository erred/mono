package o11y

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-logr/logr/funcr"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.seankhliao.com/mono/monolith/run"
	"google.golang.org/grpc"
)

var _ run.Runner = &Exporter{}

type Exporter struct {
	enabled bool
	Addr    string

	initOnce       sync.Once
	res            *resource.Resource
	otlpConn       *grpc.ClientConn
	metricExporter *otlpmetric.Exporter
	metricProvider *controller.Controller
	traceExporter  *otlptrace.Exporter
	traceProvider  *trace.TracerProvider
}

func New() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Enabled() bool { return e.enabled }

func (e *Exporter) Register(flags *flag.FlagSet) {
	flags.BoolVar(&e.enabled, "o11y", false, "enable exporters")
	flags.StringVar(&e.Addr, "o11y.otlp", "127.0.0.1:4317", "otel collector endpoint")
}

func (e *Exporter) Provider() ToolProvider {
	e.initOnce.Do(func() {
		if e.enabled {
			e.otelResources()

			ctx := context.TODO()
			err := e.otelGRPC(ctx, e.Addr)
			if err != nil {
				os.Exit(1)
			}

			e.otelTracer()
			e.otelMetrics()
		}
	})

	return ToolProvider{
		L: funcr.NewJSON(func(obj string) {
			fmt.Fprintln(os.Stderr, obj)
		}, funcr.Options{
			LogTimestamp:    true,
			TimestampFormat: time.RFC3339,
		}),
		T: otel.GetTracerProvider(), // has noop defaults
		M: global.GetMeterProvider(),
	}
}

func (e *Exporter) Run(ctx context.Context) error {
	noop := func(ctx context.Context) error {
		<-ctx.Done()
		return nil
	}

	if e.enabled {
		type runner struct {
			start, stop func(context.Context) error
		}
		runners := map[string]runner{
			"metric provider": {e.metricProvider.Start, e.metricProvider.Stop},
			"metric exporter": {e.metricExporter.Start, e.metricExporter.Shutdown},
			"trace provider":  {noop, e.traceProvider.Shutdown},
			"trace exporter":  {e.traceExporter.Start, e.traceExporter.Shutdown},
		}

		var wg sync.WaitGroup
		for n, r := range runners {
			err := r.start(ctx)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func(n string, r runner) {
				defer wg.Done()
				<-ctx.Done()
				r.stop(context.Background())
			}(n, r)
		}

		// extra instrumentation
		err := otelruntime.Start(otelruntime.WithMeterProvider(e.metricProvider))
		if err != nil {
			return err
		}
		wg.Wait()
	} else {
		noop(ctx)
	}

	return nil
}

func (e *Exporter) otelResources() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		e.res = resource.Default()
		return
	}

	// ignore merge errors from mismatched schemas
	e.res, _ = resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(bi.Path),
			semconv.ServiceVersionKey.String(bi.Main.Version),
			semconv.ProcessRuntimeVersionKey.String(runtime.Version()),
		),
	)
}

func (e *Exporter) otelGRPC(ctx context.Context, endpoint string) error {
	dos := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	var err error
	e.otlpConn, err = grpc.DialContext(ctx, endpoint, dos...)
	return err
}

func (e *Exporter) otelMetrics() {
	e.metricExporter = otlpmetric.NewUnstarted(
		otlpmetricgrpc.NewClient(
			otlpmetricgrpc.WithGRPCConn(e.otlpConn),
		),
	)

	// provider
	e.metricProvider = controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithExporter(e.metricExporter),
		controller.WithCollectPeriod(10*time.Second),
		controller.WithResource(e.res),
	)

	// global
	global.SetMeterProvider(e.metricProvider)
}

func (e *Exporter) otelTracer() {
	// Exporter
	e.traceExporter = otlptrace.NewUnstarted(
		otlptracegrpc.NewClient(
			otlptracegrpc.WithGRPCConn(e.otlpConn),
		),
	)

	// provider
	e.traceProvider = trace.NewTracerProvider(
		trace.WithResource(e.res),
		trace.WithBatcher(e.traceExporter),
	)

	// global
	otel.SetTracerProvider(e.traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{},
	))
}
