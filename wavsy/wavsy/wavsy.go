// A little bit silly, and barebones .wav file creation functions:
// Wav_open : which also creates the header.
// Wav_write: which is just a wrap around file.Write([]byte)
// Wav_close: which also checks the final filesize and
//            places the correct size in the header.
//
// There are panics in the code, so if you want to use this as
// a base for something properer, that'd be the place to start.
//
// There is also a neat little, unsafe, trick in Wav_open
// for making your structs available as a []byte as well!
//
// 12-Sep-2014, Lauri Peltomäki
package wavsy

import (
	"os"
	"unsafe"
)

const Sample_per_sec = int32(44100)

// http://www.topherlee.com/software/pcm-tut-wavformat.html
type wav_header struct {
	Riff_tag        [4]uint8 //1 - 4   "RIFF"  Marks the file as a riff file. Characters are each 1 byte long.
	Riff_len        int32    //5 - 8   File size (integer) Size of the overall file - 8 bytes, in bytes (32-bit integer). Typically, you'd fill this in after creation.
	Wave_tag        [4]uint8 //9 -12   "WAVE"  File Type Header. For our purposes, it always equals "WAVE".
	Fmt_tag         [4]uint8 //13-16   "fmt "  Format chunk marker. Includes trailing null
	Fmt_len         uint32   //17-20   16  Length of format data as listed above
	Audio_fmt       uint16   //21-22   1   Type of format (1 is PCM) - 2 byte integer
	Num_chans       uint16   //23-24   2   Number of Channels - 2 byte integer
	Sample_rate     uint32   //25-28   44100   Sample Rate - 32 byte integer. Common values are 44100 (CD), 48000 (DAT). Sample Rate = Number of Samples per second, or Hertz.
	Byte_rate       uint32   //29-32   176400  (Sample Rate * BitsPerSample * Channels) / 8.
	Block_align     uint16   //33-34   4   (BitsPerSample * Channels) / 8.1 - 8 bit mono2 - 8 bit stereo/16 bit mono4 - 16 bit stereo
	Bits_per_sample uint16   //35-36   16    bits per sample
	Data_tag        [4]uint8 //37-40   "data" chunk header. Marks the beginning of the data section.
	Data_len        uint32   //41-44   File size (data)    Size of the data section.
}

var header wav_header

func Wav_open(filename string) *os.File {

	var (
		srate = uint32(Sample_per_sec)
		bitss = uint16(16)
	)

	// A trick so that we get the header as a []byte
	buf := make([]byte, int(unsafe.Sizeof(header)))

	// A pointer to the start of the buffer and cast it to header
	sp := (*wav_header)(unsafe.Pointer(&buf[0]))

	// And then use the pointer for assignments.
	sp.Riff_tag = [4]uint8{'R', 'I', 'F', 'F'}
	sp.Wave_tag = [4]uint8{'W', 'A', 'V', 'E'}
	sp.Fmt_tag = [4]uint8{'f', 'm', 't', ' '}
	sp.Data_tag = [4]uint8{'d', 'a', 't', 'a'}

	sp.Riff_len = 0
	sp.Fmt_len = 16
	sp.Audio_fmt = 1
	sp.Num_chans = 1
	sp.Sample_rate = srate
	sp.Byte_rate = srate * (uint32(bitss) / 8)
	sp.Block_align = bitss / 8
	sp.Bits_per_sample = bitss
	sp.Data_len = 0

	//file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	if _, err := file.Write(buf); err != nil {
		panic(err)
	}

	return file
}

// ..., maybe take a bytes.Buffer instead of
// []byte and do the "conversion" here ?
func Wav_write(file *os.File, data []byte) {
	file.Write(data)
}

func Wav_close(file *os.File) {
	// not just close, but also set the size of the file appropriate

	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}
	file_len := fi.Size()

	data_len := file_len - int64(unsafe.Sizeof(header))
	file.Seek(int64(unsafe.Sizeof(header))-4, os.SEEK_SET) // size of int32
	if _, err := file.Write([]byte(string(data_len))); err != nil {
		panic(err)
	}

	riff_len := file_len - 8
	file.Seek(4, os.SEEK_SET)
	file.Write([]byte(string(riff_len)))

	file.Close()
}
