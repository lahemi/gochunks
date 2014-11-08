package main

import (
	"os"
	"strconv"
	"strings"
)

type ENV struct {
	FName   string
	Text    []rune
	Pos     int
	Numargs Stack
	Strargs Stack
}

type MACROSET map[string]string
type COMMANDSET map[string]func(*ENV)

const (
	STRINGMARKER = "`"
	MACROMARKER  = "'"
)

// This is required so that we may construct the actual
// used commands from an external config file.
// So horribly dynamic!
var COMMANDTABLE = COMMANDSET{
	"moveChar":       moveChar,
	"jumpChar":       jumpChar,
	"searchCharF":    searchCharF,
	"searchCharB":    searchCharB,
	"deleteChar":     deleteChar,
	"insertChar":     insertChar,
	"printChar":      printChar,
	"quit":           quit,
	"writeFile":      writeFile,
	"changeFile":     changeFile,
	"addition":       addition,
	"subtraction":    subtraction,
	"multiplication": multiplication,
	"division":       division,
	"repeatCmd":      repeatCmd,
	"eof":            eof,
}

var (
	G_env    = ENV{}
	MACROS   = MACROSET{}
	COMMANDS = COMMANDSET{}

	inputFile string
)

func eval(cmds []string, env *ENV) {
	for _, c := range cmds {
		if cmd, ok := COMMANDS[c]; ok {
			cmd(env)
		} else if isInt(c) {
			// Yes, but we test for its validity already with isInt.
			i, _ := strconv.Atoi(c)
			env.Numargs.Push(i)
		} else if strings.HasPrefix(c, STRINGMARKER) {
			cc := string(c[len(STRINGMARKER):])
			env.Strargs.Push(cc)
		} else if strings.HasPrefix(c, MACROMARKER) {
			if m, ok := MACROS[string(c[len(MACROMARKER):])]; ok {
				eval(cmdList([]rune(m)), env)
			}
		}
	}
}

func init() {
	if len(os.Args) != 2 {
		stderr("Need an input file.\n")
		os.Exit(1)
	}
	inputFile = os.Args[1]

	// Need proper places for these files.
	COMMANDS = loadCommands("commands.rc")
	MACROS = loadMacros("macros.rc")
}

func main() {
	if inputFile == "" {
		return
	}
	text := readInputFile(inputFile)
	G_env.FName = inputFile
	G_env.Text = text
	G_env.Pos = 0
	for {
		cmds := cmdList(readCommands(os.Stdin))
		eval(cmds, &G_env)
	}
}
