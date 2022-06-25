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
	changedFiles := splitByLinePrefix(s, "diff --git")
	for i := 0; i < len(changedFiles); i++ {
		changedSnippets := splitByLinePrefix(changedFiles[i], "@@")
		filename := parseMetadata(changedSnippets[0])
		for j := 1; j < len(changedSnippets); j++ {
			lines := strings.Split(changedSnippets[j], "\n")
			minusLineNumber, plusLineNumber := extractLineNumber(lines[0])
			for k := 1; k < len(lines); k++ {
				line := lines[k]
				if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
					if line[0] == '-' {
						diffLine := diff{
							filename:   filename,
							content:    line[1:],
							isAdded:    false,
							lineNumber: minusLineNumber,
						}
						diffs = append(diffs, &diffLine)
						minusLineNumber++
					} else {
						diffLine := diff{
							filename:   filename,
							content:    line[1:],
							isAdded:    true,
							lineNumber: plusLineNumber,
						}
						diffs = append(diffs, &diffLine)
						plusLineNumber++
					}
				} else {
					minusLineNumber++
					plusLineNumber++
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

func extractLineNumber(m string) (uint, uint) {
	var minus uint = 0
	var plus uint = 0
	if strings.Contains(m, "-") && strings.Contains(m, ",") {
		startIndex := strings.IndexRune(m, '-') + 1
		endIndex := strings.IndexRune(m[startIndex:], ',') + startIndex
		numberStr := m[startIndex:endIndex]
		number, err := strconv.ParseUint(numberStr, 10, 0)
		if err != nil {
			minus = 0
		}
		minus = uint(number)
	}
	if strings.Contains(m, "+") && strings.Contains(m, ",") {
		startIndex := strings.IndexRune(m, '+') + 1
		endIndex := strings.IndexRune(m[startIndex:], ',') + startIndex
		numberStr := m[startIndex:endIndex]
		number, err := strconv.ParseUint(numberStr, 10, 0)
		if err != nil {
			plus = 0
		}
		plus = uint(number)
	}
	return minus, plus
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
