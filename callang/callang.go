package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Operators map[string]string

// Got to come up with a cleaner scheme.
// Maybe a sublists, so that operators
// can be grouped according to type.
var (
	ops = Operators{
		"TERM": ";", "ASSIGN": ":",

		"PLUS": "+", "MINUS": "-",
		"TIMES": "×", "DIV": "÷",
		"MOD": "%",

		"ABS": "|", "NEG": "_",
		"INDEX": "ı",

		"CMT": "Ð", "VAR": "'",

		"REDUCE": "/",
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
	case o["ABS"], o["INDEX"]:
		return true
	default:
		return false
	}
}
func (o Operators) IsDyadic(oper string) bool {
	switch oper {
	case o["REDUCE"]:
		return true
	default:
		return false
	}
}
func (o Operators) IsMulRet(oper string) bool {
	switch oper {
	case o["INDEX"]:
		return true
	default:
		return false
	}
}
func (o Operators) IsOp(oper string) bool {
	switch {
	case o.IsArg1(oper) || o.IsArg2(oper) || o.IsDyadic(oper) || oper == o["VAR"]:
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
func (o Operators) Nop(oper string) {}

func (o Operators) RunMulRet(oper string, a1 float64) (ret Fstack) {
	switch oper {
	case o["INDEX"]:
		var i float64 = 1
		for ; i <= a1; i++ {
			ret = append(ret, i)
		}
	default:
		ret = Fstack{}
	}
	return
}

func (o Operators) RunDyadic(oper, argop string, args Fstack) (ret Fstack) {
	switch oper {
	case o["REDUCE"]:
		var (
			total float64 = args[0]
		)
		for c := 1; c < len(args); c++ {
			total = o.RunOp2(argop, total, args[c])
		}
		ret = Fstack{total}

	default:
		ret = Fstack{}
	}
	return
}

type Env map[string]Fstack

func execute(text string) string {
	const (
		RD       = 0
		INCMT    = 1
		INCOMP   = 2
		INVAR    = 4
		INDYADIC = 8
	)
	var (
		env = make(Env)
		spl = strings.Split(text, "") // To get UTF-8 chars

		parseloop func(int, int)

		buf      string
		dyadicop string
	)
	spl = append(spl, " ") // A "terminating" whitespace
	env["_G"] = Fstack{}

	parseloop = func(cp, state int) {
		if cp >= len(spl) {
			return
		}

		c := spl[cp]

		switch state {
		case INCMT:
			if c == "\n" {
				state = RD
			}

		case INCOMP:
			if c == ops["TERM"] || isWhite(c) {
				env[buf] = env["_G"]
				env["_G"] = Fstack{}
				buf = ""
				state = RD
			} else {
				buf += c
			}

		case INVAR:
			if isWhite(c) || ops.IsOp(c) || c == ops["TERM"] {
				if v, ok := env[buf]; ok {
					t := env["_G"]
					t = append(t, v...)
					env["_G"] = t
				}
				buf = ""
				state = RD
				if ops.IsOp(c) {
					cp--
				}
			} else {
				buf += c
			}

		case INDYADIC:
			if ops.IsOp(c) {
				t := env["_G"]
				r := ops.RunDyadic(dyadicop, c, t)
				t = append(t, r...)
				env["_G"] = t
			}
			dyadicop = ""
			state = RD

		case RD:
			switch {
			case isWhite(c) && buf == "":

			case c == ops["CMT"]:
				state = INCMT

			case (isWhite(c) || ops.IsOp(c) || c == ops["TERM"]) && buf != "":
				if strings.HasPrefix(buf, ops["NEG"]) {
					buf = strings.Replace(buf, "_", "-", -1)
				}
				if isNum(buf) {
					n, e := convNum(buf)
					if e == nil {
						t := env["_G"]
						t.Push(n)
						env["_G"] = t
					}
				}
				buf = ""
				if ops.IsOp(c) {
					cp--
				}

			case c == ops["VAR"]:
				state = INVAR

			case ops.IsDyadic(c):
				dyadicop = c
				state = INDYADIC

			case ops.IsArg1(c):
				t := env["_G"]
				if len(t) < 1 {
					return //"insufficient stack"
				}
				a1 := t.Pop()

				if ops.IsMulRet(c) {
					r := ops.RunMulRet(c, a1)
					t = append(t, r...)
				} else {
					r := ops.RunOp1(c, a1)
					if r == -0xffffffff {
						return //"error with operators"
					}
					t.Push(r)
				}
				env["_G"] = t
				buf = ""

			case ops.IsArg2(c):
				t := env["_G"]
				if len(t) < 2 {
					return //"insufficient stack"
				}
				a2, a1 := t.Pop(), t.Pop()

				r := ops.RunOp2(c, a1, a2)

				if r == -0xffffffff {
					return //"error with operators"
				}
				t.Push(r)
				env["_G"] = t
				buf = ""

			case c == ops["ASSIGN"]:
				buf = "" // 1 2 3:var; handle this ?
				state = INCOMP

			default:
				buf += c
			}
		}

		cp++
		parseloop(cp, state)
	}
	parseloop(0, RD)

	t := env["_G"]
	if len(t) < 1 {
		return "insufficient stack"
	}
	return strconv.FormatFloat(t.Pop(), 'g', 4, 64)
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
	fmt.Print(execute(input))
}
