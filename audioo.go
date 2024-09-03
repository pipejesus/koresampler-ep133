package main

import (
	"fmt"
	"github.com/faiface/mainthread"
)

var app *App

func run() {
	// now we can run stuff on the main thread like this
	mainthread.Call(func() {
		app.UI.Run()
	})
}
func main() {
	//runtime.GOMAXPROCS(runtime.NumCPU())

	app = NewApp()
	defer app.Destroy()
	app.PrintAudioSources()

	app.Deck.StartAudioCapture()
	mainthread.Run(run)

	sendChan := make(chan string, 1)
	fmt.Println("Waiting for signal")

	app.Deck.WaitForSignal(sendChan)
	//<-sendChan
	app.Deck.StartRecording()

	app.Deck.UntilPatternEnd(sendChan)
	//<-sendChan
	app.Deck.StopRecording()
	app.Deck.StopAudioCapture()
	app.Deck.SaveTape()
}
