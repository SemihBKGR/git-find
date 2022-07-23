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

	filesOldName := mapSlice(diff.files, func(df *diffFile) string {
		return df.oldFilename
	})
	assert.Contains(filesOldName, "find.go")
	assert.Contains(filesOldName, "parse.go")
	assert.Contains(filesOldName, "README.md")

	filesNewName := mapSlice(diff.files, func(df *diffFile) string {
		return df.newFilename
	})
	assert.Contains(filesNewName, "find.go")
	assert.Contains(filesNewName, "parse.go")
	assert.Contains(filesNewName, "README.md")

	/*
		filesDiffLines := mapSlice(diff.files, func(df *diffFile) []*diffLine {
			return df.lines
		})
	*/
	//todo: *diffLine assertion

}

func readTestdata(filename string, t *testing.T) string {
	bytes, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

func mapSlice[T, R any](s []T, f func(T) R) []R {
	r := make([]R, len(s))
	for _, i := range s {
		r = append(r, f(i))
	}
	return r
}
