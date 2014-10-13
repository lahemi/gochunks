package main

import (
	"fmt"
)

func bitap(text, pattern string) string {

	var (
		m                 = uint(len(pattern))
		maxval       uint = 255
		pattern_mask      = make([]int, maxval+1)
		R                 = ^1
		i            uint
	)

	if m == 0 {
		return text
	}
	if m > 31 {
		return "Nope len"
	}

	// Init the pattern bitmasks
	for i = 0; i <= maxval; i++ {
		pattern_mask[i] = ^0
	}
	for i = 0; i < m; i++ {
		pattern_mask[pattern[i]] &= ^(1 << i)
	}
	for i = 0; i < uint(len(text)); i++ {
		R |= pattern_mask[text[i]]
		R <<= 1

		if (R & (1 << m)) == 0 {
			return string(text[i-m+1:])
		}
	}

	return ""
}

func bitapFuzzy(text, pattern string, k uint) string {
	var (
		m                 = uint(len(pattern))
		maxval       uint = 255
		pattern_mask      = make([]int, maxval+1)
		R                 = make([]int, maxval+1)
		i, d         uint
	)

	if m == 0 {
		return text
	}
	if m > 31 {
		return "Nope len"
	}

	// Init bit array
	for i = 0; i <= k; i++ {
		R[i] = ^1
	}

	// Init pattern bitmasks
	for i = 0; i <= maxval; i++ {
		pattern_mask[i] = ^0
	}
	for i = 0; i < m; i++ {
		pattern_mask[pattern[i]] &= ^(1 << i)
	}

	for i = 0; i < uint(len(text)); i++ {
		var old_Rd1 int = R[0]

		R[0] |= pattern_mask[text[i]]
		R[0] <<= 1

		for d = 1; d <= i; d++ {
			var tmp int = R[d]
			// Substitution is all we care about.
			R[d] = (old_Rd1 & (R[d] | pattern_mask[text[i]])) << 1
			old_Rd1 = tmp
		}

		if (R[k] & (1 << m)) == 0 {
			return string(text[i-m+1:])
		}
	}

	return ""
}

func main() {
	text := "This is a text."
	fmt.Println(bitap(text, " is "))
	fmt.Println(bitapFuzzy(text, "hitm", 4))
}
