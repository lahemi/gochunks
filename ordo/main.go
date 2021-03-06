package main

import (
	"github.com/lahemi/stack"
	"os"
	"strconv"
	"strings"
)

type ENV struct {
	FName  string
	Text   []rune
	Pos    int
	Nums   stack.Stack
	Strs   stack.Stack
	GenSt  stack.Stack // for tucking away things...
	Branch bool
}

type MACROSET map[string]string
type COMMANDSET map[string]func(*ENV)

const (
	STRINGMARKER = "`"
	BRANCHSTART  = "["
	BRANCHSEP    = "|"
	BRANCHEND    = "]"
)

var COMMANDS = COMMANDSET{
	"jumpChar":       jumpChar,
	"searchCharF":    searchCharF,
	"searchCharB":    searchCharB,
	"deleteChar":     deleteChar,
	"insert":         insert,
	"printChar":      printChar,
	"quit":           quit,
	"writeFile":      writeFile,
	"changeFile":     changeFile,
	"addition":       addition,
	"subtraction":    subtraction,
	"multiplication": multiplication,
	"division":       division,
	"eof":            eof,
	"currentPos":     currentPos,
	"curLoadChar":    curLoadChar,
	"putChar":        putChar,
	"upperChar":      upperChar,
	"lowerChar":      lowerChar,

	"StrPop": func(e *ENV) {
		if r, err := e.Strs.PopE(); err == nil {
			e.GenSt.Push(r)
		}
	},
	"StrDup":  func(e *ENV) { e.Strs.Dup() },
	"StrDrop": func(e *ENV) { e.Strs.Drop() },
	"StrSwap": func(e *ENV) { e.Strs.Swap() },
	"StrOver": func(e *ENV) { e.Strs.Over() },
	"StrRot":  func(e *ENV) { e.Strs.Rot() },
	"NumPop": func(e *ENV) {
		if r, err := e.Nums.PopE(); err == nil {
			e.GenSt.Push(r)
		}
	},
	"NumDup":  func(e *ENV) { e.Nums.Dup() },
	"NumDrop": func(e *ENV) { e.Nums.Drop() },
	"NumSwap": func(e *ENV) { e.Nums.Swap() },
	"NumOver": func(e *ENV) { e.Nums.Over() },
	"NumRot":  func(e *ENV) { e.Nums.Rot() },

	// ...
	"GenPop": func(e *ENV) {
		if s, err := e.Strs.PopE(); err == nil {
			switch s.(string) {
			case "s":
				if r, err := e.GenSt.PopE(); err == nil {
					e.Strs.Push(r.(string))
				}
			case "n":
				if r, err := e.GenSt.PopE(); err == nil {
					e.Nums.Push(r.(int))
				}
			}
		}
	},
	"GenDup":  func(e *ENV) { e.GenSt.Dup() },
	"GenDrop": func(e *ENV) { e.GenSt.Drop() },
	"GenSwap": func(e *ENV) { e.GenSt.Swap() },
	"GenOver": func(e *ENV) { e.GenSt.Over() },
	"GenRot":  func(e *ENV) { e.GenSt.Rot() },
}

var (
	G_env  = ENV{}
	MACROS = MACROSET{}

	inputFile string
)

func eval(cmds []string, env *ENV) {
	for i := 0; i < len(cmds); i++ {
		c := cmds[i]

		if cmd, ok := COMMANDS[c]; ok {
			cmd(env)

		} else if isInt(c) {
			// Yes, but we test for its validity already with isInt.
			i, _ := strconv.Atoi(c)
			env.Nums.Push(i)

		} else if strings.HasPrefix(c, STRINGMARKER) {
			cc := string(c[len(STRINGMARKER):])
			env.Strs.Push(cc)

			// Ugh, so nasty... But it works, at least for simple cases.
		} else if c == BRANCHSTART {
			if env.Branch {
				for i++; i < len(cmds); i++ {
					if cmds[i] == BRANCHSEP {
						break
					}
				}
				env.Branch = false
			}

		} else if c == BRANCHSEP {
			for i++; i < len(cmds); i++ {
				if cmds[i] == BRANCHEND {
					break
				}
			}

		} else if m, ok := MACROS[c]; ok {
			eval(cmdList([]rune(m)), env)
		}
	}
}

func init() {
	if len(os.Args) != 2 {
		stderr("Need an input file.\n")
		os.Exit(1)
	}
	inputFile = os.Args[1]

	// This needs to be separate from the rest,
	// if we define it at the same time we do
	// the others, we get a warning about
	// a circular reference.
	COMMANDS["repeatCmd"] = repeatCmd

	// Need proper places for these files.
	MACROS = loadMacros("macros.rc")
}

func main() {
	if inputFile == "" {
		return
	}
	text := readInputFile(inputFile)
	G_env.FName = inputFile
	G_env.Text = text
	for {
		cmds := cmdList(readCommands(os.Stdin))
		eval(cmds, &G_env)
	}
}
