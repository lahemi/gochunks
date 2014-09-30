// Br**nfuck. Not clean, panics.
// Usage: prog <file.bf>
package main

import (
	"io/ioutil"
	"os"
)

type Stack struct {
	st []int
}

func (s *Stack) push(i int) {
	s.st = append(s.st, i)
}
func (s *Stack) pop() int {
	last := s.st[len(s.st)-1]
	s.st = s.st[:len(s.st)-1] // Doesn't work with `last`
	return last
}

func interpret(code []byte) {
	var (
		data [100000]int
		jump [100000]int
		dp   = 0 // Data Pointer

		stack = Stack{}
	)

	// Pre-compile a jump table
	for i := 0; i < len(code); i++ {
		switch code[i] {
		case '[':
			stack.push(i)
		case ']':
			jump[i] = stack.pop()
			jump[jump[i]] = i // This is quite clever, actually.
		}
	}

	// Execute. cp == Code Pointer
	for cp := 0; cp < len(code); cp++ {
		switch code[cp] {
		case '>':
			dp += 1
		case '<':
			dp -= 1
		case '+':
			data[dp] += 1
		case '-':
			data[dp] -= 1
		case '.':
			os.Stdout.Write([]byte(string(data[dp])))
		case ',':
			b := make([]byte, 1)
			_, _ = os.Stdin.Read(b)
			data[dp] = int(b[0])
		case '[':
			if data[dp] == 0 {
				cp = jump[cp]
			}
		case ']':
			if data[dp] != 0 {
				cp = jump[cp]
			}
		}
	}
}

func main() {
	cnt, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	interpret(cnt)
}
