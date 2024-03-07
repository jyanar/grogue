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
	ROFloor  renderOrder = "FLOOR" // An item on the floor. Blood, grass, etc.
	ROCorpse renderOrder = "CORPSE"
	ROItem   renderOrder = "ITEM"
	ROActor  renderOrder = "ACTOR"
)

// Entities with this component can be rendered.
type Renderable struct {
	cell  gruid.Cell
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
	LOS       int     // Perceptive radius.
	FOV       *rl.FOV // Effective FOV, which can be affected by occlusion.
	perceived []int   // List of perceived entities.
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
}

// This component represents a message, to be processed by and added to the message log.
type LogEntry struct {
	Text  string      // entry text
	Color gruid.Color // color
	Dups  int         // consecutive duplicates of same message
}

// Entities with this component are consumed on use.
type Consumable struct{}

// Entities with this component provide healing when used.
type Healing struct {
	amount int
}

// Entities with this component can be picked up and placed in inventory.
type Collectible struct{}

// Entities with this component have an inventory, and can pick up Collectible
// components.
type Inventory struct {
	items []int // A list of entities.
}

// Entities with this component can be used for ranged attacks. e.g. staffs.
type Ranged struct {
	Range int
}

// Entities with this component will perform an action.
type Action struct {
	action actionType
}

// Entities with this component can be thrown.
type Throwable struct{}

// // Entities with this component can be used for ranged attacks. e.g. staffs.
// type Zappable struct {
// 	Range int
// }

// Entities with this component have an area of effect which is activated when
// it is zapped (such as in the case of staffs) or thrown (such as in the case
// of potions)
type AreaOfEffect struct {
	radius int
}

// Entities with this component will take damage.
type DamageEffect struct {
	source int
	amount int
}

type CFrameCell struct {
	r Renderable
	p gruid.Point
}

type CAnimationFrame struct {
	framecells []CFrameCell
	nticks     int // Duration of animation, in ticks
	itick      int // Current duration. Resets to duration when 0.
}

type CAnimation struct {
	frames []CAnimationFrame
	index  int
	repeat int // -1 for infinite, 0 for no repeat, n for n repeats
}

type InterruptibleAnimation struct {
	CAnimation
}
