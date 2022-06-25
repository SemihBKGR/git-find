package main

import (
	"regexp"
	"strings"
)

func find(diffs []*diff, searchTerms []string, ignoreCase, removed, regex bool) ([]*diff, error) {
	foundDiffs := make([]*diff, 0, 10)
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

func findKeyword(diffs []*diff, keyword string, ignoreCase, removed bool) []*diff {
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
