package magnetofon

import (
	"fmt"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)
import _ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"

const DeviceName = "EP-133"
const MsgStart = "Continue"
const MsgStop = "Stop"

type EP133 struct {
	InPort  drivers.In
	OutPort drivers.Out
	Stop    func()
	Send    func(midi.Message) error
}

func NewEP133() *EP133 {
	return (&EP133{}).FindDevice()
}

func (e *EP133) FindDevice() *EP133 {
	in, err := midi.FindInPort(DeviceName)
	if err != nil {
		return e
	}
	e.InPort = in

	e.OutPort, err = midi.FindOutPort(DeviceName)
	if err != nil {
		return e
	}

	sendFunc, err := midi.SendTo(e.OutPort)
	if err != nil {
		return e
	}

	e.Send = sendFunc

	return e
}

func (e *EP133) CheckDevice() bool {
	return e.InPort != nil && e.OutPort != nil
}

func (e *EP133) CloseDevice() {
	if e.Stop != nil {
		e.Stop()
	}

	midi.CloseDriver()
}

// StartPlayback sends a MIDI Continue message to EP-133
// to start the playback from the beginning of the pattern.
func (e *EP133) StartPlayback() {
	if e.Send != nil {
		_ = e.Send(midi.Continue())
	}
}

func (e *EP133) ListenToContinueAndStop(sendChannel chan string) {
	if e.InPort == nil {
		fmt.Println("No device found")
		return
	}

	stopFunc, err := midi.ListenTo(e.InPort, func(msg midi.Message, timestampms int32) {

		//if msg.Is(midi.TimingClockMsg) {
		//	fmt.Println("KLOK!")
		//}

		messageStr := msg.Type().String()
		switch messageStr {
		case MsgStart:
			sendChannel <- MsgStart
			fmt.Println("Received START from EP-133!")
		case MsgStop:
			sendChannel <- MsgStop
			fmt.Println("Received STOP from EP-133!")
		default:
			fmt.Printf("Received unknown message: %s\n", messageStr)
		}
	}, midi.UseSysEx(), midi.UseTimeCode())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	e.Stop = stopFunc
}
