.PHONY: all
all: clean build

.PHONY: build
build:
	makepkg -cfs

.PHONY: clean
clean:
	git restore PKGBUILD
	-rm -rf src/ pkg/ w16-*.pkg.tar.zst

PKGBUILD: w16.service
	updpkgsums
