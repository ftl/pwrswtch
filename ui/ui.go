package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

type Application interface {
	PowerLevels() []string
	SetPowerLevel(index int) error
	SetTx(enabled bool) error
}

func Run(app Application) error {
	gtk.Init(nil)

	mainWindow, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return errors.Wrap(err, "cannot create main window")
	}
	mainWindow.SetTitle("Power Switch")
	mainWindow.Connect("destroy", func() {
		gtk.MainQuit()
	})

	grid, err := gtk.GridNew()
	if err != nil {
		return errors.Wrap(err, "cannot create grid")
	}
	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	grid.SetHExpand(true)

	for i, level := range app.PowerLevels() {
		button, err := gtk.ButtonNewWithLabel(level)
		if err != nil {
			return errors.Wrap(err, "cannot create button")
		}
		button.Connect("clicked", powerSetter(app, i))
		grid.Add(button)
	}

	txButton, _ := gtk.ButtonNewWithLabel("Tx")
	txButton.Connect("clicked", txSetter(app, true))
	grid.Add(txButton)
	rxButton, _ := gtk.ButtonNewWithLabel("Rx")
	rxButton.Connect("clicked", txSetter(app, false))
	grid.Add(rxButton)

	mainWindow.Add(grid)
	mainWindow.ShowAll()

	gtk.Main()
	return nil
}

func powerSetter(app Application, index int) func() {
	return func() {
		app.SetPowerLevel(index)
	}
}

func txSetter(app Application, enabled bool) func() {
	return func() {
		app.SetTx(enabled)
	}
}
