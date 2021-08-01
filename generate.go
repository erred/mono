//go:build generate
// +build generate

package main

//go:generate buf generate

//go:generate go run go.seankhliao.com/mono/go/cmd/webrender -gtm GTM-TLVN7D6 -src ./blog -dst ./go/internal/w16/static/root
//go:generate go run go.seankhliao.com/mono/go/cmd/webrender -compact -src webextension/chrome-newtab/newtab.md -dst webextension/chrome-newtab/index.html
