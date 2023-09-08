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

// InFOV returns true if p is in the field of view of an entity with FOV. We only
// keep cells within maxLOS manhattan distance from the source entity.
//
// NOTE: Currently InFOV only returns true for the player FOV.
func (g *game) InFOV(p gruid.Point) bool {
	pp := g.ECS.positions[0].Point
	los := g.ECS.fovs[0].LOS
	if g.ECS.fovs[0].FOV.Visible(p) && paths.DistanceManhattan(pp, p) <= los {
		return true
	}
	return false
}

const MonstersToSpawn = 6

func (g *game) SpawnEnemies() {
	for i := 0; i < MonstersToSpawn; i++ {
		switch {
		case g.Map.Rand.Intn(100) < 80:
			g.NewGoblin()
		default:
			g.NewTroll()
		}
	}
}

const PotionsToPlace = 5

// Places potions and other items throughout the map during gen.
func (g *game) PlaceItems() {
	for i := 0; i < PotionsToPlace; i++ {
		g.ECS.Create(
			Name{"Health Potion"},
			Renderable{glyph: '!', color: ColorHealthPotion, order: ROItem},
			Collectible{},
			Consumable{hp: 5},
			Position{g.FreeFloorTile()},
		)
	}
}

func (g *game) PickupItem() (ok bool) {
	// Right now only looking at entities that have both input and inventory (player)
	// but want to write way of doing this that doesn't care about input
	ok = false
	for _, e := range g.ECS.EntitiesWith(Input{}, Inventory{}) {
		// Check if e is standing over a collectible item.
		p := g.ECS.positions[e].Point
		for _, i := range g.ECS.EntitiesAt(p) {
			if i != e && g.ECS.HasComponent(i, Collectible{}) {
				// There is an item here that is collectible! Place a reference to it
				// in e's inventory and remove both its Position and Renderable components.
				ok = true
				name := g.ECS.names[e].string
				itemName := g.ECS.names[i].string
				g.Logf("%s picks up %s.", ColorLogSpecial, name, itemName)
				g.ECS.inventories[e].items = append(g.ECS.inventories[e].items, i)
				g.ECS.positions[i] = nil
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
