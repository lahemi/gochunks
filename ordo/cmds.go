package main

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
	stdout(string(e.Text[e.Pos]))
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
			eval(cmdList([]byte(m)), e)
		}
	}
}
