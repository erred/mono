package format

import (
	"flag"
	"fmt"
	"os"

	"go.seankhliao.com/mono/cli/cmd/fin/internal/run"
	"go.seankhliao.com/mono/cli/cmd/fin/internal/store"
)

func Run(o run.Options, args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	err := fs.Parse(args[1:])
	if err != nil {
		return err
	}

	return format(o.File)
}

func format(fn string) error {
	all, err := store.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("read %s: %w", fn, err)
	}
	tmp := fn + ".tmp"
	err = store.WriteFile(tmp, all)
	if err != nil {
		return fmt.Errorf("write %s: %w", tmp, err)
	}
	return os.Rename(tmp, fn)
}
