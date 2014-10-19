package main

import (
	"fmt"
)

// This is what you get with no meta programming and macros.
func comparisonOpers(comp string, env *ENV) {
	var a1, a2 float64
	switch env.MODE {
	case INT:
		a2, a1 = float64(env.INTS.Pop()), float64(env.INTS.Pop())
	case FLOAT:
		a2, a1 = env.FLOATS.Pop(), env.FLOATS.Pop()
	}
	switch comp {
	case "<":
		if a1 < a2 {
			env.INTS.Push(1)
			return
		}
	case "≤":
		if a1 <= a2 {
			env.INTS.Push(1)
			return
		}
	case "=":
		if a1 == a2 {
			env.INTS.Push(1)
			return
		}
	case "≠":
		if a1 != a2 {
			env.INTS.Push(1)
			return
		}
	case "≥":
		if a1 >= a2 {
			env.INTS.Push(1)
			return
		}
	case ">":
		if a1 > a2 {
			env.INTS.Push(1)
			return
		}
	}
	env.INTS.Push(0) // if it ever reaches here, it must be false
}

var BUILTINS = map[string]func(env *ENV){
	// < ≤ = ≥ > ≠
	// ⌊ ⌈
	// ∧ ∨ ~
	// ○
	// ↑ ↓
	// « » ⍲ ⍱
	// / º
	//
	// SWAP, OVER, ROT, PICK
	//
	"INT": func(env *ENV) {
		env.SetMode(INT)
	},
	"FLOAT": func(env *ENV) {
		env.SetMode(FLOAT)
	},
	"STRING": func(env *ENV) {
		env.SetMode(STRING)
	},
	"PRINT": func(env *ENV) {
		switch env.MODE {
		case INT:
			a1 := env.INTS.Pop()
			fmt.Print(a1)
		case FLOAT:
			a1 := env.FLOATS.Pop()
			fmt.Print(a1)
		case STRING:
			a1 := env.STRINGS.Pop()
			fmt.Print(a1)
		}
	},
	".S": func(env *ENV) {
		fmt.Print(env)
	},
	"NL": func(env *ENV) {
		fmt.Println()
	},

	// Got to figure a better way to do this...
	// Wish I had some proper macros.
	"+": func(env *ENV) {
		switch env.MODE {
		case INT:
			a2, a1 := env.INTS.Pop(), env.INTS.Pop()
			env.INTS.Push(a1 + a2)
		case FLOAT:
			a2, a1 := env.FLOATS.Pop(), env.FLOATS.Pop()
			env.FLOATS.Push(a1 + a2)
		case STRING:
			a2, a1 := env.STRINGS.Pop(), env.STRINGS.Pop()
			env.STRINGS.Push(a1 + a2)
		}
	},
	"-": func(env *ENV) {
		switch env.MODE {
		case INT:
			a2, a1 := env.INTS.Pop(), env.INTS.Pop()
			env.INTS.Push(a1 - a2)
		case FLOAT:
			a2, a1 := env.FLOATS.Pop(), env.FLOATS.Pop()
			env.FLOATS.Push(a1 - a2)
		}
	},
	"×": func(env *ENV) {
		switch env.MODE {
		case INT:
			a2, a1 := env.INTS.Pop(), env.INTS.Pop()
			env.INTS.Push(a1 * a2)
		case FLOAT:
			a2, a1 := env.FLOATS.Pop(), env.FLOATS.Pop()
			env.FLOATS.Push(a1 * a2)
		}
	},
	"÷": func(env *ENV) {
		switch env.MODE {
		case INT:
			a2, a1 := env.INTS.Pop(), env.INTS.Pop()
			env.INTS.Push(a1 / a2)
		case FLOAT:
			a2, a1 := env.FLOATS.Pop(), env.FLOATS.Pop()
			env.FLOATS.Push(a1 / a2)
		}
	},

	"<": func(env *ENV) {
		comparisonOpers("<", env)
	},
	"≤": func(env *ENV) {
		comparisonOpers("≤", env)
	},
	"=": func(env *ENV) {
		comparisonOpers("=", env)
	},
	"≠": func(env *ENV) {
		comparisonOpers("≠", env)
	},
	"≥": func(env *ENV) {
		comparisonOpers("≥", env)
	},
	">": func(env *ENV) {
		comparisonOpers(">", env)
	},

	"DUP": func(env *ENV) {
		switch env.MODE {
		case INT:
			a1 := env.INTS.Pop()
			env.INTS.Push(a1)
			env.INTS.Push(a1)
		case FLOAT:
			a1 := env.FLOATS.Pop()
			env.FLOATS.Push(a1)
			env.FLOATS.Push(a1)
		case STRING:
			a1 := env.STRINGS.Pop()
			env.STRINGS.Push(a1)
			env.STRINGS.Push(a1)
		}
	},
	"DROP": func(env *ENV) {
		switch env.MODE {
		case INT:
			env.INTS.Pop()
		case FLOAT:
			env.FLOATS.Pop()
		case STRING:
			env.STRINGS.Pop()
		}
	},
}
