// feedagg is a server that serves aggregates RSS/Atom feeds.
//
// It takes a config file in yaml:
//     feeds:
//         $aggregated2:
//             refresh: 1h # feed refresh frequency in time.Duration format
//             urls:
//                 $feed1: https://...
//                 $feed2: https://...
//         $aggregated2:
//              refresh: 24h
//              urls:
//                 $feed3: https://...
//
// Data is stored as a SQLite database.
// It should support ETags.
package main
