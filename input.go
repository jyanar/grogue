package main

import "github.com/anaseto/gruid"

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {

	m.target.path = nil // Remove path highlighting.

	pdelta := gruid.Point{}
	switch msg.Key {

	// Movement
	case gruid.KeyArrowLeft, "h":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, 0)}
	case gruid.KeyArrowDown, "j":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(0, 1)}
	case gruid.KeyArrowUp, "k":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(0, -1)}
	case gruid.KeyArrowRight, "l":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, 0)}
	case "y":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, -1)}
	case "u":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, -1)}
	case "b":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, 1)}
	case "n":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, 1)}

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

	// Waiting
	case ".":
		m.action = action{Type: ActionWait, Delta: pdelta.Shift(0, 0)}

	// Quitting
	case gruid.KeyEscape, "q":
		m.action = action{Type: ActionQuit}
	}
}
