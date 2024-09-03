package main

import (
	magnetofon "audioo/deck"
	"audioo/visual"
	"fmt"
)

type App struct {
	Deck  *magnetofon.Magnetofon
	EP133 *magnetofon.EP133
	UI    *visual.UI
}

func (a *App) PrintAudioSources() {
	devices, _ := a.Deck.AvailableAudioSources()
	for i, device := range devices {
		fmt.Printf("%d: %s\n", i, device.Name)
	}
}

func (a *App) Destroy() {
	a.UI.Destroy()
	a.Deck.TurnOff()
}

func NewApp() *App {
	deck := magnetofon.NewMagnetofon()
	deck.InsertTape(magnetofon.NewTape("fifi.wav"))
	deck.TurnOn()

	ep133 := magnetofon.NewEP133()

	ui := visual.NewUI()
	ui.Init()

	return &App{Deck: deck, EP133: ep133, UI: ui}
}
