package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

type Position struct {
	gruid.Point
}

type renderOrder int

const (
	ROCorpse renderOrder = iota
	ROItem
	ROActor
)

type Renderable struct {
	glyph rune
	color gruid.Color
	order renderOrder
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

// Entities with this component obstruct movement. Corpses do not.
type Obstruct struct{}

// Entities with this component have health, and can take damage.
type Health struct {
	hp, maxhp int
}

// Entities with this component can damage entities with a health component.
type Damage struct {
	int
}

// Entities with this component are marked for death (see DeathSystem).
type Death struct{}

// Entities with this component perceive other entities around them.
type Perception struct {
	radius    int   // Perception radius.
	perceived []int // List of perceived entities.
}

type creatureState int

const (
	CSWandering = iota
	CSSleeping
	CSHunting
)

// Entities with this component will be controlled by AI, and can wander,
// sleep, or hunt the player.
type AI struct {
	state creatureState
	dest  *gruid.Point
	// path  []gruid.Point
}
