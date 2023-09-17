//go:build !js

package main

import (
	"log"

	"github.com/anaseto/gruid"
	sdl "github.com/anaseto/gruid-sdl"
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
