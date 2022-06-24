package main

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
	return nil, nil
}
