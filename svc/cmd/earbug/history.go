package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.seankhliao.com/mono/internal/web/render"
	"go.seankhliao.com/mono/internal/o11y"
)

// /user/history
func (s *Server) handleUserHistory(rw http.ResponseWriter, r *http.Request) {
	ctx, span, l := o11y.Start(s.t, s.l, r.Context(), "user_history")
	defer span.End()

	id := r.Header.Get("auth-id")
	if id == "" || id == "anonymous" {
		o11y.HttpError(rw, l, span, http.StatusUnauthorized, nil, "no auth-id")
		return
	}

	startKey := path.Join(s.StorePrefix, "history", id, "playback")
	endKey := path.Join(s.StorePrefix, "history", id, "playback2")
	gr, err := s.Store.Get(ctx, startKey, clientv3.WithRange(endKey))
	if err != nil {
		o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "get key", attribute.String("start_key", startKey))
		return
	}

	trackMap := make(map[string]spotify.SimpleTrack)

	var buf bytes.Buffer
	buf.WriteString("### _Listening_ History\n\n")
	fmt.Fprintf(&buf, "%v entries for [_%s_](/user)\n\n", len(gr.Kvs), id)
	buf.WriteString(`
| idx | time | track | artist |
| --- | ---- | ----- | ------ |
`)

	for i := range gr.Kvs {
		kv := gr.Kvs[len(gr.Kvs)-1-i]

		rts := strings.TrimPrefix(string(kv.Key), startKey+"/")
		ts, err := time.Parse(time.RFC3339Nano, rts)
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "unmarshal key", attribute.String("key", string(kv.Key)))
			return
		}

		var ph PlaybackHistory
		err = json.Unmarshal(kv.Value, &ph)
		if err != nil {
			o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "unmarshal value", attribute.String("key", string(kv.Key)))
			return
		}

		track, ok := trackMap[ph.TrackID]
		if !ok {
			trackP := path.Join(s.StorePrefix, "track", ph.TrackID)
			tgr, err := s.Store.Get(ctx, trackP)
			if err != nil {
				o11y.HttpError(rw, l, span, http.StatusInternalServerError, err, "get key", attribute.String("key", trackP))
				return
			}
			err = json.Unmarshal(tgr.Kvs[0].Value, &track)
			if err != nil {
				http.Error(rw, "unmarshal track", http.StatusInternalServerError)
				return
			}

			trackMap[ph.TrackID] = track
		}

		fmt.Fprintf(
			&buf,
			"| %v | %s | %s | %v |\n",
			i+1,
			ts.Format(time.RFC3339),
			track.Name,
			track.Artists[0].Name,
		)
	}

	err = render.Compact(
		rw,
		"earbug view",
		"view listening history",
		s.CanonicalURL+"/user/history",
		buf.Bytes(),
	)
	if err != nil {
		l.Error(err, "render")
		return
	}
}
