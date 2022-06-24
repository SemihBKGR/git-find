package main

import (
	"strings"
)

type Diff struct {
	Filename   string
	Content    string
	LineNumber uint
	IsAdded    bool
}

func parseDiff(s string) ([]*Diff, error) {
	diffs := make([]*Diff, 0, 10)
	changeFiles := strings.Split(s, "diff --git")
	for _, changeFile := range changeFiles {
		changes := strings.Split(changeFile, "@@")
		filename := parseMetadata(changes[0])
		for i := 1; i < len(changes); i++ {
			lines := strings.Split(changes[i], "\n")
			//todo: calculate line number
			for _, line := range lines {
				if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
					diffLine := Diff{
						Filename: filename,
						Content:  line[1:],
						IsAdded:  line[0] == '+',
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
			return l[7:]
		}
	}
	return ""
}

func find(diffs []*Diff, keyword string) []*Diff {
	foundDiffs := make([]*Diff, 0, 10)
	for _, d := range diffs {
		if strings.Contains(d.Content, keyword) {
			foundDiffs = append(foundDiffs, d)
		}
	}
	return foundDiffs
}
