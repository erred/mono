package httpmid

import "net/http"

var _ http.ResponseWriter = &responseWrapper{}

type responseWrapper struct {
	base    http.ResponseWriter
	status  int
	written int
}

func newResponseWrapper(base http.ResponseWriter) *responseWrapper {
	return &responseWrapper{base: base, status: 200}
}

func (r *responseWrapper) WriteHeader(status int) {
	r.base.WriteHeader(status)
}

func (r *responseWrapper) Write(b []byte) (int, error) {
	n, err := r.base.Write(b)
	r.written += n
	return n, err
}

func (r *responseWrapper) Header() http.Header {
	return r.base.Header()
}
