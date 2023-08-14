package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

type model struct {
	grid   gruid.Grid // The drawing grid.
	game   game       // The game state.
	action action     // The current UI action.
}

type game struct {
	ECS *ECS
	Map *Map
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
		// Initialize map.
		size := m.grid.Size()
		m.game.Map = NewMap(size)
		// Initialize entities.
		m.game.ECS = NewECS()
		m.game.ECS.Create(
			Position{m.game.Map.RandomFloor()},
			Name{"Player"},
			Renderable{'@', gruid.ColorDefault},
			FOV{LOS: 10, FOV: rl.NewFOV(gruid.NewRange(-10, -10, 10+1, 10+1))},
			Input{},
		)
		m.game.ECS.Create(
			Position{m.game.Map.RandomFloor()},
			Name{"Goblin"},
			Renderable{'g', gruid.ColorDefault},
			// FOV{LOS: 3, FOV: rl.NewFOV(gruid.NewRange(-3, -3, 3+1, 3+1))},
		)
		m.game.ECS.drawgrid = &m.grid
		m.game.ECS.Map = m.game.Map
		m.game.ECS.Update()

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

const (
	ColorFOV gruid.Color = iota + 1
)

// Draw implements gruid.Model.Draw. It draws a simple map that spans the whole
// grid.
func (m *model) Draw() gruid.Grid {
	// m.grid.Fill(gruid.Cell{Rune: ' '})
	// m.game.Map.Draw(&m.grid) // Draw the map.
	// m.game.ECS.Draw(&m.grid) // Draw the entities.
	return m.grid
}

// InFOV returns true if p is in the field of view of an entity with FOV. We only
// keep cells within maxLOS manhattan distance from the source entity.
func (g *game) InFOV(p gruid.Point) bool {
	pp := g.ECS.positions[0].Point
	// While we iterate over all entities with FOV components, typically only the
	// player has an FOV component, so generally the iterator returns on index 0.
	for _, f := range g.ECS.fovs {
		if f.FOV.Visible(p) && paths.DistanceManhattan(pp, p) <= 10 {
			return true
		}
	}
	return false
}
