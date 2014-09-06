package main

import (
	"bufio"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	writechan = make(chan string, 512)
	linerex   = regexp.MustCompile("^:[^ ]+?!([^ ]+? ){3}:.+")

	curTarget string

	irccmdlist = []string{
		// PRIVMSG,
		"QUIT",
		"JOIN",
		"PART",
		"NICK",
		"USER",
	}
)

func stdout(s string) {
	os.Stdout.Write([]byte(s + "\n"))
}
func stderr(s string) {
	os.Stderr.Write([]byte(s + "\n"))
}

func handleOut(s string) {
	if linerex.MatchString(s) {
		spl := strings.SplitN(s, " ", 4)
		nick := spl[0][1:strings.Index(s, "!")]
		// cmd := spl[1]
		target := spl[2]
		msg := spl[3]

		sep := " | "

		stdout(nick + sep + target + sep + msg)

	} else {
		stdout(s)
	}
}

func handleCmds(s string) bool {
	cmd := s[:regexp.MustCompile("( |$)").FindStringIndex(s)[0]] // oh my
	for _, k := range irccmdlist {
		switch cmd {
		case "QUIT":
			writechan <- "QUIT"
			os.Exit(0)
		case "JOIN":
			curTarget = s[5:] // Need a target list.
			fallthrough
		case k:
			return true
		}
	}
	return false
}

func main() {

	var (
		server = "irc.freenode.net"
		port   = "6667"
		nick   = "gecannels"
	)

	conn, err := net.Dial("tcp", server+":"+port)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)
	in := bufio.NewReader(os.Stdin)

	go func() {
		for {
			str, err := r.ReadString(byte('\n'))
			if err != nil {
				stderr("read error")
				break
			}
			if str[:4] == "PING" {
				writechan <- "PONG" + str[4:len(str)-2]
			} else {
				handleOut(str[:len(str)-2])
			}
		}
	}()
	go func() {
		for {
			str := <-writechan
			if _, err := w.WriteString(str + "\r\n"); err != nil {
				stderr("write error")
				break
			}
			w.Flush()
		}
	}()

	writechan <- "USER " + nick + " * * :" + nick
	writechan <- "NICK " + nick

	for {
		input, err := in.ReadString('\n')
		if err != nil {
			stderr("error input")
			break
		}
		inp := input[:len(input)-1]
		if handleCmds(inp) {
			writechan <- inp
		} else {
			// Default to PRIVMSG
			writechan <- "PRIVMSG " + curTarget + " :" + inp
		}
	}
}

// Add a target list.
// Direct IRC output plus user input into files.
