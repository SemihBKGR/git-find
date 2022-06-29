# git-find

[![ci workflow](https://github.com/semihbkgr/git-find/actions/workflows/ci.yml/badge.svg)](https://github.com/semihbkgr/git-find/actions/workflows/ci.yml)
[![go doc](https://godoc.org/github.com/semihbkgr/git-find?status.png)](https://pkg.go.dev/github.com/semihbkgr/git-find)

'git-find' is a tool that enables you to search on a specific git commit.

Installation

```shell
go install github.com/semihbkgr/git-find@v1.0.0
```

Finds 'todo' and 'func' keywords on commit 9aaf0b4

```shell
git-find --commit=9aaf0b4 --ignore-case --removed todo func
```
Lists all available options

```shell
git-find --help
```
