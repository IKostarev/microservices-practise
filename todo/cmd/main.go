package main

import (
	"todo/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		panic(err)
	}

	app.RunAPI()
}
