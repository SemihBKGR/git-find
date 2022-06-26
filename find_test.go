package main

import (
	"os"
	"testing"
)

func TestFind(t *testing.T) {

	bytes, err := os.ReadFile("testdata/9aaf0b4.diff")

	if err != nil {
		t.Fatal(err)
	}

	diffs, err := parseDiff(string(bytes))

	if err != nil {
		t.Fatal(err)
	}

	expectedFoundDiffsLen := 3

	foundDiffs, err := find(diffs, []string{"todo", "func"}, true, true, false)

	if err != nil {
		t.Fatal(err)
	}

	if len(foundDiffs) != expectedFoundDiffsLen {
		t.Fatalf("expected len of found diffs is %d, but actual value is %d", expectedFoundDiffsLen, len(foundDiffs))
	}

}
