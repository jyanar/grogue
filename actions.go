package main

import (
	"github.com/anaseto/gruid"
)

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
		// Add a bump component to all entities with an Input component
		for _, e := range m.game.ECS.entities {
			if m.game.ECS.inputs[e] != nil {
				m.game.ECS.AddComponent(e, Bump{m.action.Delta})
			}
		}
		m.game.ECS.Update()

	case ActionQuit:
		return gruid.End()
	}
	return nil
}
