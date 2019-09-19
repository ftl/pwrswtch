package cfg

import (
	"strconv"

	"github.com/ftl/hamradio/cfg"

	"github.com/ftl/pwrswtch/core"
)

const (
	levels      cfg.Key = "pwrswtch.levels"
	tuningValue cfg.Key = "pwrswtch.tuningValue"
	trxHost     cfg.Key = "pwrswtch.trxHost"
)

func Load() (core.Configuration, error) {
	configuration, err := cfg.LoadDefault()
	if err != nil {
		return core.Configuration{}, err
	}

	result := core.Configuration{
		Levels: []core.PowerLevel{
			{10, "0.039216"},
			{30, "0.117647"},
			{50, "0.196078"},
			{100, "1.0"},
		},
		TuningValue: configuration.Get(tuningValue, "0.039216").(string),
		TRXHost:     configuration.Get(trxHost, "localhost:4532").(string),
	}

	ls := []core.PowerLevel{}
	configuration.GetSlice(levels, func(_ int, e map[string]interface{}) {
		rawWatts, ok := e["watts"]
		if !ok {
			return
		}
		watts, err := strconv.Atoi(rawWatts.(string))
		if err != nil {
			return
		}

		value, ok := e["value"]
		if !ok {
			return
		}

		level := core.PowerLevel{
			Watts: watts,
			Value: value.(string),
		}
		ls = append(ls, level)
	})
	if len(ls) > 0 {
		result.Levels = ls
	}

	return result, nil
}
