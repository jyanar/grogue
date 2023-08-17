package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type game struct {
	ECS *ECS
	Map *Map
}

// InFOV returns true if p is in the field of view of an entity with FOV. We only
// keep cells within maxLOS manhattan distance from the source entity.
//
// NOTE: Currently InFOV only returns true for the player FOV.
func (g *game) InFOV(p gruid.Point) bool {
	pp := g.ECS.positions[0].Point
	if g.ECS.fovs[0].FOV.Visible(p) && paths.DistanceManhattan(pp, p) <= 10 {
		return true
	}
	return false
}

const MonstersToSpawn = 6

func (g *game) SpawnEnemies() {
	for i := 0; i < MonstersToSpawn; i++ {
		switch {
		case g.Map.Rand.Intn(100) < 80:
			g.ECS.Create(
				Position{g.Map.RandomFloor()},
				Name{"Goblin"},
				Renderable{glyph: 'g', color: ColorMonster, order: ROActor},
				Health{hp: 10, maxhp: 10},
				Damage{2},
				Perception{radius: 4},
				AI{state: CSWandering},
				Obstruct{},
			)
		default:
			g.ECS.Create(
				Position{g.Map.RandomFloor()},
				Name{"Orc"},
				Renderable{glyph: 'o', color: ColorMonster, order: ROActor},
				Health{hp: 15, maxhp: 15},
				Damage{3},
				Perception{radius: 4},
				AI{state: CSWandering},
				Obstruct{},
			)
		}
		g.ECS.Create()
	}
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
