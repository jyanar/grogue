package main

import (
	"github.com/anaseto/gruid"
)

// A collection of animations that can be triggered or created.

func NewFrameCell(cell gruid.Cell, p gruid.Point) FrameCell {
	return FrameCell{Renderable{cell: cell, order: ROActor}, p}
}

func (g *game) NewExampleAnimation(p gruid.Point) int {
	return g.ECS.Create(
		Position{p},
		Animation{
			index: 0,
			frames: []Frame{
				{
					framecells: []FrameCell{
						{Renderable{cell: gruid.Cell{Rune: '1', Style: gruid.Style{Fg: ColorBlood}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
				{
					framecells: []FrameCell{
						{Renderable{cell: gruid.Cell{Rune: '0', Style: gruid.Style{Fg: ColorPlayer}}, order: ROActor}, p},
					},
					itick:  0,
					nticks: 1,
				},
				{
					framecells: []FrameCell{
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

func NewExampleIAnimation(p gruid.Point) *Animation {
	return &Animation{
		index:  0,
		repeat: 20,
		frames: []Frame{
			{
				framecells: []FrameCell{
					{Renderable{cell: gruid.Cell{Rune: '1', Style: gruid.Style{Fg: ColorBlood}}, order: ROActor}, p},
				},
				itick:  0,
				nticks: 1,
			},
			{
				framecells: []FrameCell{
					{Renderable{cell: gruid.Cell{Rune: '0', Style: gruid.Style{Fg: ColorPlayer}}, order: ROActor}, p},
				},
				itick:  0,
				nticks: 1,
			},
			{
				framecells: []FrameCell{
					{Renderable{cell: gruid.Cell{Rune: 'P', Style: gruid.Style{Fg: ColorTroll}}, order: ROActor}, p},
				},
				itick:  0,
				nticks: 1,
			},
		},
	}
}
