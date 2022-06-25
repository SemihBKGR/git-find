package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	commit     string
	ignoreCase bool
	help       bool
	removed    bool
	regex      bool
)

func init() {

	if testsRunning() {
		return
	}

	flag.StringVar(&commit, "commit", "", "the commit on which you want to find")
	flag.BoolVar(&ignoreCase, "ignore-case", false, "ignore case sensitivity")
	flag.BoolVar(&help, "help", false, "print args")
	flag.BoolVar(&removed, "removed", false, "include removed lines")
	flag.BoolVar(&regex, "regex", false, "apply regex")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

}

func find(diffs []*diff, keyword string, ignoreCase bool, removed bool) []*diff {
	if ignoreCase {
		keyword = strings.ToLower(keyword)
	}
	foundDiffs := make([]*diff, 0, 10)
	for _, d := range diffs {
		if ignoreCase {
			d.content = strings.ToLower(d.content)
		}
		if strings.Contains(d.content, keyword) {
			if !removed && !d.isAdded {
				continue
			}
			foundDiffs = append(foundDiffs, d)
		}
	}
	return foundDiffs
}

func findRegex(diffs []*diff, regexp *regexp.Regexp, removed bool) []*diff {
	foundDiffs := make([]*diff, 0, 10)
	for _, d := range diffs {
		if regexp.MatchString(d.content) {
			if !removed && !d.isAdded {
				continue
			}
			foundDiffs = append(foundDiffs, d)
		}
	}
	return foundDiffs
}

func main() {

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

	var foundDiffs = make([]*diff, 0, 10)

	if regex {
		for _, searchTerm := range searchTerms {
			regexp, err := regexp.Compile(searchTerm)
			if err != nil {
				panic(err)
			}
			foundDiffs = append(foundDiffs, findRegex(diffs, regexp, removed)...)
		}
	} else {
		for _, word := range searchTerms {
			foundDiffs = append(foundDiffs, find(diffs, word, ignoreCase, removed)...)
		}
	}

	for _, foundDiff := range foundDiffs {
		fmt.Println(foundDiff)
	}

}

func testsRunning() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
