package main

import (
	"flag"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"strings"
)

func main() {

	var (
		commit     string
		ignoreCase bool
		help       bool
		removed    bool
		regex      bool
	)

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
		//todo: warn or print all changes
		os.Exit(0)
	}

	diffOutput, err := gitDiff(commit)
	if err != nil {
		panic(err)
	}

	diff, err := parseDiff(diffOutput)
	if err != nil {
		panic(err)
	}

	fo := findOptions{
		ignoreCase: ignoreCase,
		removed:    removed,
		regex:      regex,
	}

	c := find(diff, searchTerms, fo)
	done := false

	for !done {
		select {
		case r, ok := <-c:
			done = !ok
			if ok {
				printFindResult(r)
			}
		}
	}

}

func gitDiff(commit string) (string, error) {
	c := exec.Command("git", "--no-pager", "diff")
	if commit != "" {
		c.Args = append(c.Args, commit+"~1", commit)
	}
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	c.Dir = dir
	output, err := c.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func printFindResult(r findResult) {
	if len(r.lines) == 0 {
		return
	}
	if r.file.newFilename == r.file.oldFilename {
		color.Yellowf("%s\n", r.file.newFilename)
	} else {
		color.Yellowf("%s -> %s\n", r.file.oldFilename, r.file.newFilename)
	}
	for _, l := range r.lines {
		if l.added {
			color.Greenf("+ %d\t|", l.lineNumber)
			printLineByOccurrences(l.content, r.searchTerm, color.Green, color.Magenta)
		} else {
			color.Redf("- %d\t|", l.lineNumber)
			printLineByOccurrences(l.content, r.searchTerm, color.Red, color.Magenta)
		}
	}
}

func printLineByOccurrences(s, o string, mc, oc color.Color) {
	i := strings.Index(s, o)
	if i == -1 {
		mc.Println(s)
		return
	}
	mc.Print(s[0:i])
	oc.Print(o)
	mc.Println(s[i+len(o):])
}
