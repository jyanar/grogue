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
			Position{m.game.Map.RandomFloor()},
			Name{"Player"},
			Renderable{'@', gruid.ColorDefault},
			FOV{LOS: 10, FOV: rl.NewFOV(gruid.NewRange(-10, -10, 10+1, 10+1))},
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

const (
	ColorFOV gruid.Color = iota + 1
)

// Draw implements gruid.Model.Draw. It draws a simple map that spans the whole
// grid.
func (m *model) Draw() gruid.Grid {
	g := m.game
	// Draw the map.
	m.grid.Fill(gruid.Cell{Rune: ' '})
	it := g.Map.Grid.Iterator()
	for it.Next() {
		if !g.Map.Explored[it.P()] {
			continue
		}
		c := gruid.Cell{Rune: g.Map.Rune(it.Cell())}
		if g.InFOV(it.P()) {
			c.Style.Bg = ColorFOV
		}
		m.grid.Set(it.P(), c)
	}
	// Draw the entities.
	for _, e := range m.game.ecs.entities {
		if m.game.ecs.HasComponent(e, Position{}) && m.game.ecs.HasComponent(e, Renderable{}) {
			p := m.game.ecs.positions[e]
			r := m.game.ecs.renderables[e]
			m.grid.Set(p.Point, gruid.Cell{
				Rune:  r.glyph,
				Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
			})
		}
	}
	return m.grid
}

// InFOV returns true if p is in the player's field of view. We only keep cells
// within maxLOS manhattan distance from the player, as natural given our
// current 4-way movement. With 8-way movement, the natural distance choice
func (g *game) InFOV(p gruid.Point) bool {
	pp := g.ecs.positions[0].Point
	// return true
	return g.ecs.fovs[0].FOV.Visible(p) &&
		paths.DistanceManhattan(pp, p) <= 10
}
