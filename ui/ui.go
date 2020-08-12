package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ftl/gmtry"
	"github.com/ftl/hamradio/cfg"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type Application interface {
	PowerLevels() []string
	SetPowerLevel(index int) error
	SetTx(enabled bool) error
	Tuning() bool
	ToggleTuning() (bool, error)
}

type ui struct {
	application Application
	gtkapp      *gtk.Application
	mainWindow  *gtk.ApplicationWindow
	geometry    *gmtry.Geometry
}

func (ui *ui) onActivate() {
	var err error
	ui.mainWindow, err = gtk.ApplicationWindowNew(ui.gtkapp)
	if err != nil {
		log.Fatalf("Cannot create main window: %v", err)
	}
	ui.mainWindow.SetTitle("Power Switch")

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatalf("Cannot create grid: %v", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_HORIZONTAL)
	grid.SetHExpand(true)

	for i, level := range ui.application.PowerLevels() {
		button, err := gtk.ButtonNewWithLabel(level)
		if err != nil {
			log.Fatalf("Cannot create button: %v", err)
		}
		button.Connect("clicked", powerSetter(ui.application, i))
		grid.Add(button)
	}

	rxButton, _ := gtk.ButtonNewWithLabel("Rx")
	rxButton.Connect("clicked", txSetter(ui.application, false))
	grid.Add(rxButton)
	tuneButton, _ := gtk.ButtonNewWithLabel("Tune")
	tuneButton.Connect("clicked", tuningToggler(ui.application, tuneButton))
	tuneButton.SetHExpand(true)
	grid.Add(tuneButton)

	ui.mainWindow.Add(grid)
	ui.mainWindow.SetDefaultSize(100, 100)

	connectToGeometry(ui.geometry, "main", &ui.mainWindow.Window)
	err = ui.geometry.Restore()
	if err != nil {
		log.Printf("Cannot restore the window geometry: %v", err)
	}

	ui.mainWindow.ShowAll()
}

func Run(app Application) error {
	var err error

	gdk.SetAllowedBackends("x11")

	ui := new(ui)
	ui.application = app
	ui.gtkapp, err = gtk.ApplicationNew("ft.pwrswtch", glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		return fmt.Errorf("cannot create application: %v", err)
	}
	ui.gtkapp.Connect("activate", ui.onActivate)

	configDir, err := cfg.Directory("")
	if err != nil {
		log.Fatalf("No access to configuration directory %s: %v", cfg.DefaultDirectory, err)
	}
	filename := filepath.Join(configDir, "pwrswtch.geometry")
	ui.geometry = gmtry.NewGeometry(filename)

	ui.gtkapp.Run(os.Args)

	err = ui.geometry.Store()
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
		w.SetPosition(window.GetPosition()) // in this setup, the event contains the y of the area below the title bar, not the y of the top-left corner
		w.SetSize(e.Width(), e.Height())
	})
	window.Connect("window-state-event", func(_ interface{}, event *gdk.Event) {
		e := gdk.EventWindowStateNewFromEvent(event)
		if e.ChangedMask()&gdk.WINDOW_STATE_MAXIMIZED == gdk.WINDOW_STATE_MAXIMIZED {
			geometry.Get(id).SetMaximized(e.NewWindowState()&gdk.WINDOW_STATE_MAXIMIZED == gdk.WINDOW_STATE_MAXIMIZED)
		}
	})
}
