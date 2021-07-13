package db

type Desc struct {
        Filename string
        Name string
        Version string
        Desc string
        CSize int
        ISize int
        MD5Sum string
        Sha256Sum string
        PgpSig string
        URL string
        License string
        BuildDate int
        Packager string
        Replaces []string
        Conflicts []string
        Provides []string
        Depends []string

		// format_entry "FILENAME"  "${1##*/}"
		// format_entry "NAME"      "$pkgname"
		// format_entry "BASE"      "$pkgbase"
		// format_entry "VERSION"   "$pkgver"
		// format_entry "DESC"      "$pkgdesc"
		// format_entry "GROUPS"    "${_groups[@]}"
		// format_entry "CSIZE"     "$csize"
		// format_entry "ISIZE"     "$size"
  //
		// # add checksums
		// format_entry "MD5SUM"    "$md5sum"
		// format_entry "SHA256SUM" "$sha256sum"
  //
		// # add PGP sig
		// format_entry "PGPSIG"    "$pgpsig"
  //
		// format_entry "URL"       "$url"
		// format_entry "LICENSE"   "${_licenses[@]}"
		// format_entry "ARCH"      "$arch"
		// format_entry "BUILDDATE" "$builddate"
		// format_entry "PACKAGER"  "$packager"
		// format_entry "REPLACES"  "${_replaces[@]}"
		// format_entry "CONFLICTS" "${_conflicts[@]}"
		// format_entry "PROVIDES"  "${_provides[@]}"
  //
		// format_entry "DEPENDS" "${_depends[@]}"
		// format_entry "OPTDEPENDS" "${_optdepends[@]}"
		// format_entry "MAKEDEPENDS" "${_makedepends[@]}"
		// format_entry "CHECKDEPENDS" "${_checkdepends[@]}"
}

type Files struct {
        Files []string
}

func formatEntry(b *strings.Builder, )
