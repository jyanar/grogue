package main

import (
	"fmt"

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
	case gruid.KeyEscape:
		if m.mouseActive {
			m.mouseActive = false
			m.target = nil
		} else {
			m.action = action{Type: ActionQuit}
		}
	case "q":
		m.action = action{Type: ActionQuit}

	// Movement
	default:
		if dir, ok := keyToDir(msg.Key); ok {
			m.action = action{Type: ActionBump, Delta: dir}
		}
	}
}

// updateTargeting updates targeting information in response to user input
// messages.
func (m *model) updateTargeting(msg gruid.Msg) {
	maprg := gruid.NewRange(0, 0, UIWidth, UIHeight)
	if m.target == nil {
		m.target = &targeting{}
	}
	if !m.target.pos.In(maprg) {
		pos, _ := m.game.ECS.GetComponent(0, Position{})
		m.target.pos = pos.(Position).Point.Add(maprg.Min)
	}
	p := m.target.pos
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		switch msg.Key {

		case gruid.KeyEnter:
			if m.mode == modeExamination {
				break
			}
			m.activateTarget(p)

		case gruid.KeyEscape, "q":
			m.mode = modeNormal
			m.target = nil
			return

		default:
			if dir, ok := keyToDir(msg.Key); ok {
				p = p.Add(dir)
			}
		}

		if m.target != nil {
			m.target.pos = p.Add(maprg.Min)
		}

	case gruid.MsgMouse:
		switch msg.Action {
		case gruid.MouseMove:
			m.target.pos = msg.P.Shift(-1, -1)

		case gruid.MouseMain:
			fmt.Printf("Mouse click at: %v\n", msg.P)
		}
	}

	// We only compute the full path if the player is still alive.
	if m.game.ECS.PlayerDead() {
		m.target = nil
	} else {
		if m.target != nil {
			p, _ := m.game.ECS.GetComponent(0, Position{})
			pos := p.(Position)
			m.target.path = m.pr.JPSPath(m.target.path, pos.Point, m.target.pos, m.game.Pathable, true)
		}
	}
}
