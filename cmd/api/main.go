package main

import (
	"github.com/hamzausmani302/prometheus-database-exporter/internal/initiator"
)

func main() {
	done := make(chan bool, 1)
	app := initiator.Application{Done: done}
	if err := app.Init(); err != nil {
		panic(err)
	}
	app.StartApi()

}
