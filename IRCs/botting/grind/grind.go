package grind

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"
)

var dataPath = os.Getenv("HOME") + "/.local/share/gestelle/"

func BotPrint(conn net.Conn, channel, output string) {
	wr := textproto.NewWriter(bufio.NewWriter(conn))
	wr.PrintfLine("PRIVMSG %s :%s", channel, output)
	// It's nice to see what the Bot sends, too.
	os.Stdout.Write([]byte(channel + " : SELF | " + output + "\n"))
}

func GetPage(url string) (page string) {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	cont, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(cont)
}

func getLines(file string) []string {
	lines, err := ioutil.ReadFile(dataPath + file)
	if err != nil {
		panic(err)
	}
	// ReadFile returns []byte, so this is needed.
	return strings.Split(string(lines), "\n")
}

// Processing and acting upon bot commands.
func ProcessCmds(conn net.Conn, channel, nick, cmd string) {
	rmap := map[string]*regexp.Regexp{
		"fort": regexp.MustCompile("(?i)^fortune"),
	}
	switch {
	case rmap["fort"].MatchString(cmd):
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		f := getLines("fortunes.txt")
		out := f[r.Intn(len(f)-1)]
		BotPrint(conn, channel, out)
	}
}
