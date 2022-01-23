//go:build tools

package tools

import (
	// skaffold is on too old a version of opentelemetry
	// _ "github.com/GoogleContainerTools/skaffold/cmd/skaffold"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/google/ko"
)

//go:generate mkdir -p bin
//go:generate go build -o bin/ github.com/bufbuild/buf/cmd/buf
//go:generate go build -o bin/ github.com/google/ko
