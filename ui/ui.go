package ui

import (
	"log"
	"path/filepath"

	"github.com/ftl/gmtry"
	"github.com/ftl/hamradio/cfg"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/errors"
)

type Application interface {
	PowerLevels() []string
	SetPowerLevel(index int) error
	SetTx(enabled bool) error
	Tuning() bool
	ToggleTuning() (bool, error)
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

	rxButton, _ := gtk.ButtonNewWithLabel("Rx")
	rxButton.Connect("clicked", txSetter(app, false))
	grid.Add(rxButton)
	tuneButton, _ := gtk.ButtonNewWithLabel("Tune")
	tuneButton.Connect("clicked", tuningToggler(app, tuneButton))
	tuneButton.SetHExpand(true)
	grid.Add(tuneButton)

	mainWindow.Add(grid)

	configDir, err := cfg.Directory("")
	if err != nil {
		log.Fatalf("No access to configuration directory %s: %v", cfg.DefaultDirectory, err)
	}
	filename := filepath.Join(configDir, "pwrstch.geometry")
	geometry := gmtry.NewGeometry(filename)
	connectToGeometry(geometry, "main", mainWindow)
	err = geometry.Restore()
	if err != nil {
		log.Printf("Cannot restore the window geometry: %v", err)
	}

	mainWindow.ShowAll()

	gtk.Main()

	err = geometry.Store()
	if err != nil {
		log.Printf("Cannot store the window geometry: %v", err)
	}
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

func tuningToggler(app Application, button *gtk.Button) func() {
	return func() {
		tuning, _ := app.ToggleTuning()
		if tuning {
			button.SetLabel("Tuning")
		} else {
			button.SetLabel("Tune")
		}
	}
}

func connectToGeometry(geometry *gmtry.Geometry, id gmtry.ID, window *gtk.Window) {
	geometry.Add(id, window)

	window.Connect("configure-event", func(_ interface{}, event *gdk.Event) {
		e := gdk.EventConfigureNewFromEvent(event)
		w := geometry.Get(id)
		w.SetPosition(e.X(), e.Y())
		w.SetSize(e.Width(), e.Height())
	})
	window.Connect("window-state-event", func(_ interface{}, event *gdk.Event) {
		e := gdk.EventWindowStateNewFromEvent(event)
		if e.ChangedMask()&gdk.WINDOW_STATE_MAXIMIZED == gdk.WINDOW_STATE_MAXIMIZED {
			geometry.Get(id).SetMaximized(e.NewWindowState()&gdk.WINDOW_STATE_MAXIMIZED == gdk.WINDOW_STATE_MAXIMIZED)
		}
	})
}
