package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

type Position struct {
	gruid.Point
}

type Renderable struct {
	glyph rune
	color gruid.Color
}

type Name struct {
	string
}

type FOV struct {
	LOS int
	FOV *rl.FOV
}

// Entities with this component can accept keyboard input.
type Input struct{}

// Represents a directional action.
type Bump struct {
	gruid.Point
}
