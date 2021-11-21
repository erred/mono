package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	logrtesting "github.com/go-logr/logr/testing"
	"go.seankhliao.com/mono/svc/internal/testhttp"
)

func TestVanity(t *testing.T) {
	l := logrtesting.NewTestLogger(t)
	ctx := logr.NewContext(context.Background(), l)

	svr := httptest.NewServer(New(ctx))
	defer svr.Close()

	testhttp.Expect(
		t, svr.Client(), "vanity",
		http.MethodGet, svr.URL+"/foo", http.StatusOK,
		nil,
		[]string{
			`<meta name=go-import content="go.seankhliao.com/foo git https://github.com/seankhliao/foo">`,
			`<meta name=go-source content="go.seankhliao.com/foo
    https://github.com/seankhliao/foo`,
		},
	)
}
