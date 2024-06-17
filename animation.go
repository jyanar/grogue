package main

import (
	"github.com/anaseto/gruid"
)

// A single frame in an animation, specifying both the duration of the frame
// and the cells/locations which to draw, via the FrameCells.
type AnimationFrame struct {
	framecells []gruid.FrameCell
}

// An animation, containing a set of frames.
type Animation struct {
	frames []AnimationFrame
}

type InterruptibleAnimation struct {
	CAnimation
}

func NewFrameCell(cell gruid.Cell, p gruid.Point) CFrameCell {
	return CFrameCell{Renderable{cell: cell, order: ROActor}, p}
}

func (g *game) NewExampleAnimation(p gruid.Point) int {
	return g.ECS.Create(
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
