package runsvr

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"runtime/debug"
	"strings"
)

// Desc adds the package path/version and a description to the flag usage.
// docgo should be the embedded contents of a doc.go file,
// the description is extracted from the package level documentation.
func Desc(flags *flag.FlagSet, docgo string) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("can't read buildinfo")
	}

	f, err := parser.ParseFile(token.NewFileSet(), "doc.go", docgo, parser.ParseComments)
	if err != nil {
		panic("parse doc.go")
	}

	flags.Usage = func() {
		fmt.Fprintf(flags.Output(), "%s\n%s %s\n\n%s\n\n",
			os.Args[0],
			bi.Path, bi.Main.Version,
			strings.TrimSpace(f.Doc.Text()),
		)
		flags.PrintDefaults()
	}
}
