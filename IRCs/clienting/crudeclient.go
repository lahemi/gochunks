package main

import (
	"bufio"
	"net"
	"os"
)

var writechan = make(chan string, 512)

func doeverything(server, port, nick string) {
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
				os.Stderr.Write([]byte("read error\n"))
				break
			}
			if str[:4] == "PING" {
				writechan <- "PONG" + str[4:len(str)-2]
			} else {
				os.Stdout.Write([]byte(str))
			}
		}
	}()
	go func() {
		for {
			str := <-writechan
			if _, err := w.WriteString(str + "\r\n"); err != nil {
				os.Stderr.Write([]byte("write error\n"))
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
			os.Stderr.Write([]byte("error input\n"))
			break
		}
		writechan <- input[:len(input)-1]
	}
}

func main() {
	doeverything("irc.freenode.net", "6667", "gecannels")
}
