package main

import (
	"strconv"
)

// Limited selection of whites, as it acts as a separator, too.
func isWhite(c string) bool {
	switch c {
	case " ", "\t", "\n":
		return true
	default:
		return false
	}
}

func isNum(n string) bool {
	var isPoint = false
	for i := 0; i < len(n); i++ {
		if i == 0 && n[i] == '-' {
			continue
		}
		if !isPoint && n[i] == '.' {
			isPoint = true
			continue
		}
		if n[i] < 48 || n[i] > 57 {
			return false
		}
	}
	return true
}

func convNum(buf string) (float64, error) {
	num, err := strconv.ParseFloat(buf, 64)
	if err != nil {
		return 0.0, err
	}
	return num, nil
}
