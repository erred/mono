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
	"time"

	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"go.seankhliao.com/mono/internal/flagwrap"
	"google.golang.org/protobuf/proto"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run() error {
	var fname string
	fset := flag.NewFlagSet("", flag.ContinueOnError)
	fset.StringVar(&fname, "earbug.data", "/var/lib/mono/earbug/earbug.pb", "path to data file")
	err := flagwrap.Parse(fset, os.Args[1:])
	if err != nil {
		return err
	}

	b, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("read %s: %w", fname, err)
	}
	var store earbugv1.Store
	err = proto.Unmarshal(b, &store)
	if err != nil {
		return fmt.Errorf("unmarshal store: %w", err)
	}

	datePrefix := time.Now().Format("2006-01-02")

	prevTrack, newTrack := make(map[string]struct{}), make(map[string]struct{})
	var newPlaybacks int
	for ts, play := range store.Playbacks {
		if !strings.HasPrefix(ts, datePrefix) {
			prevTrack[play.TrackId] = struct{}{}
		} else {
			newPlaybacks++
			newTrack[play.TrackId] = struct{}{}
		}
	}
	for k := range newTrack {
		if _, ok := prevTrack[k]; ok {
			delete(newTrack, k)
		}
	}
	newTracks := len(newTrack)

	p := Payload{fmt.Sprintf("you listened to *%d* tracks, of which _%d_ were new", newPlaybacks, newTracks)}
	b, err = json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	endpoint := os.Getenv("GCHAT_WEBHOOK")
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return errors.New("no webhook endpoint provided")
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

	return nil
}

type Payload struct {
	Text string `json:"text"`
}
