//go:build js

package main

import (
	"log"

	"codeberg.org/anaseto/gruid"
	js "codeberg.org/anaseto/gruid-js"
)

var driver gruid.Driver

func init() {
	t, err := NewTileDrawer()
	if err != nil {
		log.Fatal(err)
	}
	driver = js.NewDriver(js.Config{
		TileManager: t,
	})
}
