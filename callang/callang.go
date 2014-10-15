package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

/*
x:4;y:5; f:{_1_2+}; f(x,y);
x y +;
*/

func isWhite(c string) bool {
	switch c {
	case " ", "\t", "\n":
		return true
	default:
		return false
	}
}

type Fstack []float64

func (s *Fstack) Push(f float64) {
	(*s) = append((*s), f)
}
func (s *Fstack) Pop() float64 {
	if len((*s)) <= 0 {
		return -0xffffffff
	}
	last := (*s)[len((*s))-1]
	(*s) = (*s)[:len((*s))-1]
	return last
}

type Operators struct {
	TERM, ASSIGN string

	PLUS, MINUS string
	TIMES, DIV  string
	MOD         string

	CMT string
	VAR string
}

func (o *Operators) IsArith(oper string) bool {
	switch oper {
	case o.PLUS, o.MINUS, o.TIMES, o.DIV, o.MOD:
		return true
	default:
		return false
	}
}
func (o *Operators) IsOp(oper string) bool {
	switch {
	case o.IsArith(oper) || oper == o.VAR:
		return true
	default:
		return false
	}
}

// Feels foolish, got to check if there'd be a better way to do this.
func (o *Operators) RunOp(oper string, a1, a2 float64) (ret float64) {
	switch oper {
	case o.PLUS:
		ret = a1 + a2
	case o.MINUS:
		ret = a1 - a2
	case o.TIMES:
		ret = a1 * a2
	case o.DIV:
		ret = a1 / a2
	case o.MOD:
		ret = math.Mod(a1, a2)
	default:
		ret = -0xffffffff // A bad pseudo-error
	}
	return
}

type Env map[string]float64

var (
	ops = Operators{
		TERM:   ";",
		ASSIGN: ":",

		PLUS:  "+",
		MINUS: "-",
		TIMES: "*",
		DIV:   "/",
		MOD:   "%",

		CMT: "Ð",
		VAR: "'",
	}
)

const (
	NL = "\n"

	RD     = 0
	INCMT  = 1
	INCOMP = 2
)

func execute(text string, env Env) string {
	var (
		spl    = strings.Split(text, "") // To get UTF-8 chars
		fstack = Fstack{}
		buf    string
	)

	for i := 0; i < len(spl); i++ {
		c := spl[i]
		switch {
		case c == ops.TERM || isWhite(c):
			if num, err := strconv.ParseFloat(buf, 64); err == nil {
				fstack.Push(num)
			}
			buf = ""
		case c == ops.VAR:
			buf = ""
			for i++; i < len(spl); i++ {
				c := spl[i]
				if ops.IsOp(c) || isWhite(c) {
					if f, ok := env[buf]; ok {
						fstack.Push(f)
					}
					if ops.IsOp(c) {
						i--
					}
					buf = ""
					break
				}
				buf += c
			}
		case ops.IsOp(c):
			if len(fstack) < 2 {
				return ""
			}
			a2, a1 := fstack.Pop(), fstack.Pop()

			r := ops.RunOp(c, a1, a2)

			if r == -0xffffffff {
				return "error with operators"
			}
			fstack.Push(r)

			buf = ""

		default:
			buf += c
		}
	}

	if len(fstack) < 1 {
		return "Insufficient stack"
	}
	return strconv.FormatFloat(fstack.Pop(), 'g', 4, 64)
}

func constructEnv(text string) (string, Env) {
	var (
		env = make(Env)
		spl = strings.Split(text, "") // To get UTF-8 chars

		parseloop func(int, int)

		buf  string
		vbuf string
	)

	parseloop = func(cp, state int) {
		if cp >= len(spl) {
			return
		}

		c := spl[cp]

		switch state {
		case INCMT:
			if c == NL {
				state = RD
			}

		case INCOMP:
			if c == ops.TERM {
				if v, ok := env[buf]; ok {
					env[vbuf] = v
				} else if buf == vbuf {
					env[vbuf] = env[vbuf]
				} else if num, err := strconv.ParseFloat(buf, 64); err == nil {
					env[vbuf] = num
				}
				vbuf = ""
				buf = ""
				state = RD
			} else {
				buf += c
			}

		case RD:
			switch {
			case isWhite(c) && buf == "":

			case c == ops.CMT:
				state = INCMT

			case c == ops.ASSIGN && buf != "":
				vbuf = buf
				buf = ""
				state = INCOMP

			default:
				buf += c
			}
		}

		cp++
		parseloop(cp, state)
	}
	parseloop(0, RD)

	return text, env
}

var input string

func init() {
	flag.StringVar(&input, "s", "", "Input code as a cmdline arg.")
	flag.Parse()
}

func main() {
	if input == "" {
		return
	}
	fmt.Print(execute(constructEnv(input)))
}
