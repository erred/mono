//go:build generate

package seankhliaocom

//go:generate convert -background none -density 1200 -resize 1920x1080 static/static/map.svg -write static/static/map.png -write static/static/map.webp static/static/map.jpg
