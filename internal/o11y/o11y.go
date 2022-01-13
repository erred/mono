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

func Set(l logr.Logger, span trace.Span, attrs ...attribute.KeyValue) logr.Logger {
	span.SetAttributes(attrs...)
	l = l.WithValues(expandKvs(attrs...)...)
	return l
}

func OK(l logr.Logger, span trace.Span, msg string, attrs ...attribute.KeyValue) {
	span.SetStatus(codes.Ok, msg)
	span.SetAttributes(attribute.String("msg", msg))

	Set(l, span, attrs...).Info(msg)
}

func Error(l logr.Logger, span trace.Span, err error, msg string, attrs ...attribute.KeyValue) {
	span.SetStatus(codes.Error, msg)
	span.RecordError(err)

	Set(l, span, attrs...).Error(err, msg)
}

func HttpError(rw http.ResponseWriter, l logr.Logger, span trace.Span, httpStatusCode int, err error, msg string, attrs ...attribute.KeyValue) {
	fmt.Fprintln(rw, msg)
	rw.WriteHeader(httpStatusCode)

	Error(l, span, err, msg, attrs...)
}

func expandKvs(attrs ...attribute.KeyValue) []interface{} {
	out := make([]interface{}, 0, len(attrs)*2)
	for _, attr := range attrs {
		out = append(out, string(attr.Key), attr.Value.Emit())
	}
	return out
}
