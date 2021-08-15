
GO := gotip
BUF := buf

.PHONY: generate
generate: proto chrome-newtab

.PHONY: proto
proto:
	$(BUF) generate

.PHONY: chrome-newtab
chrome-newtab:
	$(GO) run go.seankhliao.com/mono/go/cmd/webrender -compact -src webextension/chrome-newtab/newtab.md -dst webextension/chrome-newtab/index.html
