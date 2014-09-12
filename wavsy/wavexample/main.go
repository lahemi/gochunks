// Creates a short concert A pitched sine wave.
package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"wavsy"
)

const num_samples = wavsy.Sample_per_sec * 2

func main() {

	var (
		waveform  [num_samples]int32
		frequency = 440.0
		volume    = 32000.0
		length    = num_samples
        i int32
	)

	for i = 0; i < length; i++ {
		t := float64(i) / float64(wavsy.Sample_per_sec)
		waveform[i] = int32(volume * math.Sin(frequency*t*2.0*math.Pi))
	}

	buf := new(bytes.Buffer)
	for i := 0; i < len(waveform); i++ {
		binary.Write(buf, binary.LittleEndian, waveform[i])
	}

	file := wavsy.Wav_open("sound.wav")
	wavsy.Wav_write(file, buf.Bytes())
	wavsy.Wav_close(file)
}
