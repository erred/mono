// earbug runs a background process that logs a user's spotify listening history
// to an etcd instance.
package main

import _ "embed"

//go:embed doc.go
var docgo string
