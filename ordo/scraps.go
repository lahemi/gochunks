package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// In case something needs to be added to these.
func stdout(str string) {
	fmt.Fprint(os.Stdout, str)
}
func stderr(str string) {
	fmt.Fprint(os.Stderr, str)
}

const (
	TERMINATE = ''
)

func readInputFile(name string) []rune {
	cnt, err := ioutil.ReadFile(name)
	if err != nil {
		stderr("File reading error.")
		return []rune{}
	}
	return []rune(string(cnt))
}

func readCommands(file *os.File) (buf []byte) {
	for {
		charbuf := make([]byte, 1)
		_, err := file.Read(charbuf)
		if err != nil {
			break
		}
		if charbuf[0] == TERMINATE {
			break
		}
		buf = append(buf, charbuf[0])
	}
	return
}

// Is there a better way to convert a string
// literal to its corresponding escape ?
func lit2Esc(str string) string {
	var ESCS = map[string]string{
		`\a`: "\a", `\b`: "\b", `\f`: "\f",
		`\n`: "\n", `\r`: "\r", `\t`: "\t",
		`\v`: "\v",
	}
	for l, r := range ESCS {
		str = strings.Replace(str, l, r, -1)
	}
	return str
}

func cmdList(cmds []byte) (ret []string) {
	const (
		RD = iota
		STR
	)
	var (
		buf   string
		state = RD
	)
	for i := 0; i < len(cmds); i++ {
		c := string(cmds[i])
		switch state {
		case RD:
			switch {
			case buf == "" && isWhite(c):

			case buf != "" && isWhite(c):
				ret = append(ret, buf)
				buf = ""

			case c == "`":
				buf += c
				state = STR

			default:
				buf += c
			}

		case STR:
			if c == "`" {
				state = RD
				buf = "`" + lit2Esc(string(buf[1:]))
				ret = append(ret, buf)
				buf = ""
			} else {
				buf += c
			}
		}
	}
	if buf != "" {
		ret = append(ret, buf)
	}
	return
}

// Limited selection of whites, as it acts as a separator, too.
func isWhite(c string) bool {
	switch c {
	case " ", "\t", "\n":
		return true
	default:
		return false
	}
}

func isInt(n string) bool {
	for i := 0; i < len(n); i++ {
		if i == 0 && n[i] == '-' {
			continue
		}
		if n[i] < 48 || n[i] > 57 {
			return false
		}
	}
	return true
}

func compileMacros(text []rune) MACROS {
	const (
		FUNSTART = '('
		FUNEND   = ')'
		FUNNAME  = '→'

		RD = iota
		COMP
	)
	var (
		state = RD
		buf   []rune
		bbuf  []rune
		ret   = MACROS{}
	)
	for _, r := range text {
		switch state {
		case RD:
			if r == FUNSTART {
				state = COMP
			}

		case COMP:
			switch r {
			case FUNEND:
				ret[strings.Trim(string(buf), " ")] = string(bbuf)
				buf = []rune{}
				bbuf = []rune{}
				state = RD

			case FUNNAME:
				bbuf = buf
				buf = []rune{}

			default:
				buf = append(buf, r)
			}
		}
	}
	return ret
}

func loadMacros(file string) MACROS {
	return compileMacros(readInputFile(file))
}
