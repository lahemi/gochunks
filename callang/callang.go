package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// different operator types
const (
	NOP     OPType = 0
	SPECIAL OPType = iota
	ONEARG
	TWOARG
	DYADIC
	MULRET
	VAR
)

type OPType int
type OP struct {
	Name string
	Type OPType
}
type Operators map[string]OP

var (
	ops = Operators{
		"TERM":   OP{";", SPECIAL},
		"ASSIGN": OP{":", SPECIAL},
		"PLUS":   OP{"+", TWOARG},
		"MINUS":  OP{"-", TWOARG},
		"TIMES":  OP{"×", TWOARG},
		"DIV":    OP{"÷", TWOARG},
		"MOD":    OP{"%", TWOARG},
		"RSFT":   OP{"«", TWOARG},
		"LSFT":   OP{"»", TWOARG},

		"ABS":   OP{"|", ONEARG},
		"NEG":   OP{"_", SPECIAL},
		"INDEX": OP{"ı", MULRET},

		"CMT": OP{"Ð", SPECIAL},
		"VAR": OP{"'", VAR},

		"REDUCE": OP{"/", DYADIC},
	}
)

func (o Operators) WhichType(oper string) OPType {
	for _, v := range ops {
		if v.Name == oper {
			return v.Type
		}
	}
	return NOP
}
func (o Operators) IsOp(oper string) bool {
	t := o.WhichType(oper)
	switch t {
	case ONEARG, TWOARG, DYADIC, MULRET, VAR:
		return true
	default:
		return false
	}
}

// Feels foolish, got to check if there'd be a better way to do this.
func (o Operators) RunOp2(oper string, a1, a2 float64) (ret float64) {
	switch oper {
	case o["PLUS"].Name:
		ret = a1 + a2
	case o["MINUS"].Name:
		ret = a1 - a2
	case o["TIMES"].Name:
		ret = a1 * a2
	case o["DIV"].Name:
		ret = a1 / a2
	case o["MOD"].Name:
		ret = math.Mod(a1, a2)
	default:
		ret = -0xffffffff // A bad pseudo-error
	}
	return
}
func (o Operators) RunOp1(oper string, a1 float64) (ret float64) {
	switch oper {
	case o["ABS"].Name:
		ret = math.Abs(a1)
	default:
		ret = -0xffffffff
	}
	return
}

func (o Operators) RunMulRet(oper string, a1 float64) (ret Fstack) {
	switch oper {
	case o["INDEX"].Name:
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
	case o["REDUCE"].Name:
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
		RD = iota
		INCMT
		INCOMP
		INVAR
		INDYADIC
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
			if c == ops["TERM"].Name || isWhite(c) {
				env[buf] = env["_G"]
				env["_G"] = Fstack{}
				buf = ""
				state = RD
			} else {
				buf += c
			}

		case INVAR:
			if isWhite(c) || ops.IsOp(c) || c == ops["TERM"].Name {
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
				t = r
				env["_G"] = t
			}
			dyadicop = ""
			state = RD

		case RD:
			switch {
			case isWhite(c) && buf == "":

			case c == ops["CMT"].Name:
				state = INCMT

			case (isWhite(c) || ops.IsOp(c) || c == ops["TERM"].Name) && buf != "":
				if strings.HasPrefix(buf, ops["NEG"].Name) {
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

			case ops.WhichType(c) == VAR:
				state = INVAR

			case ops.WhichType(c) == DYADIC:
				dyadicop = c
				state = INDYADIC

			case ops.WhichType(c) == ONEARG || ops.WhichType(c) == MULRET:
				t := env["_G"]
				if len(t) < 1 {
					return //"insufficient stack"
				}
				a1 := t.Pop()

				if ops.WhichType(c) == MULRET {
					r := ops.RunMulRet(c, a1)
					t = r
				} else {
					r := ops.RunOp1(c, a1)
					if r == -0xffffffff {
						return //"error with operators"
					}
					t.Push(r)
				}
				env["_G"] = t
				buf = ""

			case ops.WhichType(c) == TWOARG:
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

			case c == ops["ASSIGN"].Name:
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
