//go:build generate

package internal

//go:generate convert -background none -density 1200 -resize 1920x1080 static/static/map.svg -write static/static/map.png static/static/map.webp
