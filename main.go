package main

import (
	"context"
	"log"

	"github.com/anaseto/gruid"
	sdl "github.com/anaseto/gruid-sdl"
)

func main() {
	// Construct grid.
	gd := gruid.NewGrid(80, 24)

	// Construct a model.
	m := NewModel(gd)

	// Construct a TileDrawer.
	t, err := NewTileDrawer()
	if err != nil {
		log.Fatal(err)
	}

	// Fetch SDL driver (can replace with terminal or JS driver).
	dr := sdl.NewDriver(sdl.Config{
		TileManager: t,
	})

	// Define new application using the SDL2 gruid driver and our model.
	app := gruid.NewApp(gruid.AppConfig{
		Model:  m,
		Driver: dr,
	})

	if err := app.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
