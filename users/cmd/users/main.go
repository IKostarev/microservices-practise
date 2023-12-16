package main

import (
	"users/app"
	"users/config"
)

func main() {
	cfg := config.NewFromEnv()

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}

	err = app.RunApp()
	if err != nil {
		panic(err)
	}
}
