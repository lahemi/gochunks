package main

import (
	"bufio"
	"encoding/json"
	"gestelle/grind"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strings"
)

type Conf struct {
	// port is string for convenience
	Server, Port, Nick string
	Channels           []string
}

var botConf = readConf()

var configPath = os.Getenv("HOME") + "/.config/gestelle/"
var dataPath = os.Getenv("HOME") + "/.local/share/gestelle/"

var regs = map[string]*regexp.Regexp{
	"nick": regexp.MustCompile("^:(.+?)!"),
	"ping": regexp.MustCompile("^PING :.+"),
	"chan": regexp.MustCompile("PRIVMSG (#.+?) :"),
	"line": regexp.MustCompile("PRIVMSG #.+? :(.+)"),
	"url":  regexp.MustCompile("(?i)<title>(.+?)</title>"),
	"cmd":  regexp.MustCompile("^ı(.+)"),
}

func readConf() (conf Conf) {
	f, e := ioutil.ReadFile(configPath + "conf.json")
	if e != nil {
		panic(e)
	}
	err := json.Unmarshal(f, &conf)
	if err != nil {
		panic(err)
	}
	return
}

func parseLine(rawline string) (channel, nick, line string) {
	// The actual capture is the 2nd element here, 1st is the whole line.
	channel = regs["chan"].FindStringSubmatch(rawline)[1]
	nick = regs["nick"].FindStringSubmatch(rawline)[1]
	line = regs["line"].FindStringSubmatch(rawline)[1]
	return
}

func parseURL_getTitle(line string) (ret []string) {
	words := strings.Split(line, " ")
	for _, w := range words {
		if m, _ := regexp.MatchString("^https?://.+", w); m {
			p := grind.GetPage(w)
			if p == "" {
				continue
			}
			t := regs["url"].FindStringSubmatch(p)[1]
			ret = append(ret, t)
		}
	}
	return
}

func connect() (conn net.Conn) {
	conn, err := net.Dial("tcp", botConf.Server+":"+botConf.Port)
	if err != nil {
		log.Fatal("Unable to connect ", err)
	}

	wr := textproto.NewWriter(bufio.NewWriter(conn))

	wr.PrintfLine("USER %s 8 * :%s", botConf.Nick, botConf.Nick)
	wr.PrintfLine("NICK %s", botConf.Nick)
	for _, ch := range botConf.Channels {
		wr.PrintfLine("JOIN %s", ch)
	}

	return conn
}

func read(conn net.Conn) {
	wr := textproto.NewWriter(bufio.NewWriter(conn))
	tp := textproto.NewReader(bufio.NewReader(conn))
	for {
		rawline, err := tp.ReadLine()
		if err != nil {
			break
		}

		switch {
		case regs["ping"].MatchString(rawline):
			wr.PrintfLine("PONG :" + rawline[4:])
			// " #" needed to skip messages from server when connecting.
		case strings.Contains(rawline, "PRIVMSG #"):
			ch, n, line := parseLine(rawline)
			os.Stdout.Write([]byte(ch + " : " + n + " | " + line + "\n"))
			//fmt.Printf("%s : %s | %s\n", ch, n, line)
			switch {
			case strings.Contains(line, "http"):
				r := parseURL_getTitle(line)
				if len(r) > 0 {
					for _, title := range r {
						grind.BotPrint(conn, ch, title)
					}
				}
			case regs["cmd"].MatchString(line):
				cmd := regs["cmd"].FindStringSubmatch(line)[1]
				grind.ProcessCmds(conn, ch, n, cmd)
			}
		default:
			os.Stdout.Write([]byte(rawline + "\n"))
			//fmt.Printf("%s\n", rawline)
		}
	}
}

func main() {
	conn := connect()
	read(conn)
}
