package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

const (
	Wall rl.Cell = iota
	Floor
)

// Map represents the rectangular grid of the game's level.
type Map struct {
	Grid rl.Grid
}

func NewMap(size gruid.Point) *Map {
	m := &Map{}
	m.Grid = rl.NewGrid(size.X, size.Y)
	m.Grid.Fill(Floor)
	for i := 0; i < 3; i++ {
		// We add a few extra walls. We'll deal with map generation
		// in the next part of the tutorial.
		m.Grid.Set(gruid.Point{X: 30 + i, Y: 12}, Wall)
	}
	return m
}

func (m *Map) Walkable(p gruid.Point) bool {
	return m.Grid.At(p) == Floor
}

func (m *Map) Rune(c rl.Cell) (r rune) {
	switch c {
	case Wall:
		r = '#'
	case Floor:
		r = '.'
	}
	return r
}
