// `cat`, simple cat, like in Plan9, with no evil bloat.
package main

import (
	"os"
)

func cat(file *os.File) {
	for {
		buf := make([]byte, 8192)
		_, err := file.Read(buf)
		if err != nil {
			break
		}
		os.Stdout.Write(buf)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cat(os.Stdin)
	} else {
		for _, file := range args {
			f, e := os.Open(file)
			if e != nil {
				panic(e)
			}
			cat(f)
			f.Close()
		}
	}
}
