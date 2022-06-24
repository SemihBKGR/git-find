package main

import (
	"fmt"
	"os/exec"
)

func main() {

	c := exec.Command("git", "--no-pager", "diff", "HEAD~1", "HEAD")

	r, err := c.Output()

	if err != nil {
		panic(err)
	}

	s := string(r)

	fmt.Println(s)

}
