// This file is a set of factories which allow easy spawning of
// different kinds of enemies.

package main

import "github.com/anaseto/gruid"

func (g *game) NewPlayer(p gruid.Point) int {
	return g.ECS.Create(
		Name{"player"},
		Position{p},
		Renderable{glyph: '@', fg: ColorPlayer, order: ROActor},
		Health{hp: 18, maxhp: 18},
		Damage{5},
		FOV{LOS: 20},
		Inventory{},
		Input{},
		Obstruct{},
	)
}

func (g *game) NewGoblin(p gruid.Point) int {
	return g.ECS.Create(
		Name{"goblin"},
		Position{p},
		Renderable{glyph: 'g', fg: ColorMonster, order: ROActor},
		Health{hp: 10, maxhp: 10},
		Damage{2},
		Perception{LOS: 8},
		AI{state: CSWandering},
		Obstruct{},
	)
}

func (g *game) NewTroll(p gruid.Point) int {
	return g.ECS.Create(
		Name{"troll"},
		Position{p},
		Renderable{glyph: 'T', fg: ColorTroll, order: ROActor},
		Health{hp: 20, maxhp: 20},
		Damage{5},
		Perception{LOS: 6},
		AI{state: CSWandering},
		Obstruct{},
	)
}

func (g *game) NewHealthPotion(p gruid.Point) int {
	return g.ECS.Create(
		Name{"health potion"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: '!', fg: ColorHealthPotion, order: ROItem},
		Collectible{},
		Consumable{hp: 5},
	)
}

func (g *game) NewCorpse(p gruid.Point) int {
	return g.ECS.Create(
		Name{"corpse"},
		Position{p},
		Renderable{glyph: '%', fg: ColorCorpse, order: ROCorpse},
		Collectible{},
		Consumable{hp: 2},
	)
}

func (g *game) NewBlood(p gruid.Point) int {
	return g.ECS.Create(
		Name{"blood"},
		Position{p},
		Renderable{glyph: '.', fg: ColorBlood, order: ROFloor},
	)
}
