# mono

[@seankhliao][githubseankhliao] is in the **mono**repo phase of life

[![Go Reference][badgepkgsite]][pkgsitemono]
[![MIT LICENSE][badgelicense]][filelicense]

Why have a dozen repos when everything can live together in harmony?
No more changes that span across all the repos when you cange one thing,
it's just one change now.

## layout

The top level should be reserved for global config files,
and code is partitioned primarily by purpose / language.

## projects

see the directories under [`./go/cmd`](cmd/go) for the maintained runnable things.

## other repos

Some code just doesn't suit living in here,
in particular the public repos:

- [seankhliao/config][repoconfig] basically a git repo repo of `~/.config`
- [seankhliao/seankhliao][reposeankhliao] the readme that shows up on my github profile page, has to be its own repo
- [seankhliao/seankhliao.github.io][repogithubio] just a redirect, also has to be its own repo
- [erred/][githuberred] all of my old repos

[badgelicense]: https://img.shields.io/github/license/seankhliao/mono?style=flat-square
[badgepkgsite]: https://pkg.go.dev/badge/go.seankhliao.com/mono.svg
[filelicense]: LICENSE
[githuberred]: https://github.com/erred
[githubseankhliao]: https://github.com/seankhliao
[pkgsitemono]: https://pkg.go.dev/go.seankhliao.com/mono
[repoconfig]: https://github.com/seankhliao/config
[repogithubio]: https://github.com/seankhliao/seankhliao.github.io
[reposeankhliao]: https://github.com/seankhliao/seankhliao
