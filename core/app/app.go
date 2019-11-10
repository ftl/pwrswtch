package app

import (
	"fmt"

	"github.com/ftl/pwrswtch/core"
)

func New(trx TRX, config core.Configuration) *App {
	result := App{
		trx:         trx,
		powerLevels: config.Levels,
		tuningValue: config.TuningValue,
	}

	result.powerLevelLabels = make([]string, len(result.powerLevels))
	for i, level := range result.powerLevels {
		result.powerLevelLabels[i] = fmt.Sprintf("%dW", level.Watts)
	}
	return &result
}

type App struct {
	trx              TRX
	powerLevels      []powerLevel
	powerLevelLabels []string
	tuning           bool
	tuningValue      string
	lastValue        string
	lastMode         string
	lastPassband     string
}

type TRX interface {
	SetPowerLevel(value string) error
	GetPowerLevel() (string, error)
	SetMode(mode, passband string) error
	GetMode() (string, string, error)
	SetTx(enabled bool) error
}

type powerLevel = core.PowerLevel

func (a *App) PowerLevels() []string {
	return a.powerLevelLabels
}

func (a *App) SetPowerLevel(index int) error {
	if a.tuning {
		err := a.endTuning()
		if err != nil {
			return err
		}
	}
	return a.trx.SetPowerLevel(a.powerLevels[index].Value)
}

func (a *App) SetTx(enabled bool) error {
	return a.trx.SetTx(enabled)
}

func (a *App) Tuning() bool {
	return a.tuning
}

func (a *App) ToggleTuning() (bool, error) {
	var err error
	if a.tuning {
		err = a.endTuning()
	} else {
		err = a.beginTuning()
	}

	return a.tuning, err
}

func (a *App) beginTuning() error {
	var err error
	if err == nil {
		a.lastValue, err = a.trx.GetPowerLevel()
	}
	if err == nil {
		a.lastMode, a.lastPassband, err = a.trx.GetMode()
	}
	if err == nil {
		err = a.trx.SetPowerLevel(a.tuningValue)
	}
	if err == nil {
		err = a.trx.SetMode("FM", "8000")
	}
	if err == nil {
		err = a.trx.SetTx(true)
	}
	if err != nil {
		return err
	}
	a.tuning = true
	return nil
}

func (a *App) endTuning() error {
	var err error
	if err == nil {
		err = a.trx.SetTx(false)
	}
	if err == nil {
		err = a.trx.SetMode(a.lastMode, a.lastPassband)
	}
	if err == nil {
		err = a.trx.SetPowerLevel(a.lastValue)
	}
	if err != nil {
		return err
	}
	a.tuning = false
	return nil
}
