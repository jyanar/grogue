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
		// Add a bump component to all entities with an Input component
		m.game.ecs.AddComponent(0, Bump{dx: m.action.Delta.X, dy: m.action.Delta.Y})
		m.game.ecs.Update()

		// np := m.game.ECS.Positions[m.game.ECS.PlayerID]
		// np = np.Add(m.action.Delta)
		// if m.game.Map.Walkable(np) {
		// 	m.game.ECS.MovePlayer(np)
		// }

	case ActionQuit:
		return gruid.End()
	}
	return nil
}
