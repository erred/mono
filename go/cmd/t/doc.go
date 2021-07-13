// t runs ripgrep and outputs a file to be sourced in the shell
// for quick access to matches.
//
// Example shell function:
//     function t() {
//         command t -i "$@"
//         source /tmp/t_aliases 2>/dev/null
//     }
//
// The matches are grouped by file (files in no particular order)
// and can be accessed by `eN` where `N` is the number in brackets.
// This opens the `nvim` at the specified line.
//     cmd/paste/doc.go
//     [29] 4:9:package main
package main
