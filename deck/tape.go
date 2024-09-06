package magnetofon

import (
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"os"
)

type Tape struct {
	FileName string
	Buf      []float32
	VACBuf   []float32
}

func (t *Tape) Store() error {
	f, err := os.Create(t.FileName)
	if err != nil {
		return err
	}

	defer f.Close()

	enc := wav.NewEncoder(f, 44100, 16, 1, 1)
	defer enc.Close()

	// after converting to int we need to trim the silence from the beginning of the audio data buffer.
	// we can do this by finding the first non-zero sample and then copying the rest of the buffer to a new buffer.
	// we can then write this new buffer to the wav file. Code please:
	// find the first non-zero sample
	var firstNonZeroSample int
	for i, sample := range t.Buf {
		if sample >= 0.002 {
			firstNonZeroSample = i
			break
		}
	}

	fmt.Println("First non-zero sample: ", firstNonZeroSample)

	// copy the rest of the buffer to a new buffer
	trimmedBuf := t.Buf[firstNonZeroSample:]

	intBuf := make([]int, len(trimmedBuf))

	for i, sample := range trimmedBuf {
		intBuf[i] = int(sample * 32768.0)
	}

	err = enc.Write(&audio.IntBuffer{Data: intBuf, Format: &audio.Format{SampleRate: 44100, NumChannels: 1}})

	return err
}

func NewTape(fileName string) *Tape {
	return &Tape{FileName: fileName}
}
