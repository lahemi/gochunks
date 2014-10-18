package main

// Limited selection of whites, as it acts as a separator, too.
func isWhite(c string) bool {
	switch c {
	case " ", "\t", "\n":
		return true
	default:
		return false
	}
}

func isFloat(n string) bool {
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

func isInt(n string) bool {
	for i := 0; i < len(n); i++ {
		if i == 0 && n[i] == '-' {
			continue
		}
		if n[i] < 48 || n[i] > 57 {
			return false
		}
	}
	return true
}

// Not general.
func isString(s string) bool {
	if s[0] == '"' && s[len(s)-1] == '"' {
		return true
	}
	return false
}
