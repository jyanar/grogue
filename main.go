package main

import (
	"context"
	"log"

	"github.com/anaseto/gruid"
)

const (
	UIWidth   = 60
	UIHeight  = 27
	MapWidth  = UIWidth - 2
	MapHeight = UIHeight - 5
	LogLines  = 5
)

func main() {
	// Construct the drawgrid, and a new model.
	gd := gruid.NewGrid(UIWidth, UIHeight)
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
