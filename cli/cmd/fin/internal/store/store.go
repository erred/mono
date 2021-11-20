package store

import (
	"fmt"
	"io"
	"os"

	fin "go.seankhliao.com/mono/proto/seankhliao/fin/v1alpha1"
	"google.golang.org/protobuf/encoding/prototext"
)

func ReadFile(name string) (*fin.All, error) {
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

func Read(r io.Reader) (*fin.All, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("readall: %w", err)
	}

	var all fin.All
	err = prototext.Unmarshal(b, &all)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &all, nil
}

func WriteFile(name string, all *fin.All) error {
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

func Write(w io.Writer, all *fin.All) error {
	o := prototext.MarshalOptions{
		Multiline: true,
	}
	b, err := o.Marshal(all)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
