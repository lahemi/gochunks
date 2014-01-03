// This mkdir respects umask, unlike GNU mkdir, when -m is used.
// -m and -p are implemented.
package main

import (
	"os"
	"strconv"
)

func dierr(e string) {
	os.Stderr.Write([]byte(e))
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	var mflag, pflag = false, false
	var mode = uint32(0755)
	var dirs []string

	for _, a := range args {
		switch {
		case a == "-p":
			pflag = true
		case a == "-m":
			mflag = true
		case mflag:
			m, e := strconv.ParseUint(a, 8, 32)
			if e != nil {
				dierr("Improper mode.\n")
			}
			if m > 0777 {
				dierr("Improper mode.\n")
			}
			mode = uint32(m)
			mflag = false
		default:
			dirs = append(dirs, a)
		}
	}

	for _, d := range dirs {
		if pflag {
			// Such fancy helper, saves us maybe 10LOC.
			err := os.MkdirAll(d, os.FileMode(mode))
			if err != nil {
				dierr("Cannot create dir.\n")
			}
		} else {
			err := os.Mkdir(d, os.FileMode(mode))
			if err != nil {
				dierr("Cannot create dir.\n")
			}
		}
	}
}
