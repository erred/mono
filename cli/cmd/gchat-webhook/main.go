package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.seankhliao.com/mono/internal/gchat"
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

	client := &gchat.WebhookClient{
		Endpoint: endpoint,
		Client:   http.DefaultClient,
	}

	err := client.Post(context.TODO(), gchat.WebhookPayload{
		Text: msg,
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	fmt.Println("ok")

	return nil
}
