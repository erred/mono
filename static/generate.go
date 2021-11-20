//go:build generate

package static

//go:generate convert -background none -density 1200 -resize 1920x1080 root/static/map.svg -write root/static/map.png -write root/static/map.webp root/static/map.jpg
