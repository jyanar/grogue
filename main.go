package main

import (
	"context"
	"log"

	"github.com/anaseto/gruid"
)

func main() {
	// Construct the drawgrid, and a new model.
	gd := gruid.NewGrid(80, 24)
	m := NewModel(gd)

	// Instantiate new app. driver is generated in sdl.go, or in
	// js.go if application is built with js flags (see README).
	app := gruid.NewApp(gruid.AppConfig{
		Model:  m,
		Driver: driver,
	})

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
