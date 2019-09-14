package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gtk"

	"github.com/ftl/pwrswtch/core/trx"
)


type powerLevel struct {
	watts int
	value float32
}

var powerLevels = []powerLevel{
	{10, 0.039216},
	{30, 0.039216}, // TODO use real value
	{50, 0.039216}, // TODO use real value
	{100, 1.0},
}

type tx interface {
	SetPowerLevel(value float32) error
	SetTx(enabled bool) error
}

func main() {
	tx := trx.Open("")

	gtk.Init(nil)

	mainWindow, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Cannot create main window:", err)
	}
	mainWindow.SetTitle("Power Switch")
	mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Cannot create grid: ", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	grid.SetHExpand(true)

	for _, level := range powerLevels {
		label := fmt.Sprintf("%dW", level.watts)
		button, err := gtk.ButtonNewWithLabel(label)
		if err != nil {
			log.Fatal("Cannot create button: ", err)
		}
		button.Connect("clicked", powerSetter(tx, level.value))
		grid.Add(button)
	}
	
	txButton, _ := gtk.ButtonNewWithLabel("Tx")
	txButton.Connect("clicked", txSetter(tx, true))
	grid.Add(txButton)
	rxButton, _ := gtk.ButtonNewWithLabel("Rx")
	rxButton.Connect("clicked", txSetter(tx, false))
	grid.Add(rxButton)

	mainWindow.Add(grid)
	mainWindow.ShowAll()

	gtk.Main()
}

func powerSetter(tx tx, value float32) func() {
	return func() {
		tx.SetPowerLevel(value)
	}
}

func txSetter(tx tx, enabled bool) func() {
	return func() {
		tx.SetTx(enabled)
	}
}