package magnetofon

import (
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

	intBuf := make([]int, len(t.Buf))
	for i, sample := range t.Buf {
		intBuf[i] = int(sample * 32768.0)
	}

	err = enc.Write(&audio.IntBuffer{Data: intBuf, Format: &audio.Format{SampleRate: 44100, NumChannels: 1}})

	return err
}

func NewTape(fileName string) *Tape {
	return &Tape{FileName: fileName}
}
