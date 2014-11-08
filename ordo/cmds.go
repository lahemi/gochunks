package main

import (
	"io/ioutil"
)

// Relative to current position.
func moveChar(e *ENV) {
	i, err := e.Numargs.PopE()
	if err != nil {
		return
	}
	tp := e.Pos + i.(int)
	if tp < len(e.Text) && tp >= 0 {
		e.Pos = tp
	}
}

// Irrelative to current position, absolute.
func jumpChar(e *ENV) {
	a, err := e.Numargs.PopE()
	if err != nil {
		return
	}
	i := a.(int)
	if i >= 0 && i < len(e.Text) {
		e.Pos = i
	}
}

// This will be for yanking.
func curLoadChar(e *ENV) {
	e.Strargs.Push(e.Text[e.Pos])
}

func searchCharF(e *ENV) {
	a, err := e.Strargs.PopE()
	if err != nil {
		return
	}
	for i := e.Pos; i < len(e.Text); i++ {
		c := e.Text[i]
		if string(c) == a.(string) {
			e.Numargs.Push(i - e.Pos)
			return
		}
	}
}

func searchCharB(e *ENV) {
	a, err := e.Strargs.PopE()
	if err != nil {
		return
	}
	for i := e.Pos; i >= 0; i-- {
		c := e.Text[i]
		if string(c) == a.(string) {
			e.Numargs.Push(-i)
			return
		}
	}
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

func insertChar(e *ENV) {
	prev, next := e.Text[:e.Pos], e.Text[e.Pos:]
	a, err := e.Strargs.PopE()
	if err != nil {
		return
	}
	ts := prev
	ts = append(ts, []rune(a.(string))...)
	ts = append(ts, next...)
	e.Text = ts
}

func printChar(e *ENV) {
	if len(e.Text) >= 0 && e.Pos >= 0 && e.Pos < len(e.Text) {
		stdout(string(e.Text[e.Pos]))
	}
}

func repeatCmd(e *ENV) {
	rn, err := e.Numargs.PopE()
	if err != nil {
		return
	}
	cmd, err := e.Strargs.PopE()
	if err != nil {
		return
	}
	if c, ok := COMMANDS[cmd.(string)]; ok {
		for count := 0; count < rn.(int); count++ {
			c(e)
		}
	}
	if m, ok := MACROTABLE[cmd.(string)]; ok {
		for count := 0; count < rn.(int); count++ {
			eval(cmdList([]rune(m)), e)
		}
	}
}

func writeFile(e *ENV) {
	if err := ioutil.WriteFile(e.FName, []byte(string(e.Text)), 0666); err != nil {
		stderr("File write error.")
	}
}

func changeFile(e *ENV) {
	a, err := e.Strargs.PopE()
	if err != nil {
		return
	}
	e.FName = a.(string)
	e.Text = []rune("") // Should we preserve the text from previous file?
}
