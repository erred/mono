APP := singlepage-medea
RSYNCHOST := medea

.PHONY: all
all: clean build rsync

.PHONY: build
build:
	makepkg -cfs

.PHONY: rsync
rsync:
	rsync -rP *.pkg.tar.zst $(RSYNCHOST):pkgs/

.PHONY: clean
clean:
	-rm -rf src/ pkg/ $(APP)-*.pkg.tar.zst
