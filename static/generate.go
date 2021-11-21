//go:build generate

package static

//go:generate convert -background none -density 1200 -resize 1920x1080 seankhliao.com/static/map.svg -write seankhliao.com/static/map.png -write seankhliao.com/static/map.webp seankhliao.com/static/map.jpg
