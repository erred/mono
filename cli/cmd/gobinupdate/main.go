package main

import (
	"debug/buildinfo"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%s [paths...]
        Looks for binaries built by go and attempts to install their latest version
        by looking at their embedded build info.
        [path] can be either a file or a directory (non recursive).
        Defaults to $GOBIN or $GOPATH/bin if no arguments are passed.
`,
			os.Args[0],
		)
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	paths := flag.Args()
	if len(paths) == 0 {
		bin := os.Getenv(`GOBIN`)
		if bin == "" {
			bin = os.Getenv("GOPATH")
			if bin == "" {
				fmt.Fprintln(os.Stderr, "GOBIN and GOPATH unset")
			}
			bin = filepath.Join(bin, "bin")
		}
		paths = []string{bin}
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, p := range paths {
		wg.Add(1)
		work(&wg, p)
	}
}

func work(wg *sync.WaitGroup, p string) {
	defer wg.Done()
	fi, err := os.Stat(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: stat: %v\n", p, err)
		return
	}
	if !fi.IsDir() {
		try(wg, p)
	}
	des, err := os.ReadDir(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: readdir: %v\n", p, err)
	}
	for _, de := range des {
		wg.Add(1)
		go try(wg, filepath.Join(p, de.Name()))
	}
}

func try(wg *sync.WaitGroup, p string) {
	defer wg.Done()

	bi, err := buildinfo.ReadFile(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: read build info: %v\n", p, err)
		return
	}
	if bi.Path == "" {
		fmt.Fprintf(os.Stderr, "%s: no module info\n", p)
		return
	}

	b, err := exec.Command("go", "install", bi.Path+"@latest").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: updating %s err=%v out=%s\n", p, bi.Path, err, b)
		return
	}
	fmt.Fprintf(os.Stderr, "%s: updated %s\n", p, bi.Path)
}
