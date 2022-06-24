package main

import "strings"

type Diff struct {
	StartCommit  string
	EndCommit    string
	ChangedFiles []DiffFile
}

type DiffFile struct {
	Filename     string
	ChangesLines []DiffLine
}

type DiffLine struct {
	IsAdded    bool
	Content    string
	LineNumber uint
}

func parseDiff(s string) (*Diff, error) {
	diff := Diff{}
	changeFiles := strings.Split(s, "diff --git")
	diff.ChangedFiles = make([]DiffFile, 0, len(changeFiles))
	for _, changeFile := range changeFiles {
		diffFile := DiffFile{}
		changes := strings.Split(changeFile, "@@")

		fileMetadata := changes[0]
		filename := parseMetadata(fileMetadata)
		diffFile.Filename = filename

		diffFile.ChangesLines = make([]DiffLine, 0, len(changes)-1)
		for i := 1; i < len(changes); i++ {
			lines := strings.Split(changes[i], "\n")
			//todo: calculate line number
			for _, line := range lines {
				if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
					diffLine := DiffLine{
						Content: line[1:],
						IsAdded: line[0] == '+',
					}
					diffFile.ChangesLines = append(diffFile.ChangesLines, diffLine)
				}
			}
		}
		diff.ChangedFiles = append(diff.ChangedFiles, diffFile)
	}
	return &diff, nil
}

func parseMetadata(m string) string {
	for _, l := range strings.Split(m, "\n") {
		if strings.HasPrefix(l, "+++") {
			return l[7:]
		}
	}
	return ""
}
