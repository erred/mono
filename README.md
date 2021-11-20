# mono

[@seankhliao][githubseankhliao] is in the **mono**repo phase of life

[![Go Reference][badgepkgsite]][pkgsitemono]
[![MIT LICENSE][badgelicense]](LICENSE)

Why have a dozen repos when everything can live together in harmony?
No more changes that span across all the repos when you cange one thing,
it's just one change now.

## secrets

secrets are hidden through a git filter:

```txt
[filter "ageencrypt"]
	clean    = age -r age14mg08panez45c6lj2cut2l8nqja0k5vm2vxmv5zvc4ufqgptgy2qcjfmuu -a -
	smudge   = age -d -i ~/.ssh/age.key -
	required = true
```

## other repos

Some code just doesn't suit living in here,
in particular the public repos:

- [seankhliao/config][repoconfig] basically a git repo repo of `~/.config`
- [seankhliao/seankhliao][reposeankhliao] the readme that shows up on my github profile page, has to be its own repo
- [seankhliao/seankhliao.github.io][repogithubio] just a redirect, also has to be its own repo
- [erred/][githuberred] all of my old repos

[badgelicense]: https://img.shields.io/github/license/seankhliao/mono?style=flat-square
[badgepkgsite]: https://pkg.go.dev/badge/go.seankhliao.com/mono.svg
[githuberred]: https://github.com/erred
[githubseankhliao]: https://github.com/seankhliao
[pkgsitemono]: https://pkg.go.dev/go.seankhliao.com/mono
[repoconfig]: https://github.com/seankhliao/config
[repogithubio]: https://github.com/seankhliao/seankhliao.github.io
[reposeankhliao]: https://github.com/seankhliao/seankhliao
[seankhliaocom]: https://seankhliao.com/?utm_source=github&utm_medium=mono
