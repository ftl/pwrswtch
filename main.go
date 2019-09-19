package main

import (
	"log"

	"github.com/ftl/pwrswtch/core/app"
	"github.com/ftl/pwrswtch/core/cfg"
	"github.com/ftl/pwrswtch/core/trx"
	"github.com/ftl/pwrswtch/ui"
)

func main() {
	config, err := cfg.Load()
	if err != nil {
		log.Fatal(err)
	}

	tx := trx.Open(config.TRXHost)
	app := app.New(tx, config)

	if err := ui.Run(app); err != nil {
		log.Fatal(err)
	}
}
