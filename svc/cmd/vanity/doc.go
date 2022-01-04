// vanity is a custom go import path redirector.
package main

import _ "embed"

//go:embed doc.go
var docgo string
