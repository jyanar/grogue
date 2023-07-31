package main

import "github.com/anaseto/gruid"

type action struct {
	Type  actionType
	Delta gruid.Point
}

type actionType int

const (
	NoAction       actionType = iota
	ActionMovement            // Movement request.
	ActionQuit                // Quit the game.
)

func (m *model) handleAction() gruid.Effect {
	switch m.action.Type {
	case ActionMovement:
		m.game.PlayerPos = m.game.PlayerPos.Add(m.action.Delta)
	case ActionQuit:
		return gruid.End()
	}
	return nil
}
