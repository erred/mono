package main

import (
	"flag"
	"fmt"
	"os"

	"go.seankhliao.com/mono/internal/fin/run"
	"go.seankhliao.com/mono/internal/fin/run/format"
	"go.seankhliao.com/mono/internal/fin/run/importcsv"
	"go.seankhliao.com/mono/internal/fin/run/summary"
)

func main() {
	o := run.NewOptions(flag.CommandLine)
	flag.Parse()

	subCommands := map[string]run.Cmd{
		"importcsv": importcsv.Run,
		"summary":   summary.Run,
		"format":    format.Run,
	}
	commands := func() {
		scs := make([]string, 0, len(subCommands))
		for k := range subCommands {
			scs = append(scs, k)
		}
		fmt.Fprintf(os.Stderr, "commands:\n")
		for _, c := range scs {
			fmt.Fprintf(os.Stderr, "\t%s\n", c)
		}
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "no command given")
		commands()
		os.Exit(1)
	}

	cmd, ok := subCommands[flag.Arg(0)]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", flag.Arg(0))
		commands()
		os.Exit(1)
	}
	err := cmd(o, flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
