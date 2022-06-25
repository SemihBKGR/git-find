package main

import (
	"os"
	"testing"
)

func TestFind(t *testing.T) {

	bytes, err := os.ReadFile("testdata/diff.txt")

	if err != nil {
		t.Fatal(err)
	}

	diffs, err := parseDiff(string(bytes))

	if err != nil {
		t.Fatal(err)
	}

	expectedFoundDiffsLen := 6

	foundDiffs := find(diffs, "Diff", false, true)

	if len(foundDiffs) != expectedFoundDiffsLen {
		t.Fatalf("expected len of found diffs is %d, but actual value is %d", expectedFoundDiffsLen, len(foundDiffs))
	}

}
