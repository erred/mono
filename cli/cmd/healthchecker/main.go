package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	healthcheckerv1 "go.seankhliao.com/mono/apis/healthchecker/v1"
	"go.seankhliao.com/mono/internal/gchat"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	c, err := GetConfig()
	if err != nil {
		return fmt.Errorf("get config: %w", err)
	}

	state, err := GetState(c.StateDir)
	if err != nil {
		return fmt.Errorf("get state: %w", err)
	}

	httpClient := &http.Client{
		Transport: &trWrap{c.UserAgent, http.DefaultTransport},
	}

	var httpResMu sync.Mutex
	httpRes := make(map[string]error)

	var wg sync.WaitGroup
	for _, check := range c.Http {
		wg.Add(1)
		go func(check *healthcheckerv1.HttpCheckConfig) {
			defer wg.Done()
			err := CheckHttp(httpClient, check)
			httpResMu.Lock()
			defer httpResMu.Unlock()
			httpRes[check.Name] = err
		}(check)
	}

	wg.Wait()

	client := gchat.WebhookClient{
		Endpoint: c.Notify.Gchat,
		Client:   httpClient,
	}

	for name, err := range httpRes {
		old, ok := state.Results[name]
		if !ok {
			t := timestamppb.Now()
			res := &healthcheckerv1.CheckResult{
				Protocol:  "http",
				TsInitial: t,
				TsLatest:  t,
			}
			if err != nil {
				res.Pass = false
				res.Details = err.Error()
			} else {
				res.Pass = true
			}
			state.Results[name] = res
			if err != nil {
				err := client.Post(context.TODO(), gchat.WebhookPayload{
					Text: fmt.Sprintf("*FAIL:* %s: %s", name, res.Details),
				})
				if err != nil {
					log.Println(err)
				}
			}
			continue
		}
		if !old.Pass && err != nil { // still failing
			old.TsLatest = timestamppb.Now()
			old.Details = err.Error()
		} else if !old.Pass { // fail -> pass
			t := timestamppb.Now()
			d := t.AsTime().Sub(old.TsInitial.AsTime())
			old.TsInitial = t
			old.TsLatest = t
			old.Pass = true
			old.Details = ""
			err := client.Post(context.TODO(), gchat.WebhookPayload{
				Text: fmt.Sprintf("*PASS:* %s (failed for %v)", name, d),
			})
			if err != nil {
				log.Println(err)
			}
		} else if err != nil { // pass -> fail
			t := timestamppb.Now()
			old.TsInitial = t
			old.TsLatest = t
			old.Pass = false
			old.Details = err.Error()
			err := client.Post(context.TODO(), gchat.WebhookPayload{
				Text: fmt.Sprintf("*FAIL:* %s: %s", name, old.Details),
			})
			if err != nil {
				log.Println(err)
			}
		} else { // still passsing
			old.TsLatest = timestamppb.Now()
		}
		state.Results[name] = old
	}

	return nil
}

func GetConfig() (*healthcheckerv1.Config, error) {
	confPath := os.Getenv("HEALTHCHECKER_CONFIG_PATH")
	if confPath == "" {
		confPath = "/etc/mono/healthchecker/config.prototext"
	}
	b, err := os.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("read config from %s: %w", confPath, err)
	}
	var config healthcheckerv1.Config
	err = prototext.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}

func GetState(stateDir string) (*healthcheckerv1.State, error) {
	statePath := filepath.Join(stateDir, "state.prototext")
	b, err := os.ReadFile(statePath)
	if errors.Is(err, fs.ErrNotExist) {
		return &healthcheckerv1.State{Results: make(map[string]*healthcheckerv1.CheckResult)}, nil
	} else if err != nil {
		return nil, fmt.Errorf("read state from %s: %w", statePath, err)
	}
	var state healthcheckerv1.State
	err = prototext.Unmarshal(b, &state)
	if err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return &state, nil
}

func CheckHttp(client *http.Client, c *healthcheckerv1.HttpCheckConfig) error {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, c.Url, nil)
	if err != nil {
		return fmt.Errorf("prepare request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	s := string(b)
	if c.MatchExact != "" {
		if !strings.Contains(s, c.MatchExact) {
			return fmt.Errorf("missing exact match for %q\nbody: %s", c.MatchExact, limitString(s))
		}
	}
	if c.MatchRegex != "" {
		rgx, err := regexp.Compile(c.MatchRegex)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}
		if !rgx.MatchString(s) {
			return fmt.Errorf("missing regex match for %q\nbody: %s", c.MatchRegex, limitString(s))
		}
	}

	return nil
}

type trWrap struct {
	ua string
	http.RoundTripper
}

func (t *trWrap) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("user-agent", t.ua)
	return t.RoundTripper.RoundTrip(r)
}

func limitString(s string) string {
	if len(s) > 200 {
		return fmt.Sprintf("%s ...(%d bytes omitted)... %s", s[:100], len(s)-200, s[len(s)-100:])
	}
	return s
}
