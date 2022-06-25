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
	search     string
)

func init() {

	if testsRunning() {
		return
	}

	flag.StringVar(&commit, "commit", "", "commit on which you want to find")
	flag.BoolVar(&ignoreCase, "ignore-case", false, "case sensitivity")
	flag.BoolVar(&help, "help", false, "print args")
	flag.BoolVar(&removed, "removed", false, "include removed lines")
	flag.BoolVar(&regex, "regex", false, "apply regex")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	search = flag.Arg(0)

	if search == "" {
		fmt.Fprintln(os.Stderr, "Search var cannot be empty")
		os.Exit(1)
	}

	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("ignoreCase: %v\n", ignoreCase)
	fmt.Printf("help: %v\n", help)
	fmt.Printf("removed: %v\n", removed)
	fmt.Printf("regex: %v\n", regex)
	fmt.Printf("search: %s\n", search)

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

	//

	var foundDiffs []*diff

	if regex {
		regexp, err := regexp.Compile(search)
		if err != nil {
			panic(err)
		}
		foundDiffs = findRegex(diffs, regexp, removed)
	} else {
		foundDiffs = find(diffs, flag.Arg(0), ignoreCase, removed)
	}

	for _, foundDiff := range foundDiffs {
		fmt.Println(foundDiff)
	}

}

func testsRunning() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
