package main

import "github.com/anaseto/gruid"

type Bar struct {
	Position gruid.Point // Point on grid at which to draw Bar
	Style    gruid.Style // Bar style
	Width    int         // Bar width
}

func (b *Bar) Draw(gd gruid.Grid) {
	for i := 0; i < b.Width; i++ {
		gd.Set(b.Position.Shift(i, 0), gruid.Cell{Rune: ' ', Style: b.Style})
	}
}

type HpBar struct {
	top  Bar
	back Bar
}
