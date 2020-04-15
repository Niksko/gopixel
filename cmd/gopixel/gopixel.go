package main

import (
	"github.com/niksko/gopixel/pkg/cmd"
	"os"
)

func usage() {
	os.Stdout.WriteString("Usage: gopixel <input-file> > <output-file>\n")
}

func main() {
	if len(os.Args) != 2 {
		os.Stderr.WriteString("Invalid number of arguments\n")
		usage()
		os.Exit(1)
	}
	filename := os.Args[1]
	ok, _ := cmd.Sort(filename)
	if ok {
		os.Exit(0)
	} else {
		os.Stderr.WriteString("Error sorting file")
		os.Exit(1)
	}
}
