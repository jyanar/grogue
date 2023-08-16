package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/rl"
)

type model struct {
	grid   gruid.Grid // The drawing grid.
	game   game       // The game state.
	action action     // The current UI action.
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
		m.game.ECS = NewECS()
		m.game.ECS.Map = m.game.Map
		m.game.ECS.Create(
			Position{m.game.Map.RandomFloor()},
			Name{"Player"},
			Renderable{glyph: '@', color: ColorPlayer, order: ROActor},
			Health{hp: 20, maxhp: 20},
			FOV{LOS: 10, FOV: rl.NewFOV(gruid.NewRange(-10, -10, 10+1, 10+1))},
			Input{},
			Damage{5},
		)
		m.game.SpawnEnemies()
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
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, 0)}
	case gruid.KeyArrowDown, "j":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(0, 1)}
	case gruid.KeyArrowUp, "k":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(0, -1)}
	case gruid.KeyArrowRight, "l":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, 0)}
	case gruid.KeyEscape, "q":
		m.action = action{Type: ActionQuit}
	case "y":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, -1)}
	case "u":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, -1)}
	case "b":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(-1, 1)}
	case "n":
		m.action = action{Type: ActionBump, Delta: pdelta.Shift(1, 1)}
	}
}

// Draw implements gruid.Model.Draw. It draws a simple map that spans the whole
// grid.
func (m *model) Draw() gruid.Grid {
	ECS := m.game.ECS
	Map := m.game.Map

	m.grid.Fill(gruid.Cell{Rune: ' '}) // Clear the map.

	// Draw the map.
	it := Map.Grid.Iterator()
	for it.Next() {
		if !Map.Explored[it.P()] {
			continue
		}
		c := gruid.Cell{Rune: Map.Rune(it.Cell())}
		if m.game.InFOV(it.P()) {
			c.Style.Bg = ColorFOV
		}
		m.grid.Set(it.P(), c)
	}
	// Draw the entities.
	// TODO Refactor this ugly mess.
	// Collect list of entities to draw.
	corpsesToDraw := []int{}
	itemsToDraw := []int{}
	actorsToDraw := []int{}
	for _, e := range ECS.EntitiesWith(Position{}, Renderable{}) {
		p := ECS.positions[e]
		if !m.game.Map.Explored[p.Point] || !m.game.InFOV(p.Point) {
			continue
		}
		// Entity is in a FOV. Add them to the list.
		switch ECS.renderables[e].order {
		case ROCorpse:
			corpsesToDraw = append(corpsesToDraw, e)
		case ROItem:
			itemsToDraw = append(itemsToDraw, e)
		case ROActor:
			actorsToDraw = append(actorsToDraw, e)
		}
	}
	// // Sort them according to drawing order.
	// fmt.Println(entitiesToDraw)
	// for _, e := range entitiesToDraw {
	// 	ECS.printDebug(e)
	// }
	// sort.Slice(entitiesToDraw, func(i, j int) bool { // Why is this segfaulting?
	// 	return ECS.renderables[i].order < ECS.renderables[j].order
	// })
	// Draw.
	for _, e := range corpsesToDraw {
		p := ECS.positions[e]
		r := ECS.renderables[e]
		m.grid.Set(p.Point, gruid.Cell{
			Rune:  r.glyph,
			Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
		})
	}
	for _, e := range itemsToDraw {
		p := ECS.positions[e]
		r := ECS.renderables[e]
		m.grid.Set(p.Point, gruid.Cell{
			Rune:  r.glyph,
			Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
		})
	}
	for _, e := range actorsToDraw {
		p := ECS.positions[e]
		r := ECS.renderables[e]
		m.grid.Set(p.Point, gruid.Cell{
			Rune:  r.glyph,
			Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
		})
	}
	// for _, e := range entitiesToDraw {
	// 	p := ECS.positions[e]
	// 	r := ECS.renderables[e]
	// 	m.grid.Set(p.Point, gruid.Cell{
	// 		Rune:  r.glyph,
	// 		Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
	// 	})
	// }
	return m.grid
}
