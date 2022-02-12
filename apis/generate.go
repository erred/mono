//go:build generate

package apis

//go:generate -command buf go run github.com/bufbuild/buf/cmd/buf@latest
//go:generate buf generate
