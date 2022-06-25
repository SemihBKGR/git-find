package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	commit     string
	ignoreCase bool
	help       bool
	removed    bool
	regex      bool
)

func main() {

	flag.StringVar(&commit, "commit", "", "the commit on which you want to findKeyword")
	flag.BoolVar(&ignoreCase, "ignore-case", false, "ignore case sensitivity")
	flag.BoolVar(&help, "help", false, "print args")
	flag.BoolVar(&removed, "removed", false, "include removed lines")
	flag.BoolVar(&regex, "regex", false, "apply regex")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	searchTerms := flag.Args()

	if len(searchTerms) == 0 {
		//todo: warn
		os.Exit(0)
	}

	var c *exec.Cmd

	if commit == "" {
		c = exec.Command("git", "--no-pager", "diff")
	} else {
		c = exec.Command("git", "--no-pager", "diff", commit+"~1", commit)
	}

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	c.Dir = dir

	r, err := c.Output()

	if err != nil {
		panic(err)
	}

	s := string(r)

	diffs, err := parseDiff(s)

	if err != nil {
		panic(err)
	}

	foundDiffs, err := find(diffs, searchTerms, ignoreCase, removed, regex)

	if err != nil {
		panic(err)
	}

	for _, foundDiff := range foundDiffs {
		fmt.Println(foundDiff)
	}

}
