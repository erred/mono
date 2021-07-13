// vanity is a dumb server for serving Go custom / vanity import paths.
//
// vanity serves a page with the `go-import` meta tag
// https://pkg.go.dev/cmd/go#hdr-Remote_import_paths
// and `go-source` meta tags
// https://github.com/golang/gddo/wiki/Source-Code-Links
// to allow go to resolve custom import paths.
//
// Currently it just passes the first path element to the template,
// which assumes Github url format.
package main
