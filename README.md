# git-find

[![ci workflow](https://github.com/semihbkgr/git-find/actions/workflows/ci.yml/badge.svg)](https://github.com/semihbkgr/git-find/actions/workflows/ci.yml)
[![go doc](https://godoc.org/github.com/semihbkgr/git-find?status.png)](https://pkg.go.dev/github.com/semihbkgr/git-find)

'git-find' is a tool that enables you to search texts on a specific git commit.

```shell
git find [args --commit=<hash> --ignore-case ...] [search terms]
```

to list all available args: 'git-find --help'

### Installation

```shell
go install github.com/semihbkgr/git-find@latest
```

### Usage

to find 'todo' and 'func' keywords on commit '9aaf0b4'

```shell
git find --commit=9aaf0b4 --ignore-case --removed todo func
# --commit: defines which commit you want to search on
# --ignore-case: apply case-insensitive search
# --removed: prints removed lines as well
```

output

![output](./output.png)

to list all available args

```shell
git-find --help
```
