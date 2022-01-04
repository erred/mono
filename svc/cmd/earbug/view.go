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
	"go.seankhliao.com/mono/internal/web/render"
)

// /view/user/<user>
func (s *Server) viewUser(rw http.ResponseWriter, r *http.Request) {
	l := s.l.WithName("view").WithValues("page", "user")
	ctx := r.Context()

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 || parts[1] != "view" || parts[2] != "user" {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}
	user := parts[3]

	startKey := path.Join(s.StorePrefix, "history", user, "playback")
	endKey := path.Join(s.StorePrefix, "history", user, "playback2")
	gr, err := s.Store.Get(ctx, startKey, clientv3.WithRange(endKey))
	if err != nil {
		l.Error(err, "get user history", "startKey", startKey)
		http.Error(rw, "get history", http.StatusInternalServerError)
		return
	}

	trackMap := make(map[string]spotify.SimpleTrack)

	var buf bytes.Buffer
	buf.WriteString("### _Listening_ History\n\n")
	fmt.Fprintf(&buf, "%v entries for _%s_\n\n", len(gr.Kvs), user)
	buf.WriteString(`
| idx | time | track | artist |
| --- | ---- | ----- | ------ |
`)

	for i := range gr.Kvs {
		kv := gr.Kvs[len(gr.Kvs)-1-i]

		rts := strings.TrimPrefix(string(kv.Key), startKey+"/")
		ts, err := time.Parse(time.RFC3339Nano, rts)
		if err != nil {
			l.Error(err, "unmarshal timestamp", "key", string(kv.Key))
			http.Error(rw, "unmarshal history", http.StatusInternalServerError)
			return
		}

		var ph PlaybackHistory
		err = json.Unmarshal(kv.Value, &ph)
		if err != nil {
			l.Error(err, "unmarshal history", "key", string(kv.Key))
			http.Error(rw, "unmarshal history", http.StatusInternalServerError)
			return
		}

		track, ok := trackMap[ph.TrackID]
		if !ok {
			trackP := path.Join(s.StorePrefix, "track", ph.TrackID)
			tgr, err := s.Store.Get(ctx, trackP)
			if err != nil {
				l.Error(err, "get track", "key", trackP)
				http.Error(rw, "get track", http.StatusInternalServerError)
				return
			}
			err = json.Unmarshal(tgr.Kvs[0].Value, &track)
			if err != nil {
				l.Error(err, "unmarshal track", "key", trackP)
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
		s.CanonicalURL+"/view/user/"+user,
		buf.Bytes(),
	)
	if err != nil {
		l.Error(err, "render")
		return
	}
}
