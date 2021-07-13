// mirrorrank is a NIH version of reflector
// for ranking Arch Linux / pacman repository mirrors.
//
// mirrorrank downloads the filtered (via flags) mirrorlist
// from archlinux.org/mirrorlist and times the download time of
// community.db.
// Downloads are in parallel (5 at a time) and times out after 5 seconds.
//
// By default it writes directly to the mirrorlist read by pacman.
package main
