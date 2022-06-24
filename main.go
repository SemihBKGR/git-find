package main

import (
	"fmt"
	"os"
)

func main() {

	//c := exec.Command("git", "--no-pager", "diff", "HEAD~1", "HEAD")
	//
	//r, err := c.Output()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//s := string(r)
	//
	//fmt.Println(s)

	bytes, err := os.ReadFile("testdata/diff.txt")

	if err != nil {
		panic(err)
	}

	diffs, err := parseDiff(string(bytes))

	if err != nil {
		panic(err)
	}
	foundDiffs := find(diffs, "Diff")

	for _, foundDiff := range foundDiffs {
		fmt.Println(foundDiff)
	}

}
