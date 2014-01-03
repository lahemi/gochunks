// `echo`, with only the -n flag. Like in Plan9 echo.
package main

import (
	"os"
)

func main() {
	args := os.Args[1:]

	nflag := false
	if args[0] == "-n" {
		nflag = true
		args = args[1:]
	}

	for k, v := range args {
		os.Stdout.Write([]byte(v))
		if k < len(args)-1 {
			os.Stdout.Write([]byte(" "))
		}
	}
	if !nflag {
		os.Stdout.Write([]byte("\n"))
	}
}
