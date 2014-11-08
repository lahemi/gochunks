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
		"j": jumpChar,
		"s": searchCharF,
		"r": searchCharB,
		"d": deleteChar,
		"i": insertChar,
		"p": printChar,
		"q": func(e *ENV) { os.Exit(0) },
		"w": writeFile,
		"c": changeFile, // empties e.Text !
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
				eval(cmdList([]rune(m)), env)
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
	G_env.FName = inputFile
	G_env.Text = text
	G_env.Pos = 0
	for {
		cmds := cmdList(readCommands(os.Stdin))
		eval(cmds, &G_env)
	}
}

/*
char:
    move char
        relative        done
        absolute        done

    search forward      done
    search back         done

    delete              done

    insert              done


search a word ?


macros                          done, for pre-def'd

interactive mode macros defs    done


repeat command                  done

ie.
    (1 m p → pm)
    3 `pm` rep
        → 1 m p 1 m p 1 m p

        basically just string/command
        duplication.


^G to mean just execute,    done
not  exec & quit

quit command                done


if given non-existing       done
file, create it

save file                   done

*/
