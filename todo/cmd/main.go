package main

import (
	"todo/app"
	"todo/config"
)

func main() {
	cfg := config.NewFromEnv()

	app, err := app.NewApp(cfg)
	if err != nil {
		panic(err)
	}

	app.RunAPI()
}
