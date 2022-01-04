// ghdefaults is a github app that sets default settings on new repos,
// either on creation or ownership change
package main

import _ "embed"

//go:embed doc.go
var docgo string
