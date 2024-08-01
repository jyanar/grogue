package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type game struct {
	ECS *ECS
	Map *Map
	Log []LogEntry
}

const (
	MonstersToSpawn = 6
	ScrollsToPlace  = 5
	PotionsToPlace  = 5
)

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
	// g.SpawnCorpses()
	// pp := g.ECS.positions[0].Point
	// g.NewExampleAnimation(pp.Add(gruid.Point{X: 1, Y: 0}))
	// g.NewWaterTile(pp.Add(gruid.Point{X: 2, Y: 0}))
	g.ECS.Initialize()
}

// InFOV returns true if p is in the field of view of an entity with FOV. We only
// keep cells within maxLOS manhattan distance from the source entity.
func (g *game) InFOV(p gruid.Point) bool {
	for _, e := range g.ECS.EntitiesWith(Position{}, FOV{}) {
		pp, _ := g.ECS.GetComponent(e, Position{})
		f, _ := g.ECS.GetComponent(e, FOV{})
		pos := pp.(Position)
		fov := f.(FOV)
		if fov.FOV.Visible(p) && paths.DistanceManhattan(pos.Point, p) <= fov.LOS {
			return true
		}
	}
	return false
}

func (g *game) Pathable(p gruid.Point) bool {
	if g.Map.Walkable(p) && g.Map.Explored[p] {
		return true
	}
	return false
}

func (g *game) SpawnEnemies() {
	for i := 0; i < MonstersToSpawn; i++ {
		switch {
		case g.Map.Rand.Intn(100) < 80:
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
	for _, e := range g.ECS.EntitiesWith(Input{}, Inventory{}) {
		// Check if e is standing over a collectible item.
		pos, _ := g.ECS.GetComponent(e, Position{})
		p := pos.(Position).Point
		for _, i := range g.ECS.EntitiesAt(p) {
			if i != e && g.ECS.HasComponent(i, Collectible{}) {
				// There is an item here that is collectible! Place a reference to it
				// in e's inventory and remove its Position component.
				ok = true
				n, _ := g.ECS.GetComponent(e, Name{})
				in, _ := g.ECS.GetComponent(i, Name{})
				entity_name := n.(Name).string
				item_name := in.(Name).string
				if e == 0 {
					g.Logf("You pick up the %s.", ColorLogSpecial, item_name)
				} else {
					g.Logf("%s picks up %s.", ColorLogSpecial, entity_name, item_name)
				}

				inventory, _ := g.ECS.GetComponent(e, Inventory{})
				inv := inventory.(Inventory)
				inv.items = append(inv.items, i)
				g.ECS.AddComponent(e, inv)
				g.ECS.RemoveComponent(i, Position{})
				// g.ECS.inventories[e].items = append(g.ECS.inventories[e].items, i)
				// g.ECS.positions[i] = nil
			}
		}
	}
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
