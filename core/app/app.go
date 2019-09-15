package app

import (
	"github.com/ftl/pwrswtch/core"
)

func New(trx TRX) *App {
	return &App{
		trx: trx,
		powerLevels: []core.PowerLevel{
			{10, "0.039216"},
			{30, "0.117647"},
			{50, "0.196078"},
			{100, "1.0"},
		},
	}
}

type App struct {
	trx         TRX
	powerLevels []core.PowerLevel
}

type TRX interface {
	SetPowerLevel(value string) error
	SetTx(enabled bool) error
}

func (a *App) PowerLevels() []core.PowerLevel {
	return a.powerLevels
}

func (a *App) SetPowerLevel(value string) error {
	return a.trx.SetPowerLevel(value)
}

func (a *App) SetTx(enabled bool) error {
	return a.trx.SetTx(enabled)
}
