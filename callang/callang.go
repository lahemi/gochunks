package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Limited selection of whites, as it acts as a separator, too.
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
		return -0xffffffff // pseudo-error
	}
	last := (*s)[len((*s))-1]
	(*s) = (*s)[:len((*s))-1]
	return last
}

type Operators map[string]string

var (
	ops = Operators{
		"TERM": ";", "ASSIGN": ":",

		"PLUS": "+", "MINUS": "-",
		"TIMES": "×", "DIV": "÷",
		"MOD": "%",

		"ABS": "|", "NEG": "_",

		"CMT": "Ð", "VAR": "'",
	}
)

func (o Operators) IsArg2(oper string) bool {
	switch oper {
	case o["PLUS"], o["MINUS"], o["TIMES"], o["DIV"], o["MOD"]:
		return true
	default:
		return false
	}
}
func (o Operators) IsArg1(oper string) bool {
	switch oper {
	case o["ABS"]:
		return true
	default:
		return false
	}
}
func (o Operators) IsOp(oper string) bool {
	switch {
	case o.IsArg1(oper) || o.IsArg2(oper) || oper == o["VAR"]:
		return true
	default:
		return false
	}
}

// Feels foolish, got to check if there'd be a better way to do this.
func (o Operators) RunOp2(oper string, a1, a2 float64) (ret float64) {
	switch oper {
	case o["PLUS"]:
		ret = a1 + a2
	case o["MINUS"]:
		ret = a1 - a2
	case o["TIMES"]:
		ret = a1 * a2
	case o["DIV"]:
		ret = a1 / a2
	case o["MOD"]:
		ret = math.Mod(a1, a2)
	default:
		ret = -0xffffffff // A bad pseudo-error
	}
	return
}
func (o Operators) RunOp1(oper string, a1 float64) (ret float64) {
	switch oper {
	case o["ABS"]:
		ret = math.Abs(a1)
	default:
		ret = -0xffffffff
	}
	return
}

type Env map[string]float64

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
	spl = append(spl, " ") // Adds a "terminating" whitespace

	for i := 0; i < len(spl); i++ {
		c := spl[i]
		switch {
		case c == ops["TERM"] || isWhite(c) || c == ops["VAR"]:
			if strings.HasPrefix(buf, ops["NEG"]) {
				buf = strings.Replace(buf, "_", "-", -1)
			}
			if num, err := strconv.ParseFloat(buf, 64); err == nil {
				fstack.Push(num)
			}
			buf = ""

			// VAR both terminates previous exp, and references a variable
			if c == ops["VAR"] {
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
			}

			// duplication, will change
		case ops.IsArg1(c):
			if len(fstack) < 1 {
				return "insufficient stack"
			}
			a1 := fstack.Pop()

			r := ops.RunOp1(c, a1)

			if r == -0xffffffff {
				return "error with operators"
			}
			fstack.Push(r)

			buf = ""

		case ops.IsArg2(c):
			if len(fstack) < 2 {
				return "insufficient stack"
			}
			a2, a1 := fstack.Pop(), fstack.Pop()

			r := ops.RunOp2(c, a1, a2)

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
		return "insufficient stack"
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
			if c == ops["TERM"] {
				if v, ok := env[buf]; ok {
					env[vbuf] = v
				} else if buf == vbuf {
					env[vbuf] = env[vbuf]
				} else {
					if strings.HasPrefix(buf, ops["NEG"]) {
						buf = strings.Replace(buf, "_", "-", -1)
					}
					if num, err := strconv.ParseFloat(buf, 64); err == nil {
						env[vbuf] = num
					}
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

			case c == ops["CMT"]:
				state = INCMT

			case c == ops["ASSIGN"] && buf != "":
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
