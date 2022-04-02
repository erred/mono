package flagwrap

import (
	"flag"
	"os"
	"strings"
)

// Wrapper for flag.FlagSet.Parse which will also read values from the env,
// with -hello.world read from HELLO_WORLD
func Parse(fset *flag.FlagSet, args []string) error {
	rep := strings.NewReplacer(".", "_", "-", "_")
	fset.VisitAll(func(f *flag.Flag) {
		key := strings.ToUpper(rep.Replace(f.Name))
		v, ok := os.LookupEnv(key)
		if ok {
			val := strings.TrimSpace(v)
			f.Value.Set(val)
		}
	})

	err := fset.Parse(args)
	if err != nil {
		return err
	}

	return nil
}
