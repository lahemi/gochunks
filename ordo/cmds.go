package main

import (
	"io/ioutil"
	"os"
	"strings"
)

// Non-relative to current position, absolute.
// Use currentPos for relative movement.
func jumpChar(e *ENV) {
	a, err := e.Nums.PopE()
	if err != nil {
		return
	}
	i := a.(int)
	if i >= 0 && i < len(e.Text) {
		e.Pos = i
	}
}

func curLoadChar(e *ENV) {
	e.Strs.Push(string(e.Text[e.Pos]))
}

func searchCharF(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	for i := e.Pos + 1; i < len(e.Text); i++ {
		c := e.Text[i]
		if string(c) == a.(string) {
			e.Nums.Push(i - e.Pos)
			return
		}
	}
	// If it ever gets here, it means the char were not found.
	e.Branch = true
}

func searchCharB(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	for i := e.Pos - 1; i >= 0; i-- {
		c := e.Text[i]
		if string(c) == a.(string) {
			e.Nums.Push(e.Pos - i)
			return
		}
	}
	e.Branch = true
}

func deleteChar(e *ENV) {
	if e.Pos == 0 {
		e.Text = e.Text[e.Pos+1:]
		return
	}
	prev := e.Text[:e.Pos]
	next := e.Text[e.Pos+1:]
	prev = append(prev, next...)
	e.Text = prev
}

func insert(e *ENV) {
	prev, next := string(e.Text[:e.Pos]), string(e.Text[e.Pos:])
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	ts := prev
	ts += a.(string)
	ts += next
	e.Text = []rune(ts)
}

func printChar(e *ENV) {
	if len(e.Text) >= 0 && e.Pos >= 0 && e.Pos < len(e.Text) {
		stdout(string(e.Text[e.Pos]))
	}
}

func repeatCmd(e *ENV) {
	rn, err := e.Nums.PopE()
	if err != nil {
		return
	}
	cmd, err := e.Strs.PopE()
	if err != nil {
		return
	}
	scmd := cmd.(string)

	if c, ok := COMMANDS[scmd]; ok {
		for count := 0; count < rn.(int); count++ {
			c(e)
		}
		return
	}

	if m, ok := MACROS[scmd]; ok {
		for count := 0; count < rn.(int); count++ {
			eval(cmdList([]rune(m)), e)
		}
	}
}

func writeFile(e *ENV) {
	if err := ioutil.WriteFile(e.FName, []byte(string(e.Text)), 0666); err != nil {
		stderr("File write error.\n")
	}
}

func changeFile(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	e.FName = a.(string)
	e.Text = []rune("") // Should we preserve the text from previous file?
}

func quit(e *ENV) {
	os.Exit(0)
}

func eof(e *ENV) {
	e.Nums.Push(len(e.Text) - 1)
}

func currentPos(e *ENV) {
	e.Nums.Push(e.Pos)
}

func putChar(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	s := a.(string)
	e.Strs.Push(s)
	insert(e)
	e.Strs.Push(s)
}

func upperChar(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	e.Strs.Push(strings.ToUpper(a.(string)))
}

func lowerChar(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	e.Strs.Push(strings.ToLower(a.(string)))
}

func testChar(e *ENV) {
	a, err := e.Strs.PopE()
	if err != nil {
		return
	}
	curLoadChar(e)
	c := e.Strs.Pop()
	if c.(string) == a.(string) {
		// and then somehow use the ret of test.
	}
}
