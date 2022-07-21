package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"os"
	"os/exec"
	"strings"
)

var (
	commit     = flag.String("commit", "", "the commit on which you want to findKeyword")
	ignoreCase = flag.Bool("ignore-case", false, "ignore case sensitivity")
	removed    = flag.Bool("removed", false, "include removed lines")
)

func main() {
	flag.Parse()

	args := notBlankStrings(flag.Args())
	searchTerms := deduplicate(args, *ignoreCase)
	if len(searchTerms) == 0 {
		fmt.Fprintln(os.Stderr, "missing search terms")
		os.Exit(1)
	}

	diffOutput, err := gitDiff(*commit)
	if err != nil {
		panic(err)
	}

	diff, err := parseDiff(diffOutput)
	if err != nil {
		panic(err)
	}

	fo := findOptions{
		ignoreCase: *ignoreCase,
		removed:    *removed,
	}

	c := find(diff, searchTerms, fo)
	done := false

	for !done {
		r, ok := <-c
		done = !ok
		if ok {
			printFindResult(r)
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
	if len(r.occurrences) == 0 {
		return
	}
	if r.file.newFilename == r.file.oldFilename {
		color.Bluef("%s\n", r.file.newFilename)
	} else {
		if r.file.newFilename == "" {
			color.Bluef("[old] %s\n", r.file.oldFilename)
		} else if r.file.oldFilename == "" {
			color.Bluef("[new] %s\n", r.file.newFilename)
		} else {
			color.Bluef("%s -> %s\n", r.file.oldFilename, r.file.newFilename)
		}
	}
	for _, lo := range r.occurrences {
		if lo.line.added {
			color.Greenf("+ %d\t|", lo.line.lineNumber)
			printLineByHighlightingOccurrences(lo.line.content, lo.occurrences, color.Green, color.Yellow)
		} else if *removed {
			color.Redf("- %d\t|", lo.line.lineNumber)
			printLineByHighlightingOccurrences(lo.line.content, lo.occurrences, color.Red, color.Yellow)
		}
	}
}

func printLineByHighlightingOccurrences(s string, lo map[string][]uint, pc, sc color.Color) {
	occurrenceIndexPairs := make([][2]uint, 0, 10)
	i := 0
	for st, indexes := range lo {
		stLen := uint(len(st))
		for _, index := range indexes {
			occurrenceIndexPairs = append(occurrenceIndexPairs, [2]uint{index, index + stLen})
			i++
		}
	}
	for i := 0; i < len(occurrenceIndexPairs); i++ {
		for j := 0; j < len(occurrenceIndexPairs); j++ {
			if occurrenceIndexPairs[i][0] < occurrenceIndexPairs[j][0] {
				pair := occurrenceIndexPairs[j]
				occurrenceIndexPairs[j] = occurrenceIndexPairs[i]
				occurrenceIndexPairs[i] = pair
			}
		}
	}
	index := uint(0)
	for _, p := range occurrenceIndexPairs {
		if index >= p[1] {
			continue
		}
		if index >= p[0] {
			sc.Print(s[index:p[1]])
		} else {
			pc.Print(s[index:p[0]])
			sc.Print(s[p[0]:p[1]])
		}
		index = p[1]
	}
	pc.Println(s[index:])
}

func deduplicate(ss []string, ignoreCase bool) []string {
	m := make(map[string]any)
	for _, s := range ss {
		if ignoreCase {
			s = strings.ToLower(s)
		}
		m[s] = nil
	}
	uss := make([]string, 0, len(m))
	for s := range m {
		uss = append(uss, s)
	}
	return uss
}

func notBlankStrings(ss []string) []string {
	ness := make([]string, 0, len(ss))
	for _, s := range ss {
		if len(strings.TrimSpace(s)) != 0 {
			ness = append(ness, s)
		}
	}
	return ness
}
