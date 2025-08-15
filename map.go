package main

import (
	"math/rand/v2"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"codeberg.org/anaseto/gruid/rl"
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
	PR       *paths.PathRange
}

func NewMap(size gruid.Point) *Map {
	m := &Map{
		Grid:     rl.NewGrid(size.X, size.Y),
		Rand:     rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		Explored: make(map[gruid.Point]bool),
		PR:       paths.NewPathRange(gruid.NewRange(0, 0, size.X, size.Y)),
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

// Generate fills the Grid attribute of m with a procedurally generated map.
func (m *Map) Generate() {
	// Using the rl package from gruid.
	mgen := rl.MapGen{Rand: m.Rand, Grid: m.Grid}
	// Cellular automata map generation with rules that give a cave-like map.
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 4, WallsOutOfRange: true},
		{WCutoff1: 5, WCutoff2: 25, Reps: 3, WallsOutOfRange: true},
	}
	for {
		mgen.CellularAutomataCave(Wall, Floor, 0.45, rules)
		freep := m.RandomFloor()
		// We put walls in the floor cells non-reachable from freep, to ensure
		// that all the cells are connected (which is not guaranteed with CA).
		pr := paths.NewPathRange(m.Grid.Range())
		pr.CCMap(&path{m: m}, freep)
		ntiles := mgen.KeepCC(pr, freep, Wall)
		const minCaveSize = 400
		if ntiles > minCaveSize {
			break
		}
		// If there were not enough free tiles, we run the map generation again.
	}
}

// RandomFloor returns a random floor cell in the map. It assumes that such a
// floor cel exists (otherwise the function does not end).
func (m *Map) RandomFloor() gruid.Point {
	size := m.Grid.Size()
	for {
		freep := gruid.Point{X: m.Rand.IntN(size.X), Y: m.Rand.IntN(size.Y)}
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
