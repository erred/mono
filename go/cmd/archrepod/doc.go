// archrepod is a http server for hosting Arch Linux / pacman repositories.
//
// Packages are served at: `/repos/$repo/$arch/$pkg`
// and can be added by POST directly to the location
// it should be served at, with repositories being created on demand.
// Repositories, arch, and packages can be removed through DELETE;
// the database file is updated automatically in both cases.
//
// POST and DELETE operations require a `authorization: Bearer $token` header,
// valid tokens should be passed as a comma separated list to the server
// via the `AR_TOKENS` env.
package main
