package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `post text messages to a google chat workspace

GCHAT_WEBHOOK
        webhook url
GCHAT_MESSAGE
        text message (if not passed via args)
`)
	}
	flag.Parse()

	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	endpoint := os.Getenv("GCHAT_WEBHOOK")
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return errors.New("no webhook provided")
	}

	msg := os.Getenv("GCHAT_MESSAGE")
	if len(flag.Args()) > 0 {
		msg = strings.Join(flag.Args(), " ")
	}
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return errors.New("no message provided")
	}

	p := Payload{msg}
	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("post request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("post status: %s", res.Status)
	}
	// defer res.Body.Close()
	// b, err = io.ReadAll(res.Body)
	// if err != nil {
	// 	return fmt.Errorf("read response: %w", err)
	// }
	fmt.Println("ok")

	return nil
}

type Payload struct {
	Text string `json:"text"`
}
