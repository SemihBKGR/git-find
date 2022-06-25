package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {

	commit := flag.String("commit", "", "commit")
	flag.Parse()

	var c *exec.Cmd

	if *commit == "" {
		c = exec.Command("git", "--no-pager", "diff")
	} else {
		c = exec.Command("git", "--no-pager", "diff", *commit+"~1", *commit)
	}

	r, err := c.Output()

	if err != nil {
		panic(err)
	}

	s := string(r)

	diffs, err := parseDiff(s)

	if err != nil {
		panic(err)
	}

	foundDiffs := find(diffs, "Diff")

	for _, foundDiff := range foundDiffs {
		if foundDiff.isAdded {
			fmt.Println(foundDiff)
		}
	}

}
