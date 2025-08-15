package main

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/k0kubun/pp/v3"
)

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {

	m.target = nil
	// m.target.path = nil // Remove path highlighting.

	switch msg.Key {

	// Movement
	case gruid.KeyArrowLeft, "h":
		m.action = action{Type: ActionBump, Delta: gruid.Point{-1, 0}}
	case gruid.KeyArrowDown, "j":
		m.action = action{Type: ActionBump, Delta: gruid.Point{0, 1}}
	case gruid.KeyArrowUp, "k":
		m.action = action{Type: ActionBump, Delta: gruid.Point{0, -1}}
	case gruid.KeyArrowRight, "l":
		m.action = action{Type: ActionBump, Delta: gruid.Point{1, 0}}
	case "y":
		m.action = action{Type: ActionBump, Delta: gruid.Point{-1, -1}}
	case "u":
		m.action = action{Type: ActionBump, Delta: gruid.Point{1, -1}}
	case "b":
		m.action = action{Type: ActionBump, Delta: gruid.Point{-1, 1}}
	case "n":
		m.action = action{Type: ActionBump, Delta: gruid.Point{1, 1}}

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

	// Waiting
	case ".":
		m.action = action{Type: ActionWait}

	// Quitting
	case gruid.KeyEscape, "q":
		m.action = action{Type: ActionQuit}
	}
}

// updateTargeting updates targeting information in response to user input
// messages.
func (m *model) updateTargeting(msg gruid.Msg) {
	// Initialize targeting right next to palyer if nil
	maprange := gruid.NewRange(0, 0, UIWidth, UIHeight)
	if m.target == nil {
		// Start cursor position at first visible object
		per := m.game.ECS.GetComponentUnchecked(0, Perception{}).(Perception)
		fmt.Println(per.perceived)
		if len(per.perceived) > 0 {
			other_pos := m.game.ECS.GetComponentUnchecked(per.perceived[0], Position{}).(Position).Point
			m.target = &targeting{
				pos:    other_pos,
				radius: 1,
			}
		} else {
			p := m.game.PlayerPosition()
			m.target = &targeting{
				pos:    p.Add(gruid.Point{X: 2, Y: 2}),
				radius: 1,
			}
		}
	}
	p := m.target.pos
	switch msg := msg.(type) {
	case gruid.MsgKeyDown:
		switch msg.Key {

		case gruid.KeyArrowLeft, "h":
			if p.X > maprange.Min.X {
				p = p.Shift(-1, 0)
			}
		case gruid.KeyArrowRight, "l":
			if p.X < maprange.Max.X {
				p = p.Shift(1, 0)
			}
		case gruid.KeyArrowDown, "j":
			if p.Y < maprange.Max.Y {
				p = p.Shift(0, 1)
			}
		case gruid.KeyArrowUp, "k":
			if p.Y > maprange.Min.Y {
				p = p.Shift(0, -1)
			}
		case "y":
			if p.X > maprange.Min.X && p.Y > maprange.Min.Y {
				p = p.Shift(-1, -1)
			}
		case "u":
			if p.X < maprange.Max.X && p.Y > maprange.Min.Y {
				p = p.Shift(1, -1)
			}
		case "b":
			if p.X > maprange.Min.X && p.Y < maprange.Max.Y {
				p = p.Shift(-1, 1)
			}
		case "n":
			if p.X < maprange.Max.X && p.Y < maprange.Max.Y {
				p = p.Shift(1, 1)
			}

		case gruid.KeyEnter:
			if m.mode == modeExamination {
				break
			}
			m.activateTarget(p)

		case gruid.KeyEscape, "q":
			m.mode = modeNormal
			m.target = nil
			return
		}

		if m.target != nil {
			m.target.pos = p.Add(maprange.Min)
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
			p := m.game.PlayerPosition()
			m.target.path = m.pr.JPSPath(m.target.path, p, m.target.pos, m.game.Pathable, true)
		}
	}
}
