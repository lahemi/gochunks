package main

type Fstack []float64

func (s *Fstack) Push(f float64) {
	(*s) = append((*s), f)
}
func (s *Fstack) Pop() float64 {
	if len((*s)) <= 0 {
		return -0xffffffff // pseudo-error
	}
	last := (*s)[len((*s))-1]
	(*s) = (*s)[:len((*s))-1]
	return last
}
