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
