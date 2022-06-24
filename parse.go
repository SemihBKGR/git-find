package main

import (
	"fmt"
	"strconv"
	"strings"
)

type diff struct {
	filename   string
	content    string
	lineNumber uint
	isAdded    bool
}

func (d *diff) String() string {
	r := '+'
	if !d.isAdded {
		r = '-'
	}
	return fmt.Sprintf("%c %s %d\t|%s", r, d.filename, d.lineNumber, d.content)
}

func parseDiff(s string) ([]*diff, error) {
	diffs := make([]*diff, 0, 10)
	changeFiles := strings.Split(s, "diff --git")
	for _, changeFile := range changeFiles {
		changes := strings.Split(changeFile, "@@")
		filename := parseMetadata(changes[0])
		for i := 1; i < len(changes); i++ {
			lines := strings.Split(changes[i], "\n")
			lineNumber := extractLineNumber(lines[0])
			for j := 1; j < len(lines); j++ {
				line := lines[j]
				lineNumber++
				if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
					diffLine := diff{
						filename:   filename,
						content:    line[1:],
						isAdded:    line[0] == '+',
						lineNumber: lineNumber,
					}
					diffs = append(diffs, &diffLine)
				}
			}
		}
	}
	return diffs, nil
}

func parseMetadata(m string) string {
	for _, l := range strings.Split(m, "\n") {
		if strings.HasPrefix(l, "+++") {
			return l[6:]
		}
	}
	return ""
}

func extractLineNumber(m string) uint {
	if strings.Contains(m, "+") && strings.Contains(m, ",") {
		startIndex := strings.IndexRune(m, '+') + 1
		endIndex := strings.IndexRune(m[startIndex:], ',') + startIndex
		numberStr := m[startIndex:endIndex]
		number, err := strconv.ParseUint(numberStr, 10, 0)
		if err != nil {
			return 0
		}
		return uint(number)
	}
	return 0
}

func find(diffs []*diff, keyword string) []*diff {
	foundDiffs := make([]*diff, 0, 10)
	for _, d := range diffs {
		if strings.Contains(d.content, keyword) {
			foundDiffs = append(foundDiffs, d)
		}
	}
	return foundDiffs
}
