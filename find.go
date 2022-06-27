package main

import (
	"regexp"
	"strings"
)

func find(diffs []*diffLine, searchTerms []string, ignoreCase, removed, regex bool) ([]*diffLine, error) {
	foundDiffs := make([]*diffLine, 0, 10)
	for _, searchTerm := range searchTerms {
		if regex {
			r, err := regexp.Compile(searchTerm)
			if err != nil {
				return nil, err
			}
			foundDiffs = append(foundDiffs, findRegex(diffs, r, removed)...)
		} else {
			foundDiffs = append(foundDiffs, findKeyword(diffs, searchTerm, ignoreCase, removed)...)
		}
	}
	return foundDiffs, nil
}

func findKeyword(diffs []*diffLine, keyword string, ignoreCase, removed bool) []*diffLine {
	if ignoreCase {
		keyword = strings.ToLower(keyword)
	}
	foundDiffs := make([]*diffLine, 0, 10)
	for _, d := range diffs {
		lineContent := d.content
		if ignoreCase {
			lineContent = strings.ToLower(lineContent)
		}
		if strings.Contains(lineContent, keyword) {
			if !removed && !d.added {
				continue
			}
			foundDiffs = append(foundDiffs, d)
		}
	}
	return foundDiffs
}

func findRegex(diffs []*diffLine, regexp *regexp.Regexp, removed bool) []*diffLine {
	foundDiffs := make([]*diffLine, 0, 10)
	for _, d := range diffs {
		if regexp.MatchString(d.content) {
			if !removed && !d.added {
				continue
			}
			foundDiffs = append(foundDiffs, d)
		}
	}
	return foundDiffs
}
