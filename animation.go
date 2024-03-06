package main

import (
	"time"

	"github.com/anaseto/gruid"
)

// An animation message for passing into the model's Update() method. A value of
// true indicates that the animation is ongoing, whereas false indicates that
// the animation has terminated.
type msgAnimation bool

// A single frame in an animation, specifying both the duration of the frame
// and the cells/locations which to draw, via the FrameCells.
type AnimationFrame struct {
	duration   time.Duration
	framecells []gruid.FrameCell
}

// An animation, containing a set of frames.
type Animation struct {
	frames []AnimationFrame
}

const (
	AnimDurShort       = 25 * time.Millisecond
	AnimDurShortMedium = 50 * time.Millisecond
	AnimDurMedium      = 75 * time.Millisecond
	AnimDurMediumLong  = 100 * time.Millisecond
)

func NewAnimation() *Animation {
	return &Animation{}
}

func (a *Animation) AddFrame(frame AnimationFrame) {
	a.frames = append(a.frames, frame)
}

func (a *Animation) Done() bool {
	return len(a.frames) == 0
}

func (a *Animation) animCmdStart() gruid.Cmd {
	return gruid.Cmd(func() gruid.Msg {
		t := time.NewTimer(a.frames[0].duration)
		<-t.C
		return msgAnimation(true)
	})
}

func (a *Animation) animCmdContinue(duration time.Duration) gruid.Cmd {
	return gruid.Cmd(func() gruid.Msg {
		t := time.NewTimer(duration)
		<-t.C
		return msgAnimation(true)
	})
}

func (a *Animation) animCmdEnd() gruid.Cmd {
	return gruid.Cmd(func() gruid.Msg {
		return msgAnimation(false)
	})
}

var GoreRunes = []rune{'.', '~', '`', '&'}

// A gore animation that occurs whenever someone gets killed.
// We send bits of gore flying in different directions (primarily in the
// direction from attacker to defender).
// func GoreAnimation(attacker, defender gruid.Point) *Animation {
// 	a := NewAnimation()
// 	a.AddFrame(AnimationFrame{
// 		duration:   AnimDurMedium,
// 		framecells: []gruid.FrameCell{
// 			gruid.FrameCell{Cell:gruid.Cell{Rune: ''}},
// 		},
// 	})
// }

func NewAttackedAnimation(p gruid.Point) *Animation {
	a := NewAnimation()
	a.AddFrame(AnimationFrame{
		duration: 300 * time.Millisecond,
		framecells: []gruid.FrameCell{
			gruid.FrameCell{Cell: gruid.Cell{Rune: ' '}, P: p},
		},
	})
	return a
}

// Might make sense to pass the map to this function as well -- that way we can
// decide whether a given point is a wall or whatever.
func NewDeathAnimation(p gruid.Point, m *Map) *Animation {

	// Okay, so we will need to check: where is the enemy, where is the player,
	// where are there walls, etc!

	a := NewAnimation()
	a.AddFrame(AnimationFrame{
		duration: 300 * time.Millisecond,
		framecells: []gruid.FrameCell{
			gruid.FrameCell{Cell: gruid.Cell{Rune: '~'}, P: p.Shift(0, 1)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '~'}, P: p.Shift(1, 0)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '~'}, P: p.Shift(-1, 0)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '~'}, P: p.Shift(0, -1)},
		},
	})
	a.AddFrame(AnimationFrame{
		duration: 300 * time.Millisecond,
		framecells: []gruid.FrameCell{
			gruid.FrameCell{Cell: gruid.Cell{Rune: '*'}, P: p.Shift(0, 1)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '*'}, P: p.Shift(1, 0)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '*'}, P: p.Shift(-1, 0)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '*'}, P: p.Shift(0, -1)},
		},
	})
	a.AddFrame(AnimationFrame{
		duration: 300 * time.Millisecond,
		framecells: []gruid.FrameCell{
			gruid.FrameCell{Cell: gruid.Cell{Rune: '&'}, P: p.Shift(0, 1)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '&'}, P: p.Shift(1, 0)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '&'}, P: p.Shift(-1, 0)},
			gruid.FrameCell{Cell: gruid.Cell{Rune: '&'}, P: p.Shift(0, -1)},
		},
	})
	return a
}
