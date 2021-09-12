APP := w16
RSYNCHOST := medea

.PHONY: all
all: clean build rsync

.PHONY: build
build:
	updpkgsums
	makepkg -cfs

.PHONY: rsync
rsync:
	rsync -rP *.pkg.tar.zst $(RSYNCHOST):pkgs/

.PHONY: clean
clean:
	-rm -rf src/ pkg/ $(APP)-*.pkg.tar.zst
