package httpmid

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	httpmidv1 "go.seankhliao.com/mono/apis/httpmid/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccessLogOut struct {
	Log        zerolog.Logger
	PubSubject string
	Pub        interface {
		Publish(subject string, data any) error
	}
}

func AccessLog(h http.Handler, o AccessLogOut) http.Handler {
	if o.PubSubject == "" {
		o.PubSubject = "accesslog." + filepath.Base(os.Args[0])
	}
	tracer := otel.GetTracerProvider().Tracer("accesslog")
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "accesslog")
		defer span.End()

		r = r.WithContext(ctx)

		wrw, ok := rw.(*responseWrapper)
		if !ok {
			wrw = newResponseWrapper(rw)
			rw = wrw
		}

		t0 := time.Now()
		h.ServeHTTP(rw, r)
		t1 := time.Now()
		d := t1.Sub(t0)

		al := httpmidv1.AccessLog{
			Ts:            timestamppb.New(t0),
			TraceId:       span.SpanContext().TraceID().String(),
			SpanId:        span.SpanContext().SpanID().String(),
			HttpMethod:    r.Method,
			HttpUrl:       r.URL.String(),
			HttpVersion:   r.Proto,
			HttpHost:      r.Host,
			HttpUseragent: r.UserAgent(),
			HttpReferrer:  r.Referer(),
			HandleTime:    durationpb.New(d),
			HttpStatus:    int32(wrw.status),
			BytesWritten:  int64(wrw.written),
		}

		if o.Pub != nil {
			err := o.Pub.Publish(o.PubSubject, &al)
			if err != nil {
				o.Log.Err(err).Str("subject", o.PubSubject).Msg("publish access log")
			}
		}

		if o.Log.Trace().Enabled() {
			b, err := protojson.Marshal(&al)
			if err != nil {
				o.Log.Err(err).Msg("marshal access log to json")
			}
			o.Log.Trace().RawJSON("accesslog", b).Send()
		}
	})
}
