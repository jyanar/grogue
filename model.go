package main

import "github.com/anaseto/gruid"

type model struct {
	grid   gruid.Grid // The drawing grid.
	game   game       // The game state.
	action action     // The current UI action.
}

type game struct {
	PlayerPos gruid.Point
}

func NewModel(g gruid.Grid) *model {
	return &model{
		grid: g,
	}
}

func (m *model) Update(msg gruid.Msg) gruid.Effect {
	m.action = action{}
	switch msg := msg.(type) {
	case gruid.MsgInit:
		m.game.PlayerPos = m.grid.Size().Div(2)
	case gruid.MsgKeyDown:
		m.updateMsgKeyDown(msg)
	}
	// Handle action (if any provided).
	return m.handleAction()
}

func (m *model) updateMsgKeyDown(msg gruid.MsgKeyDown) {
	pdelta := gruid.Point{}
	switch msg.Key {
	case gruid.KeyArrowLeft, "h":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(-1, 0)}
	case gruid.KeyArrowDown, "j":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(0, 1)}
	case gruid.KeyArrowUp, "k":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(0, -1)}
	case gruid.KeyArrowRight, "l":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(1, 0)}
	case gruid.KeyEscape, "q":
		m.action = action{Type: ActionQuit}
	case "y":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(-1, -1)}
	case "u":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(1, -1)}
	case "b":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(-1, 1)}
	case "n":
		m.action = action{Type: ActionMovement, Delta: pdelta.Shift(1, 1)}
	}
}

// Draw implements gruid.Model.Draw. It draws a simple map that spans the whole
// grid.
func (m *model) Draw() gruid.Grid {
	it := m.grid.Iterator()
	for it.Next() {
		switch {
		case it.P() == m.game.PlayerPos:
			it.SetCell(gruid.Cell{Rune: '@'})
		default:
			it.SetCell(gruid.Cell{Rune: ' '})
		}
	}
	return m.grid
}
