# git-find

[![ci workflow](https://github.com/semihbkgr/git-find/actions/workflows/ci.yml/badge.svg)](https://github.com/semihbkgr/git-find/actions/workflows/ci.yml)
[![go doc](https://godoc.org/github.com/semihbkgr/git-find?status.png)](https://pkg.go.dev/github.com/semihbkgr/git-find)

'git-find' is a tool that enables you to search texts on a specific git commit.

Installation

```shell
go install github.com/semihbkgr/git-find@v1.1.0
```

Examples

to find 'todo' and 'func' keywords on commit '9aaf0b4'

```shell
git-find --commit=9aaf0b4 --ignore-case --removed todo func
# 'todo' and 'funcs' are search terms
# --commit: define on which commit you want to search
# --ignore-case: ignore cases in search terms
# --removed: includes removed lines
```
to list all available options

```shell
git-find --help
```
