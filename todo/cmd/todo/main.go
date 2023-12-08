package main

import (
	"log"
	"todo/app"
	"todo/config"
)

func main() {
	cfg := config.LoadConfig()

	app := app.NewApp(cfg)
	if err := app.RunAPI(); err != nil {
		log.Fatalf("Started application RunAPI have error is - %s\n", err)
	}
}
