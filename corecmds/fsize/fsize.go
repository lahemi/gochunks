package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: fsize file ...")
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	for _, v := range os.Args[1:] {
		fo, err := os.Stat(v)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stdout, "%v: %v\n", v, fo.Size())
	}
}
