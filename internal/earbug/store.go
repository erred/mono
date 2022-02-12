package earbug

import (
	"errors"
	"io/fs"
	"os"

	earbugv1 "go.seankhliao.com/mono/apis/earbug/v1"
	"google.golang.org/protobuf/proto"
)

func (s *Server) initStore() error {
	b, err := os.ReadFile(s.fname)
	if errors.Is(err, fs.ErrNotExist) {
		s.Store.Playbacks = make(map[string]*earbugv1.Playback)
		s.Store.Tracks = make(map[string]*earbugv1.Track)
		return nil
	} else if err != nil {
		return err
	}

	err = proto.Unmarshal(b, s.Store)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Write() error {
	b, err := proto.Marshal(s.Store)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.fname+".tmp", b, 0o644)
	if err != nil {
		return err
	}

	err = os.Rename(s.fname+".tmp", s.fname)
	if err != nil {
		return err
	}

	return nil
}
