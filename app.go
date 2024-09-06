package main

import (
	"fmt"
	magnetofon "koresampler/deck"
)

type App struct {
	Deck  *magnetofon.Magnetofon
	EP133 *magnetofon.EP133
}

func (a *App) PrintAudioSources() {
	devices, _ := a.Deck.AvailableAudioSources()
	for i, device := range devices {
		fmt.Printf("%d: %s\n", i, device.Name)
	}
}

func (a *App) Destroy() {
	a.Deck.TurnOff()
}

func NewApp() *App {
	deck := magnetofon.NewMagnetofon()
	deck.InsertTape(magnetofon.NewTape("fifi.wav"))
	deck.TurnOn()

	ep133 := magnetofon.NewEP133()

	return &App{Deck: deck, EP133: ep133}
}
