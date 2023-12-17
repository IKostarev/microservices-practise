package main

import (
	"notifications/app"
	"notifications/config"
)

func main() {
	cfg := config.NewFromEnv()

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}

	err = app.RunAPI()
	if err != nil {
		panic(err)
	}
}
