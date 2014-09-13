package frequencisize

import (
	"math"
)

const (
	fixed_A  = 440.0           // concert A, f(0)
	magick_A = 1.0594630943593 // twelth root of 2, half step frequency ratio
)

// All caps so that lookup is simpler. See ByString().
var Notes = map[string]float64{
	"C0": 16.35, "C#0": 17.32, "DB0": 17.32, "D0": 18.35,
	"D#0": 19.45, "EB0": 19.45, "E0": 20.60, "F0": 21.83,
	"F#0": 23.12, "GB0": 23.12, "G0": 24.50, "G#0": 25.96,
	"AB0": 25.96, "A0": 27.50, "A#0": 29.14, "BB0": 29.14,
	"B0": 30.87, "C1": 32.70, "C#1": 34.65, "DB1": 34.65,
	"D1": 36.71, "D#1": 38.89, "EB1": 38.89, "E1": 41.20,
	"F1": 43.65, "F#1": 46.25, "GB1": 46.25, "G1": 49.00,
	"G#1": 51.91, "AB1": 51.91, "A1": 55.00, "A#1": 58.27,
	"BB1": 58.27, "B1": 61.74, "C2": 65.41, "C#2": 69.30,
	"DB2": 69.30, "D2": 73.42, "D#2": 77.78, "EB2": 77.78,
	"E2": 82.41, "F2": 87.31, "F#2": 92.50, "GB2": 92.50,
	"G2": 98.00, "G#2": 103.83, "AB2": 103.83, "A2": 110.00,
	"A#2": 116.54, "BB2": 116.54, "B2": 123.47, "C3": 130.81,
	"C#3": 138.59, "DB3": 138.59, "D3": 146.83, "D#3": 155.56,
	"EB3": 155.56, "E3": 164.81, "F3": 174.61, "F#3": 185.00,
	"GB3": 185.00, "G3": 196.00, "G#3": 207.65, "AB3": 207.65,
	"A3": 220.00, "A#3": 233.08, "BB3": 233.08, "B3": 246.94,
	"C4": 261.63, "C#4": 277.18, "DB4": 277.18, "D4": 293.66,
	"D#4": 311.13, "EB4": 311.13, "E4": 329.63, "F4": 349.23,
	"F#4": 369.99, "GB4": 369.99, "G4": 392.00, "G#4": 415.30,
	"AB4": 415.30, "A4": 440.00, "A#4": 466.16, "BB4": 466.16,
	"B4": 493.88, "C5": 523.25, "C#5": 554.37, "DB5": 554.37,
	"D5": 587.33, "D#5": 622.25, "EB5": 622.25, "E5": 659.25,
	"F5": 698.46, "F#5": 739.99, "GB5": 739.99, "G5": 783.99,
	"G#5": 830.61, "AB5": 830.61, "A5": 880.00, "A#5": 932.33,
	"BB5": 932.33, "B5": 987.77, "C6": 1046.50, "C#6": 1108.73,
	"DB6": 1108.73, "D6": 1174.66, "D#6": 1244.51, "EB6": 1244.51,
	"E6": 1318.51, "F6": 1396.91, "F#6": 1479.98, "GB6": 1479.98,
	"G6": 1567.98, "G#6": 1661.22, "AB6": 1661.22, "A6": 1760.00,
	"A#6": 1864.66, "BB6": 1864.66, "B6": 1975.53, "C7": 2093.00,
	"C#7": 2217.46, "DB7": 2217.46, "D7": 2349.32, "D#7": 2489.02,
	"EB7": 2489.02, "E7": 2637.02, "F7": 2793.83, "F#7": 2959.96,
	"GB7": 2959.96, "G7": 3135.96, "G#7": 3322.44, "AB7": 3322.44,
	"A7": 3520.00, "A#7": 3729.31, "BB7": 3729.31, "B7": 3951.07,
	"C8": 4186.01, "C#8": 4434.92, "DB8": 4434.92, "D8": 4698.63,
	"D#8": 4978.03, "EB8": 4978.03, "E8": 5274.04, "F8": 5587.65,
	"F#8": 5919.91, "GB8": 5919.91, "G8": 6271.93, "G#8": 6644.88,
	"AB8": 6644.88, "A8": 7040.00, "A#8": 7458.62, "BB8": 7458.62,
	"B8": 7902.13,
}

// http://www.phy.mtu.edu/~suits/NoteFreqCalcs.html
// f for frequency of the note
// n is the number of half steps from f0, higher positive, lower negative.
// f(0) == 440 Hz
// f(n) <- f(0) * a^n
//      where a := 2^(1/12)
func ByHalfStep(steps int) float64 {
	return fixed_A * math.Pow(magick_A, float64(steps))
}

func ByList(n []string) (freqs []float64) {
	for i := 0; i < len(n); i++ {
		freqs = append(freqs, Notes[n[i]])
	}
	return
}

// Treats whitespace, tab and newline as the separator.
// Ignores extraneous blanks and incorrent input.
// Accepts both cases, ie. Ab BB cB are all correct input.
func ByString(n string) (freqs []float64) {

	var tstr string

	isWhite := func(c byte) bool {
		if c == ' ' || c == '\n' || c == '\t' {
			return true
		}
		return false
	}
	isValid := func(c byte) bool {
		if (c >= 65 && c <= 90) ||
			(c >= 97 && c <= 122) ||
			(c >= 48 && c <= 57) ||
			c == 35 {
			return true
		}
		return false
	}

	for i := 0; i < len(n); i++ {
		s := n[i]
		switch {
		case isWhite(s):
			if len(tstr) > 0 {
				if _, ok := Notes[tstr]; ok {
					freqs = append(freqs, Notes[tstr])
				}
				tstr = ""
			}
		case isValid(s):
			if 'a' <= s && s <= 'z' {
				s -= 'a' - 'A'
			}
			tstr += string(s)
		}
	}
	freqs = append(freqs, Notes[tstr])

	return
}
