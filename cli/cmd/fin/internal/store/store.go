package store

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"go.seankhliao.com/mono/proto/finpb"
	"google.golang.org/protobuf/encoding/prototext"
)

func ReadFile(name string) (*finpb.All, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", name, err)
	}
	defer f.Close()

	trs, err := Read(f)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", name, err)
	}
	return trs, nil
}

func Read(r io.Reader) (*finpb.All, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	var all finpb.All
	err = prototext.Unmarshal(b, &all)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &all, nil
}

func WriteFile(name string, all *finpb.All) error {
	f, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("create %s: %w", name, err)
	}
	defer f.Close()

	err = Write(f, all)
	if err != nil {
		return fmt.Errorf("write %s: %w", name, err)
	}
	return nil
}

func Write(w io.Writer, all *finpb.All) error {
	o := prototext.MarshalOptions{
		Multiline: true,
	}
	b, err := o.Marshal(all)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	// collapse transactions onto single line
	b = regexp.MustCompile(`\n    `).ReplaceAll(b, []byte("\t"))
	b = regexp.MustCompile(`\n  }`).ReplaceAll(b, []byte(` }`))
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
