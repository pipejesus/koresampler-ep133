package main

import (
	"fmt"
	"runtime"
)

var app *App

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app = NewApp()
	defer app.Destroy()
	app.PrintAudioSources()

	app.Deck.StartAudioCapture()

	//holdChan := make(chan string, 1)
	//fmt.Println("Waiting for signal")
	ch := make(chan string, 1)
	app.EP133.ListenToMidiMessages(ch)
	<-ch
	//app.Deck.WaitForSignal()
	app.Deck.StartRecording()
	fmt.Println("Recording started")

	app.Deck.UntilPatternEnd()
	//ui.Run(app.Deck)
	//app.Deck.UntilKeyPressed()
	app.Deck.StopRecording()
	app.Deck.StopAudioCapture()
	app.Deck.SaveTape()
}
