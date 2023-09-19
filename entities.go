// This file is a set of factories which allow easy spawning of
// different kinds of enemies.

package main

func (g *game) NewPlayer() int {
	return g.ECS.Create(
		Name{"Player"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: '@', color: ColorPlayer, order: ROActor},
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
		Name{"Goblin"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: 'g', color: ColorMonster, order: ROActor},
		Health{hp: 10, maxhp: 10},
		Damage{2},
		Perception{radius: 8},
		AI{state: CSWandering},
		Obstruct{},
	)
}

func (g *game) NewTroll() int {
	return g.ECS.Create(
		Name{"Troll"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: 'T', color: ColorTroll, order: ROActor},
		Health{hp: 20, maxhp: 20},
		Damage{5},
		Perception{radius: 6},
		AI{state: CSWandering},
		Obstruct{},
	)
}

func (g *game) NewHealthPotion() int {
	return g.ECS.Create(
		Name{"Health Potion"},
		Position{g.FreeFloorTile()},
		Renderable{glyph: '!', color: ColorHealthPotion, order: ROItem},
		Collectible{},
		Consumable{hp: 5},
	)
}

// Useful for debugging that corpses pathing bug
// 2023 Sep 18 observed a corpse, goblin, player, and troll standing on the same tile.
func (g *game) NewCorpse() int {
	return g.ECS.Create(
		Name{"Corpse"},
		Position{g.Map.RandomFloor()},
		Renderable{glyph: '%', color: ColorCorpse, order: ROCorpse},
		Collectible{},
		Consumable{hp: 2},
	)
}
