// This file is a set of factories which allow easy spawning of
// different kinds of enemies.

package main

import (
	"github.com/anaseto/gruid"
)

// func (g *game) NewBaseCreature(p gruid.Point) int {
// 	return g.ECS.Create(
// 		Position{p},
// 		Visible{},
// 		Inventory{},
// 		Obstruct{},
// 	)
// }

// Convenience methods
func NewRenderable(r rune, fg, bg gruid.Color, order renderOrder) Renderable {
	return Renderable{cell: gruid.Cell{Rune: r, Style: gruid.Style{Fg: fg, Bg: bg}}, order: order}
}

func NewRenderableNoBg(Rune rune, fg gruid.Color, order renderOrder) Renderable {
	return Renderable{cell: gruid.Cell{Rune: Rune, Style: gruid.Style{Fg: fg, Bg: ColorNone}}, order: order}
}

func (g *game) NewPlayer(p gruid.Point) int {
	return g.ECS.Create(
		Name{"you"},
		Position{p},
		Visible{},
		NewRenderableNoBg('@', ColorPlayer, ROActor),
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
		Visible{},
		NewRenderableNoBg('g', ColorMonster, ROActor),
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
		Visible{},
		NewRenderableNoBg('T', ColorTroll, ROActor),
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
		Visible{},
		NewRenderableNoBg('!', ColorHealthPotion, ROItem),
		Collectible{},
		Consumable{},
		Healing{amount: 5},
	)
}

func (g *game) NewCorpse(p gruid.Point) int {
	return g.ECS.Create(
		Name{"corpse"},
		Position{p},
		Visible{},
		NewRenderableNoBg('%', ColorCorpse, ROCorpse),
		Collectible{},
		Consumable{},
		Healing{amount: 2},
	)
}

func (g *game) NewBlood(p gruid.Point) int {
	return g.ECS.Create(
		Name{"blood"},
		Visible{},
		Position{p},
		NewRenderable('.', ColorBlood, ColorBlood, ROFloor),
	)
}

func (g *game) NewScroll(p gruid.Point) int {
	return g.ECS.Create(
		Name{"scroll"},
		Position{p},
		Visible{},
		NewRenderableNoBg('?', ColorScroll, ROItem),
		Collectible{},
		Consumable{},
		Ranged{Range: 6},
		Damage{5},
		AreaOfEffect{radius: 3},
	)
}

func (g *game) NewWaterTile(p gruid.Point) int {
	return g.ECS.Create(
		Name{"water"},
		Visible{},
		Position{p},
		// NewRenderable('~', ColorWater1, ColorWater1, ROFloor),
		Animation{
			index:  0,
			repeat: -1,
			frames: []Frame{
				{
					itick:  0,
					nticks: 5,
					framecells: []FrameCell{
						{
							r: NewRenderable('~', ColorWater1, ColorWater1, ROFloor),
							p: p,
						},
					},
				},
				{
					itick:  0,
					nticks: 7,
					framecells: []FrameCell{
						{
							r: NewRenderable('~', ColorWater2, ColorWater2, ROFloor),
							p: p,
						},
					},
				},
			},
		},
	)
}
