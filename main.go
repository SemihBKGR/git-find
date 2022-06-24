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

	s := string(bytes)

	fmt.Println(s)

	diffs, err := parseDiff(s)

	if err != nil {
		panic(err)
	}

	fmt.Printf("diffs len: %d\n", len(diffs))

	fmt.Println("find result")

	lines := find(diffs, "Diff")

	for _, line := range lines {
		fmt.Printf("found line: %s\n", line.Content)
	}

}
