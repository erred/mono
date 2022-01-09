package o11y

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Start(t trace.Tracer, l logr.Logger, ctx context.Context, name string) (context.Context, trace.Span, logr.Logger) {
	ctx, span := t.Start(ctx, name)
	spanctx := span.SpanContext()
	l = l.WithValues("trace_id", spanctx.TraceID(), "span_id", spanctx.SpanID())
	return ctx, span, l
}

func HttpError(rw http.ResponseWriter, l logr.Logger, span trace.Span, httpStatusCode int, err error, msg string, attrs ...attribute.KeyValue) {
	fmt.Fprintln(rw, msg)
	rw.WriteHeader(httpStatusCode)

	if err != nil {
		l.Error(err, msg, expandKvs(attrs...)...)
		span.SetStatus(codes.Error, msg)
		span.RecordError(err)
	} else {
		span.SetAttributes(attribute.String("info", msg))
		l.Info(msg, expandKvs(attrs...)...)
	}

	span.SetAttributes(attrs...)
}

func expandKvs(attrs ...attribute.KeyValue) []interface{} {
	out := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		out = append(out, string(attr.Key), attr.Value.Emit())
	}
	return out
}
