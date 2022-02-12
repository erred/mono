package lg

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	l zerolog.Logger
}

func New() *Logger {
	return &Logger{
		l: zerolog.New(os.Stderr),
	}
}

func (l Logger) Info(ctx context.Context, msg string, kvs ...any) {
	e := l.l.Info()
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		e = e.Stringer("trace_id", spanCtx.TraceID()).Stringer("span_id", spanCtx.SpanID())
	}
	e.Fields(kvs).Msg(msg)
}

func (l Logger) Error(ctx context.Context, err error, msg string, kvs ...any) {
	e := l.l.Error().Err(err)
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		e = e.Stringer("trace_id", spanCtx.TraceID()).Stringer("span_id", spanCtx.SpanID())
	}
	e.Fields(kvs).Msg(msg)
}
