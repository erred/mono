package httpsvr

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/go-logr/logr"
	logrtesting "github.com/go-logr/logr/testing"
	"go.seankhliao.com/mono/svc/internal/testhttp"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = logr.NewContext(ctx, logrtesting.NewTestLogger(t))
	o := &Options{
		BaseContext: ctx,
		AdminAddr:   ":57890",
		HTTPAddr:    ":56789",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("test x"))
		}),
	}

	go o.Run()

	var startCount int
	interval := 100 * time.Millisecond
	for range time.NewTicker(interval).C {
		startCount++
		if startCount > 20 {
			t.Errorf("server taking too long to start time=%v", time.Duration(startCount)*interval)
			return
		}
		ok := func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1:57890/readyz", nil)
			if err != nil {
				t.Errorf("prepare request err=%v", err)
				return false
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return false
			}
			if res.StatusCode != 200 {
				return false
			}
			return true
		}()
		if ok {
			break
		}
	}

	testhttp.Expect(t, http.DefaultClient, "admin",
		http.MethodGet, "http://127.0.0.1:57890/readyz", http.StatusOK,
		nil,
		[]string{"ok"},
	)
	testhttp.Expect(t, http.DefaultClient, "http",
		http.MethodGet, "http://127.0.0.1:56789/", http.StatusOK,
		nil,
		[]string{"test x"},
	)
}
