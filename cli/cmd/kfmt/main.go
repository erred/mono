// kfmt is a thin wrapper around sigs.k8s.io/kustomize/kyaml/kio/filters
// to split and format k8s manifest files consistently,
// defaulting to %k_%n.yaml
package main

import (
	"flag"
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
)

func main() {
	var filenamePattern string
	flag.StringVar(&filenamePattern, "filename-pattern", "%k.%n.k8s.yaml", "pattern for filenames (%k: kind, %n: name, %s: namespace)")
	flag.Parse()

	p := flag.Arg(0)
	if p == "" {
		p = "."
	}

	rw := &kio.LocalPackageReadWriter{
		PackagePath: p,
	}

	err := kio.Pipeline{
		Inputs: []kio.Reader{rw},
		Filters: []kio.Filter{
			&filters.FormatFilter{
				UseSchema: true,
			},
			&filters.FileSetter{
				FilenamePattern: filenamePattern,
				Override:        true,
			},
		},
		Outputs: []kio.Writer{rw},
	}.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
