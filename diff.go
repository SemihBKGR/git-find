package main

import (
	"strconv"
	"strings"
)

type diff struct {
	files []*diffFile
}

func newDiff() *diff {
	return &diff{
		files: make([]*diffFile, 0, 10),
	}
}

func parseDiff(s string) (*diff, error) {
	diff := newDiff()
	changedFiles := splitByLinePrefix(s, "diff --git")
	for i := 0; i < len(changedFiles); i++ {
		changedSnippets := splitByLinePrefix(changedFiles[i], "@@")
		newFilename, oldFilename := extractMetadata(changedSnippets[0])
		diffFile := newDiffFile()
		diffFile.newFilename = newFilename
		diffFile.oldFilename = oldFilename
		for j := 1; j < len(changedSnippets); j++ {
			lines := strings.Split(changedSnippets[j], "\n")
			addLineNumber, removeLineNumber := extractLineNumber(lines[0])
			for k := 1; k < len(lines); k++ {
				line := lines[k]
				if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
					added := line[0] == '+'
					var lineNumber uint
					if added {
						lineNumber = addLineNumber
						addLineNumber++
					} else {
						lineNumber = removeLineNumber
						removeLineNumber++
					}
					diffLine := diffLine{
						content:    line[1:],
						added:      added,
						lineNumber: lineNumber,
					}
					diffFile.lines = append(diffFile.lines, &diffLine)
				} else {
					removeLineNumber++
					addLineNumber++
				}
			}
		}
		diff.files = append(diff.files, diffFile)
	}
	return diff, nil
}

type diffFile struct {
	newFilename string
	oldFilename string
	lines       []*diffLine
}

func newDiffFile() *diffFile {
	return &diffFile{
		lines: make([]*diffLine, 0, 10),
	}
}

type diffLine struct {
	content    string
	lineNumber uint
	added      bool
}

func extractMetadata(s string) (string, string) {
	newFilename := ""
	oldFilename := ""
	for _, l := range strings.Split(s, "\n") {
		if strings.HasPrefix(l, "+++") {
			i := strings.IndexRune(l, '/')
			if filename := l[i+1:]; filename != "dev/null" {
				newFilename = filename
			}
		} else if strings.HasPrefix(l, "---") {
			i := strings.IndexRune(l, '/')
			if filename := l[i+1:]; filename != "dev/null" {
				oldFilename = filename
			}
		}
	}
	return newFilename, oldFilename
}

func extractLineNumber(s string) (uint, uint) {
	var add uint = 0
	var remove uint = 0
	if strings.Contains(s, "+") && strings.Contains(s, ",") {
		startIndex := strings.IndexRune(s, '+') + 1
		endIndex := strings.IndexRune(s[startIndex:], ',') + startIndex
		if endIndex < startIndex {
			endIndex = strings.IndexRune(s[startIndex:], ' ') + startIndex
		}
		numberStr := s[startIndex:endIndex]
		number, err := strconv.ParseUint(numberStr, 10, 0)
		if err != nil {
			add = uint(number)
		}
	}
	if strings.Contains(s, "-") && strings.Contains(s, ",") {
		startIndex := strings.IndexRune(s, '-') + 1
		endIndex := strings.IndexRune(s[startIndex:], ',') + startIndex
		if endIndex < startIndex {
			endIndex = strings.IndexRune(s[startIndex:], ' ') + startIndex
		}
		numberStr := s[startIndex:endIndex]
		number, err := strconv.ParseUint(numberStr, 10, 0)
		if err != nil {
			remove = uint(number)
		}
	}
	return add, remove
}

func splitByLinePrefix(s, p string) []string {
	split := make([]string, 0, 10)
	lines := strings.Split(s, "\n")
	sb := strings.Builder{}
	for i, line := range lines {
		if strings.HasPrefix(line, p) && i != 0 {
			split = append(split, sb.String())
			sb = strings.Builder{}
		}
		sb.WriteString(line + "\n")
	}
	return append(split, sb.String())
}
