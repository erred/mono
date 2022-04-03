package earbug

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"google.golang.org/protobuf/proto"
)

func (s *Server) initStore() error {
	s.Store = &earbugv1.Store{}

	b, err := os.ReadFile(s.fname)
	if errors.Is(err, fs.ErrNotExist) {
		s.Store.Playbacks = make(map[string]*earbugv1.Playback)
		s.Store.Tracks = make(map[string]*earbugv1.Track)
		return nil
	} else if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	err = proto.Unmarshal(b, s.Store)
	if err != nil {
		return fmt.Errorf("unmarshal stored data: %w", err)
	}

	s.log.Debug().Bool("stored_token", len(s.Store.Token) > 0).Int("stored_tracks", len(s.Store.Tracks)).Int("stored_plays", len(s.Store.Playbacks)).Msg("got stored data")
	return nil
}

func (s *Server) Write() error {
	b, err := proto.Marshal(s.Store)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	err = os.WriteFile(s.fname+".tmp", b, 0o644)
	if err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}

	err = os.Rename(s.fname+".tmp", s.fname)
	if err != nil {
		return fmt.Errorf("rename files: %w", err)
	}

	return nil
}
