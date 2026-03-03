package main

import (
	"math"
	"math/rand/v2"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
)

// RoomGen generates room shapes as rl.Grid values (Wall/Floor cells only).
// Every returned grid includes a 1-cell wall border on all sides so that rooms
// placed on the main map never expose raw edge cells to the map boundary.
type RoomGen struct {
	Rand *rand.Rand
}

func (rg *RoomGen) RectRoom() rl.Grid {
	w := 3 + rg.Rand.IntN(8)
	h := 3 + rg.Rand.IntN(4)
	grid := rl.NewGrid(w+2, h+2)
	fillRect(grid, 1, 1, w, h)
	return grid
}

// RectsRoom returns a room made of two overlapping rectangles. The two
// rectangles are guaranteed to overlap by at least 2 cells on each axis so
// the room is always a single connected region.
func (rg *RoomGen) RectsRoom() rl.Grid {
	w1 := 3 + rg.Rand.IntN(8) // 3–10
	h1 := 3 + rg.Rand.IntN(4) // 3–6
	w2 := 3 + rg.Rand.IntN(8)
	h2 := 3 + rg.Rand.IntN(4)

	// Choose an offset for rect2 relative to rect1 that guarantees ≥2 cells
	// of overlap on each axis. For width: offset in [-(w2-2), w1-2].
	ox := -(w2 - 2) + rg.Rand.IntN(w1+w2-4)
	oy := -(h2 - 2) + rg.Rand.IntN(h1+h2-4)

	// Bounding box of the union of both rectangles.
	minX, minY := min(0, ox), min(0, oy)
	maxX, maxY := max(w1, ox+w2), max(h1, oy+h2)

	// Allocate grid with a 1-cell wall border (+2 on each dimension).
	grid := rl.NewGrid(maxX-minX+2, maxY-minY+2)

	// The border offset shifts both rects 1 cell in from the grid edge.
	fillRect(grid, 1-minX, 1-minY, w1, h1)
	fillRect(grid, 1+ox-minX, 1+oy-minY, w2, h2)

	return grid
}

// BlobRoom returns an organic, cave-like room produced by cellular automata.
// The result is non-deterministically shaped but generally compact.
func (rg *RoomGen) BlobRoom() rl.Grid {
	const w, h = 18, 11
	grid := rl.NewGrid(w, h)
	mgen := rl.MapGen{Rand: rg.Rand, Grid: grid}
	// Cells with 5+ wall neighbours become wall; isolated walls with ≤2 wall
	// neighbours open up. WallsOutOfRange ensures the border stays solid,
	// giving a natural 1-cell wall edge without extra bookkeeping.
	mgen.CellularAutomataCave(Wall, Floor, 0.50, []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 10, WallsOutOfRange: true},
	})
	return grid
}

// CircleRoom returns a single large circular room. The +0.5 on the radius
// threshold rounds the outermost ring inward rather than leaving a
// cross-shaped notch at each cardinal direction.
func (rg *RoomGen) CircleRoom() rl.Grid {
	radius := 3 + rg.Rand.IntN(3) // 3–5
	// diameter + 1-cell border on each side
	gw := 2*radius + 3
	gh := 2*radius + 3
	grid := rl.NewGrid(gw, gh)
	cx, cy := radius+1, radius+1 // centre, shifted in by the border
	for y := range gh {
		for x := range gw {
			dx, dy := x-cx, y-cy
			if math.Sqrt(float64(dx*dx+dy*dy)) <= float64(radius)+0.5 {
				grid.Set(gruid.Point{X: x, Y: y}, Floor)
			}
		}
	}
	return grid
}

// CirclesRoom returns a room formed by 2–4 overlapping small circles placed
// randomly within a fixed bounding area.
func (rg *RoomGen) CirclesRoom() rl.Grid {
	const areaW, areaH = 14, 9
	gw, gh := areaW+2, areaH+2
	grid := rl.NewGrid(gw, gh)
	n := 2 + rg.Rand.IntN(3) // 2–4 circles
	for range n {
		radius := 2 + rg.Rand.IntN(3) // 2–4
		// Clamp centre so the full circle fits within the bordered area.
		cx := 1 + radius + rg.Rand.IntN(max(1, areaW-2*radius))
		cy := 1 + radius + rg.Rand.IntN(max(1, areaH-2*radius))
		for y := range gh {
			for x := range gw {
				dx, dy := x-cx, y-cy
				if math.Sqrt(float64(dx*dx+dy*dy)) <= float64(radius)+0.5 {
					grid.Set(gruid.Point{X: x, Y: y}, Floor)
				}
			}
		}
	}
	return grid
}

// Random picks one of the four room types at random and returns it.
// The grid is pruned to its largest connected floor component before
// being returned, ensuring no disconnected islands remain.
func (rg *RoomGen) Random() rl.Grid {
	var g rl.Grid
	switch rg.Rand.IntN(2) {
	// case 0:
	// 	g = rg.RectRoom()
	// case 1:
	// 	g = rg.RectRoom()
	// case 2:
	// 	g = rg.CircleRoom()
	default:
		g = rg.RectRoom()
		// g = rg.BlobRoom()
	}
	pruneRoom(g)
	return g
}

// pruneRoom removes all floor cells not reachable from the largest connected
// floor component, converting them back to Wall. This prevents disconnected
// islands from making parts of the dungeon permanently inaccessible.
func pruneRoom(grid rl.Grid) {
	size := grid.Size()
	rng := grid.Range()
	width := size.X

	cellIdx := func(p gruid.Point) int { return p.Y*width + p.X }

	// Find all connected floor components via BFS.
	label := make([]int, size.X*size.Y)
	for i := range label {
		label[i] = -1
	}
	type component struct{ cells []gruid.Point }
	var comps []component

	it := grid.Iterator()
	for it.Next() {
		p := it.P()
		if it.Cell() != Floor || label[cellIdx(p)] >= 0 {
			continue
		}
		id := len(comps)
		comps = append(comps, component{})
		queue := []gruid.Point{p}
		label[cellIdx(p)] = id
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			comps[id].cells = append(comps[id].cells, cur)
			for _, d := range cardinals {
				nb := gruid.Point{X: cur.X + d.X, Y: cur.Y + d.Y}
				if !nb.In(rng) || label[cellIdx(nb)] >= 0 || grid.At(nb) != Floor {
					continue
				}
				label[cellIdx(nb)] = id
				queue = append(queue, nb)
			}
		}
	}

	if len(comps) <= 1 {
		return
	}

	// Identify the largest component.
	best := 0
	for i, c := range comps {
		if len(c.cells) > len(comps[best].cells) {
			best = i
		}
	}

	// Convert all other components back to Wall.
	for id, c := range comps {
		if id == best {
			continue
		}
		for _, p := range c.cells {
			grid.Set(p, Wall)
		}
	}
}

// fillRect sets all cells in a w×h rectangle starting at (x, y) to Floor.
func fillRect(grid rl.Grid, x, y, w, h int) {
	for dy := range h {
		for dx := range w {
			grid.Set(gruid.Point{X: x + dx, Y: y + dy}, Floor)
		}
	}
}

var cardinals = [4]gruid.Point{
	{X: 0, Y: -1}, // N
	{X: 1, Y: 0},  // E
	{X: 0, Y: 1},  // S
	{X: -1, Y: 0}, // W
}

// edgeFloors returns all floor cells in room whose immediate neighbour in
// direction dir is a wall (or out of bounds). These are the candidate
// attachment points for hallways and entrances.
func edgeFloors(room rl.Grid, dir gruid.Point) []gruid.Point {
	rng := room.Range()
	var pts []gruid.Point
	it := room.Iterator()
	for it.Next() {
		if it.Cell() != Floor {
			continue
		}
		p := it.P()
		next := gruid.Point{X: p.X + dir.X, Y: p.Y + dir.Y}
		if !next.In(rng) || room.At(next) == Wall {
			pts = append(pts, p)
		}
	}
	return pts
}

// Entrance records a passage carved through the room's border wall.
// Pos is the entrance cell in room-local coordinates (just outside the
// interior floor, in the 1-cell border). Dir is the outward direction
// (points toward the existing dungeon when the room is placed). Hall
// contains every cell to carve—for a bare entrance it is just {Pos}; for
// a hallway it lists every corridor cell from the border outward, with
// Pos as the far end.
type Entrance struct {
	Pos  gruid.Point
	Dir  gruid.Point
	Hall []gruid.Point
}

// RoomInstance pairs a room shape with its computed entrance metadata.
type RoomInstance struct {
	Grid      rl.Grid
	Entrances []Entrance
}

// Instance generates a random room and assigns entrance(s). With 50%
// probability the room gets a single hallway entrance; otherwise it gets
// one bare entrance per cardinal side.
func (rg *RoomGen) Instance() RoomInstance {
	room := rg.Random()
	if rg.Rand.IntN(2) == 0 {
		return rg.withHallway(room)
	}
	return rg.withEntrances(room)
}

// withEntrances picks one random edge floor cell per cardinal side and
// records the adjacent border wall cell as a bare entrance.
func (rg *RoomGen) withEntrances(room rl.Grid) RoomInstance {
	var entrances []Entrance
	for _, dir := range cardinals {
		cands := edgeFloors(room, dir)
		if len(cands) == 0 {
			continue
		}
		p := cands[rg.Rand.IntN(len(cands))]
		ePos := gruid.Point{X: p.X + dir.X, Y: p.Y + dir.Y}
		entrances = append(entrances, Entrance{
			Pos:  ePos,
			Dir:  dir,
			Hall: []gruid.Point{ePos},
		})
	}
	return RoomInstance{Grid: room, Entrances: entrances}
}

// withHallway picks a random side, selects a random edge floor cell on
// that side, and records a hallway of length 2–5 extending outward. The
// far end of the hallway is the entrance.
func (rg *RoomGen) withHallway(room rl.Grid) RoomInstance {
	start := rg.Rand.IntN(4)
	for i := range 4 {
		dir := cardinals[(start+i)%4]
		cands := edgeFloors(room, dir)
		if len(cands) == 0 {
			continue
		}
		p := cands[rg.Rand.IntN(len(cands))]
		length := 2 + rg.Rand.IntN(6) // 2–8
		var hall []gruid.Point
		for step := 1; step <= length; step++ {
			hall = append(hall, gruid.Point{X: p.X + dir.X*step, Y: p.Y + dir.Y*step})
		}
		return RoomInstance{
			Grid:      room,
			Entrances: []Entrance{{Pos: hall[len(hall)-1], Dir: dir, Hall: hall}},
		}
	}
	// Fallback: no valid side — use bare entrances.
	return rg.withEntrances(room)
}
