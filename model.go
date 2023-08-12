package main

import (
	"github.com/anaseto/gruid"
)

type model struct {
	grid   gruid.Grid // The drawing grid.
	game   game       // The game state.
	action action     // The current UI action.
}

type game struct {
	ecs *ECS
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
		m.game.ecs = NewECS()
		m.game.ecs.Create(
			Position{size.Div(2)},
			Name{"Player"},
			Renderable{'@', gruid.ColorDefault},
			Input{},
		)
		m.game.ecs.Map = m.game.Map
		m.game.ecs.Update()

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
	// Draw the map.
	m.grid.Fill(gruid.Cell{Rune: ' '})
	it := m.game.Map.Grid.Iterator()
	for it.Next() {
		m.grid.Set(it.P(), gruid.Cell{Rune: m.game.Map.Rune(it.Cell())})
	}
	// Draw the entities.
	for _, e := range m.game.ecs.entities {
		if m.game.ecs.HasComponent(e, Position{}) && m.game.ecs.HasComponent(e, Renderable{}) {
			p := m.game.ecs.positions[e]
			r := m.game.ecs.renderables[e]
			m.grid.Set(p.Point, gruid.Cell{
				Rune:  r.glyph,
				Style: gruid.Style{Fg: r.color},
			})
		}
	}
	return m.grid
}
