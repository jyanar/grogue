package main

import (
	"context"
	"log"

	"github.com/anaseto/gruid"
)

const (
	UIWidth   = 80
	UIHeight  = 24
	MapWidth  = UIWidth
	MapHeight = UIHeight - 3
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
