// This file is a set of factories which allow easy spawning of
// different kinds of enemies.

package main

import "github.com/anaseto/gruid"

func (g *game) NewPlayer() int {
	return g.ECS.Create(
		Name{"player"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: '@', fg: ColorPlayer, order: ROActor},
		Health{hp: 18, maxhp: 18},
		Damage{5},
		FOV{LOS: 20},
		Inventory{},
		Input{},
		Obstruct{},
	)
}

func (g *game) NewGoblin() int {
	return g.ECS.Create(
		Name{"goblin"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: 'g', fg: ColorMonster, order: ROActor},
		Health{hp: 10, maxhp: 10},
		Damage{2},
		Perception{LOS: 8},
		AI{state: CSWandering},
		Obstruct{},
	)
}

func (g *game) NewTroll() int {
	return g.ECS.Create(
		Name{"troll"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: 'T', fg: ColorTroll, order: ROActor},
		Health{hp: 20, maxhp: 20},
		Damage{5},
		Perception{LOS: 6},
		AI{state: CSWandering},
		Obstruct{},
	)
}

func (g *game) NewHealthPotion() int {
	return g.ECS.Create(
		Name{"health potion"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: '!', fg: ColorHealthPotion, order: ROItem},
		Collectible{},
		Consumable{hp: 5},
	)
}

// Useful for debugging that corpses pathing bug
// 2023 Sep 18 observed a corpse, goblin, player, and troll standing on the same tile.
func (g *game) NewCorpse() int {
	return g.ECS.Create(
		Name{"corpse"},
		Position{g.Map.RandomFloor()},
		Renderable{glyph: '%', fg: ColorCorpse, order: ROCorpse},
		Collectible{},
		Consumable{hp: 2},
	)
}

func (g *game) NewBlood(p gruid.Point) int {
	return g.ECS.Create(
		Name{"blood"},
		Position{p},
		Renderable{glyph: '.', fg: ColorBlood, bg: ColorBlood, order: ROFloor},
	)
}
