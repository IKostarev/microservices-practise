package main

import (
	"users/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	app.RunApp()
}
