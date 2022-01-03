// authd is an implementation of envoy's gRPC ext auth v3 service.
// It reads configuration from a prototext file.
package main

import _ "embed"

//go:embed doc.go
var docgo string
