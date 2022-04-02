package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.seankhliao.com/mono/internal/flagwrap"
	"go.seankhliao.com/mono/internal/gchat"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	var endpoint, msg string
	fset := flag.NewFlagSet("", flag.ContinueOnError)
	fset.StringVar(&endpoint, "gchat.webhook", "", "webhook endpoint")
	fset.StringVar(&msg, "gchat.message", "", "message")
	err := flagwrap.Parse(fset, os.Args[1:])
	if err != nil {
		return err
	}

	if endpoint == "" {
		return errors.New("no webhook provided")
	}

	if fset.NArg() > 0 {
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

	err = client.Post(context.TODO(), gchat.WebhookPayload{
		Text: msg,
	})
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	fmt.Println("ok")

	return nil
}
