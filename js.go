//go:build js
// +build js

package main

import (
	"log"

	"github.com/anaseto/gruid"
	js "github.com/anaseto/gruid-js"
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
