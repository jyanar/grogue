// This file is a set of factories which allow easy spawning of
// different kinds of enemies.

package main

func (g *game) NewPlayer() int {
	return g.ECS.Create(
		Position{g.Map.RandomFloor()},
		Name{"Player"},
		Renderable{glyph: '@', color: ColorPlayer, order: ROActor},
		Health{hp: 18, maxhp: 18},
		FOV{LOS: 20},
		Inventory{},
		Input{},
		Obstruct{},
		Damage{5},
	)
}

func (g *game) NewGoblin() int {
	return g.ECS.Create(
		Position{g.Map.RandomFloor()},
		Name{"Goblin"},
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
		Position{g.Map.RandomFloor()},
		Name{"Troll"},
		Renderable{glyph: 'T', color: ColorTroll, order: ROActor},
		Health{hp: 20, maxhp: 20},
		Damage{5},
		Perception{radius: 6},
		AI{state: CSWandering},
		Obstruct{},
	)
}
