package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseDiff(t *testing.T) {
	assert := assert.New(t)

	output := readTestdata("9aaf0b4.diff", t)

	diff, err := parseDiff(output)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(3, len(diff.files))
}

func TestSplitByLinePrefix(t *testing.T) {
	assert := assert.New(t)

	output := readTestdata("9aaf0b4.diff", t)

	splits := splitByLinePrefix(output, "diff --git")

	assert.Equal(3, len(splits))
}

func readTestdata(filename string, t *testing.T) string {
	bytes, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}
