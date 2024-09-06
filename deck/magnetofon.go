package magnetofon

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"math"
	"os"
	"time"
)

type Magnetofon struct {
	AudioSource   portaudio.DeviceInfo
	Tape          *Tape
	Stream        *portaudio.Stream
	waitChan      chan os.Signal
	Recording     bool
	Threshold     float64
	timeStart     time.Time
	timeFinish    time.Time
	Bpm           float64
	Steps         int32
	CurrentVolume float64
}

func (m *Magnetofon) TurnOn() {
	portaudio.Initialize()
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, 0, m.CaptureAudio)
	if err != nil {
		panic(err)
	}
	m.Stream = stream
}

func (m *Magnetofon) TurnOff() {
	m.Stream.Close()
	portaudio.Terminate()
}

func (m *Magnetofon) SetAudioSource(audioSource portaudio.DeviceInfo) {
	m.AudioSource = audioSource
}

func (m *Magnetofon) AvailableAudioSources() ([]*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()
	if err != nil {
		return nil, err
	}

	var audioDevices []*portaudio.DeviceInfo
	for _, device := range devices {
		if device.MaxInputChannels < 1 {
			continue
		}
		audioDevices = append(audioDevices, device)
	}

	return audioDevices, nil
}

func (m *Magnetofon) InsertTape(tape *Tape) {
	m.Tape = tape
}

func (m *Magnetofon) StartAudioCapture() {
	_ = m.Stream.Start()
}

func (m *Magnetofon) StopAudioCapture() {
	_ = m.Stream.Stop()
}

func (m *Magnetofon) UntilPatternEnd() {
	waiting := true

	for waiting {
		if time.Now().After(m.timeFinish) {
			waiting = false
		}
	}
}

func (m *Magnetofon) WaitForSignal() {
	waiting := true
	fmt.Println("Waiting for signal")

	for waiting {
		volume := m.calculateVolume(m.Tape.VACBuf)
		m.CurrentVolume = volume

		if volume > m.Threshold {
			waiting = false

			continue
		}
	}

	m.Tape.Buf = m.Tape.VACBuf
}

func (m *Magnetofon) UntilKeyPressed() {
	// wait for the user to press any key using signals

	var b = make([]byte, 1)
	os.Stdin.Read(b)
}

func (m *Magnetofon) CalculateTimeFinish() {
	m.timeStart = time.Now()
	durationMs := (60000.0 / m.Bpm) * float64(m.Steps)
	durationUs := durationMs * 1000 // Convert milliseconds to microseconds
	m.timeFinish = m.timeStart.Add(time.Duration(durationUs) * time.Microsecond)
}

func (m *Magnetofon) StartRecording() {
	m.Recording = true
	m.CalculateTimeFinish()
}

func (m *Magnetofon) StopRecording() {
	m.Recording = false
	actualDuration := time.Since(m.timeStart)
	expectedDuration := m.timeFinish.Sub(m.timeStart)
	lag := actualDuration - expectedDuration
	fmt.Printf("Recording finished. Expected duration: %v, Actual duration: %v, Lag: %v\n", expectedDuration, actualDuration, lag)
}

func (m *Magnetofon) CaptureAudio(in []float32) {

	if !m.Recording {
		m.Tape.VACBuf = in
		return
	}

	//m.CurrentVolume = m.calculateVolume(in)

	m.Tape.Buf = append(m.Tape.Buf, in...)
}

// Helper function to calculate the volume of the audio data
func (m *Magnetofon) calculateVolume(samples []float32) float64 {
	var sum float64
	for _, sample := range samples {
		sum += float64(sample * sample)
	}

	//fmt.Println("Length of samples: ", len(samples))

	mean := float64(sum) / float64(len(samples))
	return math.Sqrt(mean)
}

func (m *Magnetofon) SaveTape() error {
	return m.Tape.Store()
}

func NewMagnetofon() *Magnetofon {
	return &Magnetofon{Threshold: 0.001, Bpm: 100, Steps: 8}
}
