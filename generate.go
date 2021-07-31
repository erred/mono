//go:build generate
// +build generate

package main

//go:generate buf generate

//go:generate go run go.seankhliao.com/mono/go/cmd/webrender -src ./blog -dst ./go/internal/w16/static/root -gtm GTM-TLVN7D6

//go:generate go run go.seankhliao.com/mono/go/cmd/webrender -compact -src webextension/chrome-newtab/newtab.md -dst webextension/chrome-newtab/index.html
