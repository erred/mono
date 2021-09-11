APP := singlepage-medea

.PHONY: all
all: clean build

.PHONY: build
build:
	makepkg -cfs

.PHONY: clean
clean:
	git restore PKGBUILD
	-rm -rf src/ pkg/ $(APP)-*.pkg.tar.zst

PKGBUILD: w16.service
	updpkgsums
