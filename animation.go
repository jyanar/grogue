package main

import (
	"log"

	"github.com/anaseto/gruid"
)

func (m *model) handleMsgTick() {
	// Update background animations
	m.game.ECS.UpdateAnimation()
	// Update interruptible animation
	if m.ianimation != nil {
		log.Println("Updating interruptible animation!!")
		// Advance animation by a single tick.
		anim := m.ianimation
		anim.frames[anim.index].itick++

		// If the current frame has expired, move to the next frame.
		if anim.frames[anim.index].itick >= anim.frames[anim.index].nticks {
			anim.frames[anim.index].itick = 0
			anim.index++
		}

		// If the current animation has expired, remove it from the ECS or restart.
		if anim.index >= len(anim.frames) {
			anim.index = 0
			if anim.repeat == 0 {
				m.ianimation = nil
			} else if anim.repeat > 0 {
				anim.repeat--
			}
		}
	}
}

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
