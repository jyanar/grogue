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
	Grid       rl.Grid    // Gamemap.
	Rand       *rand.Rand // Random number generator.
	Explored   []bool     // Flat array [y*MapWidth+x]: tiles the player has ever seen.
	LightMap   []float32  // Flat array [y*MapWidth+x]: per-tile light level (0.0–1.0), updated each turn.
	VisibleNow []bool     // Flat array [y*MapWidth+x]: tiles in player FOV this turn, updated each turn.
	PR         *paths.PathRange
}

// idx converts a map point to a flat array index.
func (m *Map) idx(p gruid.Point) int {
	return p.Y*MapWidth + p.X
}

func NewMap(size gruid.Point) *Map {
	n := size.X * size.Y
	m := &Map{
		Grid:       rl.NewGrid(size.X, size.Y),
		Rand:       rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		Explored:   make([]bool, n),
		LightMap:   make([]float32, n),
		VisibleNow: make([]bool, n),
		PR:         paths.NewPathRange(gruid.NewRange(0, 0, size.X, size.Y)),
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

// pathDist returns the walking distance from src to dst using BFS, capped at
// limit. Returns limit+1 if dst is not reached within that many steps.
func (m *Map) pathDist(src, dst gruid.Point, limit int) int {
	p := m.PR.JPSPath(nil, src, dst, m.Walkable, false)
	if len(p) == 0 {
		return limit + 1
	}
	return len(p) - 1
}

// ConnectRooms scans every wall tile that has walkable tiles on both sides
// (N/S or E/W). If the walking distance between those two tiles exceeds the
// threshold, the wall is replaced with a Floor.
func (m *Map) ConnectRooms() {
	const threshold = 30
	mr := m.Grid.Range()
	it := m.Grid.Iterator()
	for it.Next() {
		p := it.P()
		if it.Cell() != Wall {
			continue
		}
		var a, b gruid.Point
		west := gruid.Point{X: p.X - 1, Y: p.Y}
		east := gruid.Point{X: p.X + 1, Y: p.Y}
		north := gruid.Point{X: p.X, Y: p.Y - 1}
		south := gruid.Point{X: p.X, Y: p.Y + 1}
		switch {
		case west.In(mr) && east.In(mr) && m.Walkable(west) && m.Walkable(east):
			a, b = west, east
		case north.In(mr) && south.In(mr) && m.Walkable(north) && m.Walkable(south):
			a, b = north, south
		default:
			continue
		}
		if m.pathDist(a, b, threshold) > threshold {
			m.Grid.Set(p, Floor)
		}
	}
}

// canPlace reports whether ri can be stamped at origin using entrance eIdx.
// Every room cell (floor and wall border) and every hallway cell must lie
// within map bounds and currently be Wall. Hallway side cells (perpendicular
// to the hall direction) are also checked to prevent hallways from running
// flush against existing rooms and merging them.
func (m *Map) canPlace(ri RoomInstance, origin gruid.Point, eIdx int) bool {
	mr := m.Grid.Range()
	it := ri.Grid.Iterator()
	for it.Next() {
		mp := gruid.Point{X: origin.X + it.P().X, Y: origin.Y + it.P().Y}
		if !mp.In(mr) || m.Grid.At(mp) != Wall {
			return false
		}
	}
	e := ri.Entrances[eIdx]
	// Perpendicular unit vectors to the hall direction.
	side1 := gruid.Point{X: -e.Dir.Y, Y: e.Dir.X}
	side2 := gruid.Point{X: e.Dir.Y, Y: -e.Dir.X}
	for _, hc := range e.Hall {
		mp := gruid.Point{X: origin.X + hc.X, Y: origin.Y + hc.Y}
		if !mp.In(mr) || m.Grid.At(mp) != Wall {
			return false
		}
		for _, s := range [2]gruid.Point{side1, side2} {
			sp := gruid.Point{X: mp.X + s.X, Y: mp.Y + s.Y}
			if sp.In(mr) && m.Grid.At(sp) != Wall {
				return false
			}
		}
	}
	return true
}

// stampRoom carves ri's floor cells and entrance/hallway cells onto the map.
func (m *Map) stampRoom(ri RoomInstance, origin gruid.Point, eIdx int) {
	mr := m.Grid.Range()
	it := ri.Grid.Iterator()
	for it.Next() {
		if it.Cell() != Floor {
			continue
		}
		m.Grid.Set(gruid.Point{X: origin.X + it.P().X, Y: origin.Y + it.P().Y}, Floor)
	}
	for _, hc := range ri.Entrances[eIdx].Hall {
		mp := gruid.Point{X: origin.X + hc.X, Y: origin.Y + hc.Y}
		if mp.In(mr) {
			m.Grid.Set(mp, Floor)
		}
	}
}

// tryPlaceRoom tries up to 500 random map positions to find a floor cell F
// such that placing ri with one of its entrances adjacent to F is valid.
// Placement formula: origin = F − entrance.Pos − entrance.Dir.
// Returns true if the room was placed.
func (m *Map) tryPlaceRoom(ri RoomInstance) bool {
	if len(ri.Entrances) == 0 {
		return false
	}
	// log.Println("Attempting to place room...")
	mapSize := m.Grid.Size()
	for range 500 {
		p := gruid.Point{X: m.Rand.IntN(mapSize.X), Y: m.Rand.IntN(mapSize.Y)}
		if m.Grid.At(p) != Floor {
			continue
		}
		for eIdx, e := range ri.Entrances {
			origin := gruid.Point{
				X: p.X - e.Pos.X - e.Dir.X,
				Y: p.Y - e.Pos.Y - e.Dir.Y,
			}
			if m.canPlace(ri, origin, eIdx) {
				m.stampRoom(ri, origin, eIdx)
				return true
			}
		}
	}
	// log.Println("Room placement failed.")
	// ri.Print()
	return false
}

// Generate fills the map using Brogue's iterative room-placement algorithm.
// A first room is stamped at the center; then rooms are added one by one,
// each connecting to any existing floor cell via its entrance, until 50
// consecutive placement attempts fail.
func (m *Map) Generate() {
	rg := RoomGen{Rand: m.Rand}

	// First room: centered, no entrance needed.
	first := rg.Random()
	size := first.Size()
	mapSize := m.Grid.Size()
	ox := (mapSize.X - size.X) / 2
	oy := (mapSize.Y - size.Y) / 2
	it := first.Iterator()
	for it.Next() {
		if it.Cell() == Floor {
			m.Grid.Set(gruid.Point{X: ox + it.P().X, Y: oy + it.P().Y}, Floor)
		}
	}

	// Iteratively add rooms until placement consistently fails.
	const maxConsecutiveFailures = 50
	failures := 0
	for failures < maxConsecutiveFailures {
		if m.tryPlaceRoom(rg.Instance()) {
			failures = 0
		} else {
			failures++
		}
	}

	// Connect nearby rooms
	m.ConnectRooms()
}

// PlaceNextRoom attempts to add one more room to the map.
func (m *Map) PlaceNextRoom() {
	rg := RoomGen{Rand: m.Rand}
	m.tryPlaceRoom(rg.Instance())
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
