package main

import (
	"log"

	"github.com/ftl/pwrswtch/core/app"
	"github.com/ftl/pwrswtch/core/trx"
	"github.com/ftl/pwrswtch/ui"
)

func main() {
	tx := trx.Open("")
	app := app.New(tx)

	if err := ui.Run(app); err != nil {
		log.Fatal(err)
	}
}
