package main

import (
	"log"
	"todo/app"
	"todo/config"
)

func main() {
	cfg := config.LoadConfig()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Init new app have error is - %s\n", err)
	}

	if err = app.RunApp(); err != nil {
		log.Fatalf("Started application RunAPI have error is - %s\n", err)
	}
}
