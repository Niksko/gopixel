package main

import (
	"github.com/niksko/gopixel/pkg/cmd"
	"os"
)

func main() {
	filename := os.Args[1]
	ok, _ := cmd.Sort(filename)
	if ok {
		os.Exit(0)
	} else {
		os.Stderr.WriteString("Error sorting file")
		os.Exit(1)
	}
}
