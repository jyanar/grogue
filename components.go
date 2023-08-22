package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

// An entity position, defined as an X, Y coordinate.
type Position struct {
	gruid.Point
}

// Determines order at which entity is drawn. Corpses are drawn first,
// followed by items and actors.
type renderOrder string

const (
	ROCorpse renderOrder = "CORPSE"
	ROItem   renderOrder = "ITEM"
	ROActor  renderOrder = "ACTOR"
)

// Entities with this component can be rendered.
type Renderable struct {
	glyph rune
	color gruid.Color
	order renderOrder
}

type Name struct {
	string
}

// Entities with this component have an FOV computed. Typically, only the
// player holds this component.
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

// I don't like this iota business. Prefer enums like this:
type creatureState string

const (
	CSWandering creatureState = "WANDERING"
	CSSleeping  creatureState = "SLEEPING"
	CSHunting   creatureState = "HUNTING"
)

// Entities with this component will be controlled by AI, and can wander,
// sleep, or hunt the player.
type AI struct {
	state creatureState
	dest  *gruid.Point
	// path  []gruid.Point
}
