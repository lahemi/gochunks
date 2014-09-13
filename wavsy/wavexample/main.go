package main

import (
	"math"
	"wavsy"
)

const num_samples = wavsy.Sample_per_sec * 2

func main() {

	var (
		waveform  []int32
		frequency = 440.0
		volume    = 32000.0
		length    = num_samples
		i         int32
	)

	for i = 0; i < length; i++ {
		t := float64(i) / float64(wavsy.Sample_per_sec)
		waveform = append(waveform, int32(volume*math.Sin(frequency*t*2.0*math.Pi)))
	}

	file := wavsy.WavOpen("sound.wav")
	wavsy.WavWrite(file, waveform)
	wavsy.WavClose(file)
}
