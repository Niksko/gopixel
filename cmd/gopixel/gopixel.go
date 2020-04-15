package main

import (
	"github.com/niksko/gopixel/pkg/cmd"
	"os"
	"strconv"
)

func usage() {
	os.Stdout.WriteString("Usage: gopixel <input-file> [sort-angle] > <output-file>\n")
}

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Invalid number of arguments\n")
		usage()
		os.Exit(1)
	}
	filename := os.Args[1]
	var sortAngle int
	var err error
	if len(os.Args) == 3 {
		sortAngle, err = strconv.Atoi(os.Args[2])
		if err != nil {
			os.Stderr.WriteString("Invalid angle")
			os.Exit(1)
		}
	} else {
		sortAngle = 270
	}
	ok, _ := cmd.Sort(filename, uint(sortAngle))
	if ok {
		os.Exit(0)
	} else {
		os.Stderr.WriteString("Error sorting file")
		os.Exit(1)
	}
}
