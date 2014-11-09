package main

func arithArgs(e *ENV) (int, int) {
	a2, err1 := e.Nums.PopE()
	a1, err2 := e.Nums.PopE()
	if err1 != nil || err2 != nil {
		stderr("Not enough args for arithmetic.\n")
		return 1, -1
	}
	return a1.(int), a2.(int)
}

func addition(e *ENV) {
	a1, a2 := arithArgs(e)
	e.Nums.Push(a1 + a2)
}

func subtraction(e *ENV) {
	a1, a2 := arithArgs(e)
	e.Nums.Push(a1 - a2)
}

func multiplication(e *ENV) {
	a1, a2 := arithArgs(e)
	e.Nums.Push(a1 * a2)
}

func division(e *ENV) {
	a1, a2 := arithArgs(e)
	if a2 == 0 {
		stderr("Division by zero.\n")
		return
	}
	e.Nums.Push(a1 / a2)
}
