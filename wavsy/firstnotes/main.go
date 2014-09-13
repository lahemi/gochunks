// Creates notes.wav, containing a few different notes.
package main

import (
	freqs "frequencisize"
	"math"
	"wavsy"
)

func main() {
	// Because of the possibly uneven duration length, and int is a must.
	odd := func(n float64) int32 {
		if math.Mod(n, 1) != 0 && math.Mod(n, 0.5) == 0 {
			return int32(0.5 + n)
		}
		return int32(n)
	}

	//notelist := []string{"C4","D4","E4","A4","C5"}
	var (
		//frequs   = freqs.ByList(notelist)
		// See comments on frequencisize.ByString
		frequs   = freqs.ByString("C4   d4  D#4, f4 a6# Ab4 cmts Bb4 F4")
		duration = odd(float64(len(frequs)))
		data_len = wavsy.Sample_per_sec * duration
		waveform []int32
		volume   = 32000.0
		i, j     int32
	)

	for j = 0; j < duration; j++ {
		frequ := frequs[j]
		omega := 2.0 * math.Pi * frequ
		for i = 0; i < data_len/duration; i++ {
			y := volume * math.Sin(omega*float64(i)/float64(wavsy.Sample_per_sec))
			waveform = append(waveform, int32(y))
		}
	}

	file := wavsy.Wav_open("notes.wav")
	wavsy.Wav_write(file, waveform)
	wavsy.Wav_close(file)
}
