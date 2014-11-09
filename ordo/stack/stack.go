package stack

import "errors"

var emptyStackError = errors.New("Attempting to Pop an empty stack")

type Stack struct {
	s []interface{}
}

func (s *Stack) Push(a interface{}) {
	s.s = append(s.s, a)
}

func (s *Stack) PopE() (interface{}, error) {
	if len(s.s) == 0 {
		return nil, emptyStackError
	}
	last := s.s[len(s.s)-1]
	s.s = s.s[:len(s.s)-1]
	return last, nil
}

func (s *Stack) Pop() interface{} {
	rv, err := s.PopE()
	if err != nil {
		return nil
	}
	return rv
}

func (s *Stack) Dup() {
	rv, err := s.PopE()
	if err != nil {
		return
	}
	s.Push(rv)
	s.Push(rv)
}

func (s *Stack) Drop() {
	s.Pop()
}

func (s *Stack) Swap() {
	r2, err2 := s.PopE()
	r1, err1 := s.PopE()
	if err2 != nil || err1 != nil {
		return
	}
	s.Push(r2)
	s.Push(r1)
}

func (s *Stack) Over() {
	r2, err2 := s.PopE()
	r1, err1 := s.PopE()
	if err2 != nil || err1 != nil {
		return
	}
	s.Push(r1)
	s.Push(r2)
	s.Push(r1)
}

func (s *Stack) Rot() {
	r3, err3 := s.PopE()
	r2, err2 := s.PopE()
	r1, err1 := s.PopE()
	if err3 != nil || err2 != nil || err1 != nil {
		return
	}
	s.Push(r2)
	s.Push(r3)
	s.Push(r1)
}
