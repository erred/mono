
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

MAP_BASE := internal/w16/static/root/static

.PHONY: map
map:
	convert -background none \
		-density 1200 -resize 1920x1080 $(MAP_BASE)/map.svg \
		-write $(MAP_BASE)/map.png \
		-write $(MAP_BASE)/map.webp \
		$(MAP_BASE)/map.jpg
