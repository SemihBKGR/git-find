package main

import (
	"flag"
	"github.com/gookit/color"
	"os"
	"os/exec"
)

func main() {

	var (
		commit     string
		ignoreCase bool
		help       bool
		removed    bool
	)

	flag.StringVar(&commit, "commit", "", "the commit on which you want to findKeyword")
	flag.BoolVar(&ignoreCase, "ignore-case", false, "ignore case sensitivity")
	flag.BoolVar(&help, "help", false, "print args")
	flag.BoolVar(&removed, "removed", false, "include removed lines")
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
	if len(r.occurrences) == 0 {
		return
	}
	if r.file.newFilename == r.file.oldFilename {
		color.Yellowf("%s\n", r.file.newFilename)
	} else {
		color.Yellowf("%s -> %s\n", r.file.oldFilename, r.file.newFilename)
	}
	for _, lo := range r.occurrences {
		if lo.line.added {
			color.Greenf("+ %d\t|", lo.line.lineNumber)
			printLineByHighlightingOccurrences(lo.line.content, lo, color.Green, color.Magenta)
		} else {
			color.Redf("- %d\t|", lo.line.lineNumber)
			printLineByHighlightingOccurrences(lo.line.content, lo, color.Red, color.Magenta)
		}
	}
}

func printLineByHighlightingOccurrences(s string, lo *lineOccurrence, mc, oc color.Color) {
	occurrenceIndexPairs := make([][2]uint, 0, 10)
	i := 0
	for st, indexes := range lo.occurrences {
		stLen := uint(len(st))
		for _, index := range indexes {
			occurrenceIndexPairs = append(occurrenceIndexPairs, [2]uint{index, index + stLen})
			i++
		}
	}

	for i := 0; i < len(occurrenceIndexPairs); i++ {
		for j := 0; j < len(occurrenceIndexPairs); j++ {
			if occurrenceIndexPairs[i][0] > occurrenceIndexPairs[j][0] {
				pair := occurrenceIndexPairs[j]
				occurrenceIndexPairs[j] = occurrenceIndexPairs[i]
				occurrenceIndexPairs[i] = pair
			}
		}
	}

	index := uint(0)
	for _, p := range occurrenceIndexPairs {
		mc.Print(s[index:p[0]])
		oc.Print(s[p[0]:p[1]])
		index = p[1]
	}
	mc.Println(s[index:])

}
