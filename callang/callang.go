package main

import (
	"flag"
	"fmt"
	"math"
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

type Env map[string]Fstack
type FunEnv map[string]string

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
		// These two do not really work prorerly.
		"RSFT": OP{"»", TWOARG},
		"LSFT": OP{"«", TWOARG},

		"ABS":   OP{"|", ONEARG},
		"CEIL":  OP{"⌈", ONEARG},
		"FLOOR": OP{"⌊", ONEARG},
		"NEG":   OP{"_", SPECIAL},
		"INDEX": OP{"ı", MULRET},
		"DROP":  OP{"Ð", SPECIAL},

		"VAR": OP{"'", VAR},
		"CMT": OP{"Ħ", SPECIAL},

		"REDUCE": OP{"/", DYADIC},
		"APPLY":  OP{"º", DYADIC},

		"PRINT":      OP{",", SPECIAL},
		"PRINTSTACK": OP{"ß", SPECIAL},

		"FUNSTART": OP{"(", SPECIAL},
		"FUNEND":   OP{")", SPECIAL},
		"FUNNAME":  OP{"←", SPECIAL},

		"HIDE":   OP{"↑", SPECIAL},
		"UNHIDE": OP{"↓", SPECIAL},
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
		// These seem to give somewhat dubious results
	case o["RSFT"].Name:
		ret = float64(uint(a1) >> uint(a2))
	case o["LSFT"].Name:
		ret = float64(uint(a1) << uint(a2))
	default:
		ret = -0xffffffff // A bad pseudo-error
	}
	return
}
func (o Operators) RunOp1(oper string, a1 float64) (ret float64) {
	switch oper {
	case o["ABS"].Name:
		ret = math.Abs(a1)
	case o["CEIL"].Name:
		ret = math.Ceil(a1)
	case o["FLOOR"].Name:
		ret = math.Floor(a1)
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
		var total float64 = args[0]
		for c := 1; c < len(args); c++ {
			total = o.RunOp2(argop, total, args[c])
		}
		ret = Fstack{total}

	case o["APPLY"].Name:
		if len(args) < 2 {
			return args
		}
		switch o.WhichType(argop) {
		case TWOARG:
			applyarg := args.Pop()
			for i := 0; i < len(args); i++ {
				ret = append(ret, o.RunOp2(argop, applyarg, args[i]))
			}
		case ONEARG:
			for i := 0; i < len(args); i++ {
				ret = append(ret, o.RunOp1(argop, args[i]))
			}
		default:
			ret = Fstack{}
		}

	default:
		ret = Fstack{}
	}
	return
}

// So monolithic...
func execute(text string) {
	const (
		RD = iota
		INCMT
		INCOMP
		INVAR
		INDYADIC
		INFUN
	)
	var (
		env  = make(Env)
		fenv = make(FunEnv)
		spl  = strings.Split(text, "") // To get UTF-8 chars

		parseloop func(int, int)

		buf      string
		ebuf     string
		dyadicop string
	)
	spl = append(spl, " ") // A "terminating" whitespace
	env["_G"] = Fstack{}
	env["_H"] = Fstack{}

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

		case INFUN:
			switch {
			case c == ops["FUNEND"].Name:
				buf = strings.TrimSpace(buf)
				fenv[buf] = ebuf
				ebuf = ""
				buf = ""
				state = RD
			case c == ops["FUNNAME"].Name:
				ebuf = buf
				buf = ""
			default:
				buf += c
			}

		case INVAR:
			if isWhite(c) || ops.IsOp(c) || c == ops["TERM"].Name {
				if v, ok := env[buf]; ok {
					t := env["_G"]
					t = append(t, v...)
					env["_G"] = t
					if ops.IsOp(c) {
						cp--
					}
				} else if v, ok := fenv[buf]; ok {
					// A bit unclear. Substition model.
					vs := strings.Split(v, "")
					tt := spl[cp:]
					vs = append(vs, tt...)
					spl = vs
					cp = -1 // Because of parseloop((cp+1), state)
				}
				buf = ""
				state = RD
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
				buf = ""
				state = INVAR

			case ops.WhichType(c) == DYADIC:
				dyadicop = c
				state = INDYADIC

			case c == ops["FUNSTART"].Name:
				buf = ""
				state = INFUN

			case c == ops["HIDE"].Name:
				if len(env["_G"]) < 1 {
					fmt.Print("insufficient stack")
					return
				}
				t := env["_G"]
				h := env["_H"]
				h.Push(t.Pop())
				env["_G"] = t
				env["_H"] = h

			case c == ops["UNHIDE"].Name:
				if len(env["_H"]) < 1 {
					fmt.Print("insufficient stack")
					return
				}
				t := env["_G"]
				h := env["_H"]
				t.Push(h.Pop())
				env["_G"] = t
				env["_H"] = h

			case ops.WhichType(c) == ONEARG || ops.WhichType(c) == MULRET:
				t := env["_G"]
				if len(t) < 1 {
					fmt.Print("insufficient stack")
					return
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
					fmt.Print("insufficient stack")
					return
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

			case c == ops["PRINT"].Name:
				t := env["_G"]
				if len(t) < 1 {
					fmt.Print("insufficient stack")
					return
				}
				fmt.Print(t[len(t)-1])

			case c == ops["PRINTSTACK"].Name:
				fmt.Print(env["_G"])

			case c == ops["DROP"].Name:
				t := env["_G"]
				if len(t) < 1 {
					fmt.Print("insufficient stack")
					return
				}
				t.Pop()
				env["_G"] = t

			default:
				buf += c
			}
		}
		parseloop((cp + 1), state)
	}
	parseloop(0, RD)
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
	execute(input)
}
