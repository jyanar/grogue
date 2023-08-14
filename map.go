package main

import (
	"math/rand"
	"time"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

const (
	Wall rl.Cell = iota
	Floor
)

// Map represents the rectangular grid of the game's level.
type Map struct {
	Grid     rl.Grid              // Gamemap.
	Rand     *rand.Rand           // Random number generator.
	Explored map[gruid.Point]bool // Explored tiles.
}

func NewMap(size gruid.Point) *Map {
	m := &Map{
		Grid:     rl.NewGrid(size.X, size.Y),
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		Explored: make(map[gruid.Point]bool),
	}
	m.Generate()
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

// Draws the Map onto the passed grid.
func (m *Map) Draw(grid *gruid.Grid) {
	it := m.Grid.Iterator()
	for it.Next() {
		if !m.Explored[it.P()] {
			continue
		}
		c := gruid.Cell{Rune: m.Rune(it.Cell())}
		grid.Set(it.P(), c)
	}
}

// Generate fills the Grid attribute of m with a procedurally generated map.
func (m *Map) Generate() {
	// Using the rl package from gruid.
	mgen := rl.MapGen{Rand: m.Rand, Grid: m.Grid}
	// Cellular automata map generation with rules that give a cave-like map.
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 4, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, Reps: 3, WallsOutOfRange: true},
	}
	mgen.CellularAutomataCave(Wall, Floor, 0.42, rules)
	freep := m.RandomFloor()
	// We put walls in floor cells nonreachable from freep, to ensure that
	// all the cells are connected (which is not guaranteed by CA map gen).
	pr := paths.NewPathRange(m.Grid.Range())
	pr.CCMap(&path{m: m}, freep)
	mgen.KeepCC(pr, freep, Wall)
}

// RandomFloor returns a random floor cell in the map. It assumes that such a
// floor cel exists (otherwise the function does not end).
func (m *Map) RandomFloor() gruid.Point {
	size := m.Grid.Size()
	for {
		freep := gruid.Point{m.Rand.Intn(size.X), m.Rand.Intn(size.Y)}
		if m.Grid.At(freep) == Floor {
			return freep
		}
	}
}

type path struct {
	m  *Map
	nb paths.Neighbors
}

func (p *path) Neighbors(q gruid.Point) []gruid.Point {
	return p.nb.Cardinal(q,
		func(r gruid.Point) bool { return p.m.Walkable(r) })
}
