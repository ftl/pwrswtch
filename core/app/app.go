package app

import (
	"fmt"
)

func New(trx TRX) *App {
	result := App{
		trx: trx,
		powerLevels: []powerLevel{
			{10, "0.039216"},
			{30, "0.117647"},
			{50, "0.196078"},
			{100, "1.0"},
		},
		tuningValue: "0.039216",
	}

	result.powerLevelLabels = make([]string, len(result.powerLevels))
	for i, level := range result.powerLevels {
		result.powerLevelLabels[i] = fmt.Sprintf("%dW", level.watts)
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
}

type TRX interface {
	SetPowerLevel(value string) error
	GetPowerLevel() (string, error)
	SetTx(enabled bool) error
}

type powerLevel struct {
	watts int
	value string
}

func (a *App) PowerLevels() []string {
	return a.powerLevelLabels
}

func (a *App) SetPowerLevel(index int) error {
	return a.trx.SetPowerLevel(a.powerLevels[index].value)
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
	a.lastValue, err = a.trx.GetPowerLevel()
	if err != nil {
		return err
	}
	err = a.trx.SetPowerLevel(a.tuningValue)
	if err != nil {
		return err
	}
	a.tuning = true
	return nil
}

func (a *App) endTuning() error {
	err := a.trx.SetPowerLevel(a.lastValue)
	if err != nil {
		return err
	}
	a.tuning = false
	return nil
}
