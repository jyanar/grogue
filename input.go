package main

import (
	"codeberg.org/anaseto/gruid"
	"github.com/k0kubun/pp/v3"
)

// keyToDir maps a movement key to a unit direction vector.
// Returns the direction and true if the key is a movement key, otherwise false.
func keyToDir(key gruid.Key) (gruid.Point, bool) {
	switch key {
	case gruid.KeyArrowLeft, "h":
		return gruid.Point{X: -1, Y: 0}, true
	case gruid.KeyArrowDown, "j":
		return gruid.Point{X: 0, Y: 1}, true
	case gruid.KeyArrowUp, "k":
		return gruid.Point{X: 0, Y: -1}, true
	case gruid.KeyArrowRight, "l":
		return gruid.Point{X: 1, Y: 0}, true
	case "y":
		return gruid.Point{X: -1, Y: -1}, true
	case "u":
		return gruid.Point{X: 1, Y: -1}, true
	case "b":
		return gruid.Point{X: -1, Y: 1}, true
	case "n":
		return gruid.Point{X: 1, Y: 1}, true
	}
	return gruid.Point{}, false
}

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {

	m.target = nil

	switch msg.Key {

	case "a":
		m.action = action{Type: ActionIAnimate}

	// Message log, inventory, pick up items, and examine
	case "m":
		m.action = action{Type: ActionViewMessages}
	case "i":
		m.action = action{Type: ActionInventory}
	case "d":
		m.action = action{Type: ActionDrop}
	case "g":
		m.action = action{Type: ActionPickup}
	case "x":
		m.action = action{Type: ActionExamine}

	case "t":
		pp.Print(m.game.ECS.GetComponentsFor(0))

	case `\`:
		m.debugRevealAll = !m.debugRevealAll

	// Waiting
	case ".":
		m.action = action{Type: ActionWait}

	// Quitting
	case gruid.KeyEscape, "q":
		m.action = action{Type: ActionQuit}

	// Movement
	default:
		if dir, ok := keyToDir(msg.Key); ok {
			m.action = action{Type: ActionBump, Delta: dir}
		}
	}
}
