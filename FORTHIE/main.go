package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type IDTYPE int
type ELEM struct {
	Type IDTYPE
	Cont string
}

const (
	INUM IDTYPE = iota
	FNUM
	STR
	INUMTBL
	FNUMTBL
	STRTBL
	IDENT
	BRANCH
	TESTBRANCH
	BRANCHEND
	ASSIGN
	ASSIGNEND
	UNKNOWN
)

var TYPES = map[IDTYPE]string{
	INUM: "INUM", FNUM: "FNUM", STR: "STR",
	INUMTBL: "INUMTBL", FNUMTBL: "FNUMTBL",
	STRTBL: "STRTBL", IDENT: "IDENT",
	BRANCH: "BRANCH", TESTBRANCH: "TESTBRANCH",
	BRANCHEND: "BRANCHEND", ASSIGN: "ASSIGN",
	ASSIGNEND: "ASSIGNEND", UNKNOWN: "UNKNOWN",
}

func whichType(el string) (typ IDTYPE) {
	switch {
	case isInt(el):
		typ = INUM
	case isFloat(el):
		typ = FNUM
	case isString(el): // doesn't work proper
		typ = STR
	case el == ":":
		typ = ASSIGN
	case el == ";":
		typ = ASSIGNEND
	case el == "IF":
		typ = TESTBRANCH
	case el == "ELSE":
		typ = BRANCH
	case el == "END":
		typ = BRANCHEND
	default:
		typ = IDENT // It might or might not be, that'll be checked later.
	}
	return
}

func findBranchEnd(startpoint int, elems []ELEM) int {
	var brcount = -1
	for i := startpoint; i < len(elems); i++ {
		if elems[i].Type == TESTBRANCH {
			brcount++
		}
		if elems[i].Type == BRANCHEND || elems[i].Type == BRANCH {
			if brcount <= 0 {
				return i
			}
			brcount--
		}
	}
	return -1 // Not found !
}

// handle STRINGs properly
func untypeds(text string) (ret []string) {
	var (
		buf string
		spl = strings.Split(text, "")
	)
	spl = append(spl, " ") // A "terminating" whitespace.
	for tp := 0; tp < len(spl); tp++ {
		switch {
		case isWhite(spl[tp]) && buf == "":
		case isWhite(spl[tp]) && buf != "":
			ret = append(ret, buf)
			buf = ""
		default:
			buf += spl[tp]
		}
	}
	return
}
func createASTree(untypeds []string) (elems []ELEM) {
	for _, e := range untypeds {
		t := whichType(e)
		el := ELEM{t, e}
		elems = append(elems, el)
	}
	return
}
func printASTree(elems []ELEM) {
	for i, e := range elems {
		if v, ok := TYPES[e.Type]; ok {
			fmt.Printf("%d | %s %s\n", i, v, e.Cont)
		}
	}
}

type MODES int

const (
	INT MODES = iota
	FLOAT
	STRING
)

// add underflow checks!
type ENV struct {
	INTS    Istack
	FLOATS  Fstack
	STRINGS Sstack
	MODE    MODES
}

func (e *ENV) SetMode(m MODES) {
	e.MODE = m
}

var (
	GENV   = ENV{}
	IDENTS = map[string][]ELEM{}
)

func compile(code []ELEM) {
	var (
		loop   func(int)
		cs     int
		idname string
	)
	loop = func(cp int) {
		if cp >= len(code) {
			return
		}
		switch code[cp].Type {
		case ASSIGN:
			cp++
			idname = code[cp].Cont
			cp++
			cs = cp
		case ASSIGNEND:
			IDENTS[idname] = code[cs:cp]
			idname = ""
		}

		loop((cp + 1))
	}
	loop(0)

	return
}

func removeComps(code []ELEM) (ret []ELEM) {
	var inComp = false
	for _, el := range code {
		switch {
		case el.Type == ASSIGNEND:
			inComp = false
		case el.Type == ASSIGN:
			inComp = true
		case !inComp:
			ret = append(ret, el)
		}
	}
	return
}

func eval(code []ELEM, env ENV) {
	var evalloop func(int, []ELEM)
	evalloop = func(cp int, code []ELEM) {
		if cp >= len(code) {
			return
		}

		c := code[cp]

		switch c.Type {
		case INUM:
			i, e := strconv.Atoi(c.Cont)
			if e != nil {
				fmt.Fprintf(os.Stderr, "row %d err at int conv", cp)
				return
			}
			env.INTS.Push(i)
		case FNUM:
			f, e := strconv.ParseFloat(c.Cont, 64)
			if e != nil {
				fmt.Fprintf(os.Stderr, "row %d err at float conv", cp)
				return
			}
			env.FLOATS.Push(f)
		case STR:
			s := c.Cont[1:]
			s = s[:len(s)-1]
			env.STRINGS.Push(s)
		case IDENT:
			if v, ok := IDENTS[c.Cont]; ok {
				evalloop(0, v)
			} else if v, ok := BUILTINS[c.Cont]; ok {
				v(&env)
			} else {
				return
			}
		case TESTBRANCH:
			testarg := env.INTS.Pop()
			branchend := findBranchEnd(cp, code)
			if testarg == 0 {
				if branchend != -1 {
					cp = branchend
				}
			}
		case BRANCH:
			branchend := findBranchEnd((cp + 1), code)
			if branchend != -1 {
				cp = branchend
			}
		case BRANCHEND:
			// nothing to do here, skip.
		default:
			return
		}

		evalloop((cp + 1), code)
	}

	evalloop(0, code)
}

func readText(file string) string {
	cnt, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return string(cnt)
}

func main() {
	cptest := createASTree(untypeds(readText(os.Args[1])))
	compile(cptest)
	for id, c := range IDENTS {
		fmt.Println(id)
		printASTree(c)
	}
	els := removeComps(cptest)
	fmt.Println()
	printASTree(els)
	eval(els, GENV)
}
