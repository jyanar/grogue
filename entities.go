// This file is a set of factories which allow easy spawning of
// different kinds of enemies.

package main

import "github.com/anaseto/gruid"

func (g *game) NewPlayer(p gruid.Point) int {
	return g.ECS.Create(
		Name{"you"},
		Position{p},
		Renderable{cell: gruid.Cell{Rune: '@', Style: gruid.Style{Fg: ColorPlayer}}, order: ROActor},
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
		Renderable{
			cell:  gruid.Cell{Rune: 'g', Style: gruid.Style{Fg: ColorMonster}},
			order: ROActor,
		},
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
		Renderable{
			cell:  gruid.Cell{Rune: 'T', Style: gruid.Style{Fg: ColorTroll}},
			order: ROActor,
		},
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
		Renderable{
			cell:  gruid.Cell{Rune: '!', Style: gruid.Style{Fg: ColorHealthPotion}},
			order: ROItem,
		},
		Collectible{},
		Consumable{},
		Healing{amount: 5},
	)
}

func (g *game) NewCorpse(p gruid.Point) int {
	return g.ECS.Create(
		Name{"corpse"},
		Position{p},
		Renderable{
			cell:  gruid.Cell{Rune: '%', Style: gruid.Style{Fg: ColorCorpse}},
			order: ROCorpse,
		},
		Collectible{},
		Consumable{},
		Healing{amount: 2},
	)
}

func (g *game) NewBlood(p gruid.Point) int {
	return g.ECS.Create(
		Name{"blood"},
		Position{p},
		Renderable{cell: gruid.Cell{Rune: '.', Style: gruid.Style{Fg: ColorBlood}}, order: ROFloor},
	)
}

func (g *game) NewScroll(p gruid.Point) int {
	return g.ECS.Create(
		Name{"scroll"},
		Position{p},
		Renderable{cell: gruid.Cell{Rune: '?', Style: gruid.Style{Fg: ColorScroll}}, order: ROItem},
		Collectible{},
		Consumable{},
		Ranged{Range: 6},
		Damage{5},
		AreaOfEffect{radius: 3},
	)
}

func NewFrameCell(cell gruid.Cell, p gruid.Point) CFrameCell {
	return CFrameCell{Renderable{cell: cell, order: ROActor}, p}
}

func (g *game) NewExampleAnimation(p gruid.Point) int {
	return g.ECS.Create(
		Name{"example animation"},
		Position{p},
		CAnimation{
			index: 0,
			frames: []CAnimationFrame{
				{
					framecells: []CFrameCell{
						{Renderable{cell: gruid.Cell{Rune: '1', Style: gruid.Style{Fg: ColorBlood}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
				{
					framecells: []CFrameCell{
						{Renderable{cell: gruid.Cell{Rune: '0', Style: gruid.Style{Fg: ColorPlayer}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
				{
					framecells: []CFrameCell{
						{Renderable{cell: gruid.Cell{Rune: 'P', Style: gruid.Style{Fg: ColorTroll}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
			},
			repeat: 20,
		},
	)
}

func NewExampleIAnimation(p gruid.Point) *InterruptibleAnimation {
	return &InterruptibleAnimation{
		CAnimation{
			index: 0,
			frames: []CAnimationFrame{
				{
					framecells: []CFrameCell{
						{Renderable{cell: gruid.Cell{Rune: '1', Style: gruid.Style{Fg: ColorBlood}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
				{
					framecells: []CFrameCell{
						{Renderable{cell: gruid.Cell{Rune: '0', Style: gruid.Style{Fg: ColorPlayer}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
				{
					framecells: []CFrameCell{
						{Renderable{cell: gruid.Cell{Rune: 'P', Style: gruid.Style{Fg: ColorTroll}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
			},
			repeat: 20,
		},
	}
}
