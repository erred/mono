// preqsink logs all the requests it receives to the /watch endpoint.
// Note intermediary proxies should not limit response timeout.
package main

import _ "embed"

//go:embed doc.go
var docgo string
