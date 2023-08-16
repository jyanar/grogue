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

// Obstructs movement
type Obstruct struct{}

// Entities with this component have health, and can take damage.
type Health struct {
	hp, maxhp int
}

// Entities with this component can damage entities with a health component.
type Damage struct {
	int
}

// Entities with this component will be processed as dead.
type Death struct{}
