package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"
)

var showhidden = true
var listinfo = false
var fuutimeformat = regexp.MustCompile(`(\s)\s+`)

func dierr(errmsg string) {
	fmt.Fprintln(os.Stderr, errmsg)
	os.Exit(1)
}
func out(output string) {
	fmt.Fprint(os.Stdout, output)
}

func lfun(names []string, path string) {
	var s []string // data slice
	// Collect and stringify all the ls -l fields nicely.
	for _, v := range names {
		fo, _ := os.Stat(path + "/" + v)
		s = append(s, fmt.Sprintf("%v %v %v %v %v %v",
			fo.Mode(), fo.Sys().(*syscall.Stat_t).Uid,
			fo.Sys().(*syscall.Stat_t).Gid, fo.Size(),
			// time.Format messes up everything, temp "fix". Need to check if you
			// can force more "raw" output from it using layouts or something.
			fuutimeformat.ReplaceAllString(fo.ModTime().Format(time.Stamp), "$1"),
			fo.Name()))
	}

	// Doesn't work with just a []int.
	var lgests = make(map[int]int)
	for _, v := range s {
		ss := strings.SplitAfterN(v, " ", 8)
		for i, vv := range ss {
			l := len(strings.TrimSpace(vv)) // Tidy up the split.
			if lgests[i] < l {
				lgests[i] = l
			}
		}
	}
	// Need to have these two separate, since we need to know the longest
	// element out of all the entries before we can determine what to do.
	for _, v := range s {
		ss := strings.SplitAfterN(v, " ", 8)
		for i, vv := range ss {
			vv = strings.TrimSpace(vv)
			if len(vv) >= lgests[i] {
				out(vv)
			} else {
				spca := lgests[i] - len(vv)
				out(vv + strings.Repeat(" ", spca))
			}
			out(" ")
		}
		out("\n")
	}
}

func main() {
	args := os.Args[1:]
	var oargs []string
	var pargs = make(map[string]bool)
	var argdir string

	for _, a := range args {
		if a[0] == '-' {
			pargs[a[1:]] = true
		} else {
			oargs = append(oargs, a)
		}
	}

	pwd, err := os.Getwd()
	if err != nil {
		dierr("Cannot use ls in current directory.")
	}

	// We silently ignore any extra input.
	switch {
	case len(oargs) == 0:
		argdir = pwd
	case oargs[0][0] != '/':
		argdir = pwd + "/" + oargs[0]
	default:
		argdir = oargs[0]
	}

	fh, err := os.Open(argdir)
	if err != nil {
		dierr("Cannot access file or directory.")
	}
	defer fh.Close()

	names, err := fh.Readdirnames(0)
	if err != nil {
		dierr("Cannot read directory or file!")
	}
	// Maybe write a case-insensitive sort?
	sort.Strings(names)

	// If no hidden files flag, then remove them from filenames.
	if !pargs["a"] {
		var anames []string
		for _, v := range names {
			if v[0] != '.' {
				anames = append(anames, v)
			}
		}
		names = anames
	}
	if pargs["l"] {
		lfun(names, argdir)
	} else {
		for _, v := range names {
			fmt.Fprintln(os.Stdout, v)
		}
	}
}
