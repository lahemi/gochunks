// `p`, paginate, "toned down" implementation.
// Displays given files to stdout, 22 lines at a time
// by default, which can be changed with -[number] flag
// when starting the program. After each displayed block
// the program waits for input from the user, either a
// newline or q for quiting. This program does _not_
// have the ! command for passing stuff to the shell.
package main

import (
	"bufio"
	"os"
	"strconv"
)

var defpglen = 22

// Could be better; keep a temp buffer, which you empty after each
// "stop", meaning there'd be no need to worry about large files.
func printfile(file string, pglen int) {
	f, e := os.Open(file)
	if e != nil {
		panic(e)
	}
	defer f.Close()

	var lines []string
	b := bufio.NewReader(f)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			break
		}
		lines = append(lines, line)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for k, v := range lines {
		if k != 0 {
			if (k % pglen) == 0 {
				for scanner.Scan() {
					st := scanner.Text()
					if st == "q" {
						os.Exit(0)
					}
					if st == "" {
						break
					}
				}
			}
		}
		os.Stdout.Write([]byte(v))
	}
}

func main() {
	args := os.Args[1:]
	pglen := defpglen

	for i := 0; i < len(args); i++ {
		if args[i][0] == '-' {
			s, e := strconv.Atoi(args[i][1:])
			if e != nil {
				panic(e)
			}
			pglen = s

		} else {
			printfile(args[i], pglen)
		}
	}
}
