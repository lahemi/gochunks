package main

import (
	"os"
	"strconv"
	"strings"
)

type ENV struct {
	Text    []rune
	Pos     int
	Numargs Stack
	Strargs Stack
}

type MACROS map[string]string

const (
	STRINGMARKER = "`"
	MACROMARKER  = "'"
)

var (
	MACROTABLE = MACROS{}
	G_env      = ENV{}

	COMMANDS = map[string]func(*ENV){
		"m": moveChar,
		"s": searchCharF,
		"r": searchCharB,
		"d": deleteChar,
		"i": insertChar,
		"p": printChar,
	}

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
			if m, ok := MACROTABLE[string(c[len(MACROMARKER):])]; ok {
				eval(cmdList([]byte(m)), env)
			}
		}
	}
}

func init() {
	if len(os.Args) != 2 {
		stderr("Need a input file.\n")
		os.Exit(1)
	}
	inputFile = os.Args[1]

	// This is separately here, because if we
	// place it in the same spot as the rest of
	// the commands, we get an error about circular
	// reference between COMMANDS and repeatCmd.
	COMMANDS["rep"] = repeatCmd

	MACROTABLE = loadMacros("macros.rc")
}

func main() {
	if inputFile == "" {
		return
	}
	text := readInputFile(inputFile)
	cmds := cmdList(readCommands(os.Stdin))
	G_env.Text = text
	G_env.Pos = 0
	eval(cmds, &G_env)
}

/*
char:
    move char           done

    search forward      done
    search back         done

    delete              done

    insert              done


search a word ?


macros                  done, for pre-def'd

(1 m p → pm)
(`\n` s m → $)

Add macro defs and a macro table,
so that shorhands like that can be used.


Also, interactive def of macros, not only
in a separate file loaded at start-up.

repeat command          done, for COMMANDS

ie.
    (1 m p → pm)
    3 `pm` rep
        → 1 m p 1 m p 1 m p

        basically just string/command
        duplication.
*/
