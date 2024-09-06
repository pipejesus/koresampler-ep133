package main

import (
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

	app.Deck.WaitForSignal()
	app.Deck.StartRecording()
	//fmt.Println("Recording started")

	app.Deck.UntilPatternEnd()
	//ui.Run(app.Deck)
	//app.Deck.UntilKeyPressed()
	app.Deck.StopRecording()
	app.Deck.StopAudioCapture()
	app.Deck.SaveTape()
}
