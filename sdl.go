//go:build !js

package main

import (
	"log"

	"codeberg.org/anaseto/gruid"
	sdl "codeberg.org/anaseto/gruid-sdl"
)

var driver gruid.Driver

func init() {
	t, err := NewTileDrawer()
	if err != nil {
		log.Fatal(err)
	}
	dr := sdl.NewDriver(sdl.Config{
		WindowTitle: "grogue",
		TileManager: t,
	})
	dr.PreventQuit()
	driver = dr
}
