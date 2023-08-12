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
			Position{size.Div(2).X, size.Div(2).Y},
			Name{"Player"},
			Renderable{'@', gruid.ColorDefault},
			Input{},
		)
		m.game.ecs.Update()
		// Initialization: create a player entity centered on the map.
		// m.game.ecs.PlayerID = m.game.ECS.AddEntity(&Player{}, size.Div(2))

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
	m.grid.Fill(gruid.Cell{Rune: ' '})
	// We draw the map tiles first.
	it := m.game.Map.Grid.Iterator()
	for it.Next() {
		m.grid.Set(it.P(), gruid.Cell{Rune: m.game.Map.Rune(it.Cell())})
	}
	// We draw the entities.
	for _, i := range m.game.ecs.entities {
		if p := m.game.ecs.positions[i]; p != nil {
			if r := m.game.ecs.renderables[i]; r != nil {
				m.grid.Set(gruid.Point{X: p.x, Y: p.y}, gruid.Cell{
					Rune:  r.glyph,
					Style: gruid.Style{Fg: gruid.ColorDefault},
				})
			}
		}
		// m.grid.Set(m.game.ecs.Positions[i], gruid.Cell{
		// 	Rune:  e.Rune(),
		// 	Style: gruid.Style{Fg: e.Color()},
		// })
	}
	return m.grid
}
