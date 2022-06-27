package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFind(t *testing.T) {
	assert := assert.New(t)

	output := readTestdata("9aaf0b4.diff", t)

	diffs, err := parseDiff(output)
	if err != nil {
		t.Fatal(err)
	}

	diffLines := make([]*diffLine, 0, 10)
	for _, diffFile := range diffs.files {
		diffLines = append(diffLines, diffFile.lines...)
	}

	foundDiffs, err := find(diffLines, []string{"todo", "func"}, true, true, false)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(3, len(foundDiffs))

}
