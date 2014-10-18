package main

import (
	"fmt"
)

var BUILTINS = map[string]func(env *ENV){
	// + - × ÷ < ≤ = ≥ > ≠
	// ⌊ ⌈
	// ∧ ∨ ~
	// ○
	// ↑ ↓
	// « » ⍲ ⍱
	// / º
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

	},
	".S": func(env *ENV) {
		fmt.Println(env)
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
