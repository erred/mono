package testhttp

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func Expect(t *testing.T, c *http.Client, name, method, url string, code int, expectedHeaders map[string]string, expectedBodies []string) {
	t.Run(name, func(t *testing.T) {
		t.Logf("%s %s", method, url)
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			t.Errorf("create request err=%v", err)
			return
		}
		res, err := c.Do(req)
		if err != nil {
			t.Errorf("Do err=%v", err)
			return
		}
		defer res.Body.Close()
		if res.StatusCode != code {
			t.Errorf("statuscode=%v expected=%v", res.StatusCode, code)
			return
		}
		for k, v := range expectedHeaders {
			vv := res.Header.Get(k)
			if vv != v {
				t.Errorf("expectedHeader %v=%v expected=%v", k, vv, v)
			}

		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("read body err=%v", err)
			return
		}
		bs := string(b)
		for i, eb := range expectedBodies {
			if !strings.Contains(bs, eb) {
				t.Errorf("expectedBody[%d] expected=%s", i, eb)
			}
		}
		if t.Failed() {
			t.Logf("body=%s", bs)
		}
	})
}
