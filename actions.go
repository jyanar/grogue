package main

import (
	"github.com/anaseto/gruid"
)

type action struct {
	Type  actionType  // Kind of action (bump, quit, open inventory, etc)
	Delta gruid.Point // direction for ActionBump
}

type actionType int

const (
	NoAction   actionType = iota
	ActionBump            // Movement request.
	ActionWait            // Step forward one tick.
	ActionQuit            // Quit the game.
)

func (m *model) handleAction() gruid.Effect {
	switch m.action.Type {
	case ActionBump:
		// Add a bump component to all entities with an Input component
		for _, e := range m.game.ECS.EntitiesWith(Input{}) {
			m.game.ECS.AddComponent(e, Bump{m.action.Delta})
		}
		m.game.ECS.Update()

	case ActionWait:
		m.game.ECS.Update()

	case ActionQuit:
		return gruid.End()
	}
	return nil
}
