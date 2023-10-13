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
