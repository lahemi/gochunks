package main

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
