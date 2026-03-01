package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/rl"
)

type game struct {
	ECS *ECS
	Map *Map
	Log []LogEntry
}

const (
	MonstersToSpawn = 4
	ScrollsToPlace  = 3
	PotionsToPlace  = 3
	TorchesToPlace  = 5
)

var Directions = []gruid.Point{
	gruid.Point{X: 0, Y: -1},  // N
	gruid.Point{X: 1, Y: 0},   // E
	gruid.Point{X: 0, Y: 1},   // S
	gruid.Point{X: -1, Y: 0},  // W
	gruid.Point{X: 1, Y: -1},  // NE
	gruid.Point{X: 1, Y: 1},   // SE
	gruid.Point{X: -1, Y: -1}, // NW
	gruid.Point{X: -1, Y: 1},  // SW
}

func (g *game) Initialize() {
	// Initialize map and ECS.
	g.Map = NewMap(gruid.Point{X: MapWidth, Y: MapHeight})
	g.ECS = NewECS()
	g.ECS.Map = g.Map
	// Place player on a random floor.
	g.NewPlayer(g.FreeFloorTile())
	// Spawn enemies, place items, and advance a tick.
	g.SpawnEnemies()
	g.SpawnPotions()
	g.SpawnScrolls()
	g.SpawnGrass()
	g.SpawnTorches()
	// g.SpawnCorpses()
	g.ECS.Initialize()
}

func (g *game) Pathable(p gruid.Point) bool {
	if g.Map.Walkable(p) && g.Map.Explored[g.Map.idx(p)] {
		return true
	}
	return false
}

func (g *game) SpawnEnemies() {
	for i := 0; i < MonstersToSpawn; i++ {
		switch {
		case g.Map.Rand.IntN(100) < 80:
			g.NewGoblin(g.FreeFloorTile())
		default:
			g.NewTroll(g.FreeFloorTile())
		}
	}
}

// Places potions and other items throughout the map during gen.
func (g *game) SpawnPotions() {
	for i := 0; i < PotionsToPlace; i++ {
		g.NewHealthPotion(g.FreeFloorTile())
	}
}

func (g *game) SpawnScrolls() {
	for i := 0; i < ScrollsToPlace; i++ {
		g.NewScroll(g.FreeFloorTile())
	}
}

func (g *game) SpawnGrass() {
	// Generate a clump mask using cellular automata on a scratch grid.
	// Cells that come out as GrassFloor define where grass can grow.
	const GrassFloor rl.Cell = 2
	scratch := rl.NewGrid(MapWidth, MapHeight)
	mgen := rl.MapGen{Rand: g.Map.Rand, Grid: scratch}
	rules := []rl.CellularAutomataRule{
		{WCutoff1: 5, WCutoff2: 2, Reps: 3, WallsOutOfRange: true},
	}
	mgen.CellularAutomataCave(Wall, GrassFloor, 0.65, rules)

	it := scratch.Iterator()
	for it.Next() {
		p := it.P()
		if it.Cell() == GrassFloor && g.Map.Grid.At(p) == Floor && g.ECS.NoBlockingEntityAt(p) {
			g.NewGrass(p)
		}
	}
}

func (g *game) SpawnTorches() {
	for i := 0; i < TorchesToPlace; i++ {
		g.NewTorch(g.FreeFloorTile())
	}
}

func (g *game) SpawnCorpses() {
	const corpsesToSpawn = 10
	for i := 0; i < corpsesToSpawn; i++ {
		g.NewCorpse(g.Map.RandomFloor())
	}
}

func (g *game) PickupItem() (ok bool) {
	// Right now only looking at entities that have both input and inventory (player)
	// but want to write way of doing this that doesn't care about input
	ok = false
	p := g.PlayerPosition()
	inv := g.PlayerInventory()
	for _, i := range g.ECS.EntitiesAtPWith(p, Collectible{}) {
		// There is an item here that is collectible! Place a reference to it
		// in e's inventory and remove its Position component.
		ok = true
		item_name := GetComponent[Name](g.ECS, i).string
		g.Logf("You pick up the %s.", ColorLogSpecial, item_name)
		inv.items[inv.nextKey()] = i
		g.ECS.RemoveComponent(i, Position{})
	}
	g.ECS.AddComponent(0, inv)
	return ok
}

// Returns a free floor tile in the map.
func (g *game) FreeFloorTile() gruid.Point {
	for {
		p := g.Map.RandomFloor()
		if g.ECS.NoBlockingEntityAt(p) {
			return p
		}
	}
}

// Return the position of the player.
func (g *game) PlayerPosition() gruid.Point {
	if p, hasPos := g.ECS.GetComponent(0, Position{}); hasPos {
		return p.(Position).Point
	}
	return gruid.Point{}
}

func (g *game) PlayerInventory() Inventory {
	if inv, hasInv := g.ECS.GetComponent(0, Inventory{}); hasInv {
		return inv.(Inventory)
	}
	return Inventory{}
}
