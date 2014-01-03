// A simple play "shell".
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var chre = regexp.MustCompile(`cd(\s+.+)?`)

var bpath = "" // Your amazing path here.

func getbins() ([]string, error) {
	fh, e := os.Open(bpath)
	if e != nil {
		fmt.Println("Cannot open BIN dir.")
		return nil, e
	}
	defer fh.Close()
	names, e := fh.Readdirnames(0)
	if e != nil {
		fmt.Println("Cannot read BIN dir.")
		return nil, e
	}

	return names, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	bins, e := getbins()
	if e != nil {
		fmt.Println("Error reading bin dir.")
		panic(e)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Print("[" + pwd + "]$ ")
	for scanner.Scan() {
		st := scanner.Text()
		switch {

		case st == "exit":
			os.Exit(0)

		case st == "pwd":
			fmt.Println(pwd)

		case chre.MatchString(st):
			if chre.FindStringSubmatch(st)[1] == "" {
				if err := os.Chdir(os.Getenv("HOME")); err != nil {
					panic(err)
				}
			} else {
				argdir := strings.TrimSpace(chre.FindStringSubmatch(st)[1])
				if m, _ := regexp.MatchString("^/.?", argdir); m == false {
					argdir = pwd + "/" + argdir
				}
				if err := os.Chdir(argdir); err != nil {
					fmt.Println("No such directory!")
				}
			}
		default:
			// Break the args to external programs proper.
			re := regexp.MustCompile(`\s+`)
			cargs := re.Split(st, -1)
			for _, v := range bins {
				if cargs[0] == v {
					var procAttr os.ProcAttr
					procAttr.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
					pid, err := os.StartProcess(bpath+v, cargs, &procAttr)
					if err != nil {
						fmt.Printf("%v", err)
					}
					pid.Wait()
					break
				}
			}
		}
		pwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Print("[" + pwd + "]$ ")

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
