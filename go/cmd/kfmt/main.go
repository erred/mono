// kfmt is a thin wrapper around sigs.k8s.io/kustomize/kyaml/kio/filters
// to split and format k8s manifest files consistently,
// defaulting to %k_%n.yaml
package main

import (
	"os"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/kio/filters"
)

const (
	FilenamePattern = "%k_%n.yaml"
)

func main() {
	p := "."
	if len(os.Args) > 1 {
		p = os.Args[1]
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
				FilenamePattern: FilenamePattern,
				Override:        true,
			},
		},
		Outputs: []kio.Writer{rw},
	}.Execute()
	if err != nil {
		panic(err)
	}
}
