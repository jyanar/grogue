// This file defines the main model of the game: the Update function that
// updates the model state in response to user input, and the Draw function,
// which draw the final grid.
//
// That is, it fulfills the following interface:
//
//	  // Model contains the application's state.
//	  type Model interface {
//		  // Update is called when a message is received. Use it to update your
//		  // model in response to messages and/or send commands or subscriptions.
//		  // It is always called the first time with a MsgInit message.
//		  Update(Msg) Effect
//
//		  // Draw is called after every Update. Use this function to draw the UI
//		  // elements in a grid to be returned. If only parts of the grid are to
//		  // be updated, you can return a smaller grid slice, or an empty grid
//		  // slice to skip any drawing work. Note that the contents of the grid
//		  // slice are then compared to the previous state at the same bounds,
//		  // and only the changes are sent to the driver anyway.
//		  Draw() Grid
//	  }
//

package main

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type model struct {
	grid      gruid.Grid  // The drawing grid.
	game      game        // The game state.
	action    action      // The current UI action.
	mode      mode        // The current UI mode.
	log       *ui.Label   // Label for the log.
	status    *ui.Label   // Label for the status.
	desc      *ui.Label   // Label for position description.
	viewer    *ui.Pager   // Message's history viewer.
	inventory *ui.Menu    // Inventory menu.
	mousePos  gruid.Point // Mouse position.
}

type mode int

const (
	modeNormal            mode = iota // Controlling the player.
	modeEnd                           // Win or death (currently only death).
	modeMessageViewer                 // Currently viewing messages.
	modeInventoryActivate             // Browsing inventory, in order to use an item.
	modeInventoryDrop                 // Browsing inventory, in order to drop an item.
)

func NewModel(g gruid.Grid) *model {
	return &model{
		grid: g,
	}
}

// Update implements gruid.Model.update. It handles keyboard and mouse input
// messages and updates the model in response to them.
func (m *model) Update(msg gruid.Msg) gruid.Effect {

	m.action = action{}
	switch m.mode {
	case modeEnd:
		switch msg := msg.(type) {
		case gruid.MsgKeyDown:
			switch msg.Key {
			case "q", gruid.KeyEscape:
				// You died: quit on "q" or "escape"
				return gruid.End()
			case ".":
				// Otherwise, allow player to continue watching sim.
				m.updateMsgKeyDown(msg)
			}
		}
		return nil

	case modeMessageViewer:
		m.viewer.Update(msg) // e.g., scrolling.
		if m.viewer.Action() == ui.PagerQuit {
			m.mode = modeNormal
		}
		return nil

	case modeInventoryActivate, modeInventoryDrop:
		m.updateInventory(msg)
		return nil

	case modeNormal:
		switch msg := msg.(type) {

		case gruid.MsgInit:
			m.log = &ui.Label{}
			m.status = &ui.Label{}
			m.desc = &ui.Label{Box: &ui.Box{}}
			m.InitializeMessageViewer()
			m.game = game{}
			// Initialize map.
			m.game.Map = NewMap(gruid.Point{X: MapWidth, Y: MapHeight})
			m.game.ECS = NewECS()
			m.game.ECS.Map = m.game.Map
			// Place player on a random floor.
			m.game.ECS.Create(
				Position{m.game.Map.RandomFloor()},
				Name{"Player"},
				Renderable{glyph: '@', color: ColorPlayer, order: ROActor},
				Health{hp: 18, maxhp: 18},
				FOV{LOS: 20},
				Inventory{},
				Input{},
				Obstruct{},
				Damage{5},
			)
			// Spawn enemies, place items, and advance a tick.
			m.game.SpawnEnemies()
			m.game.PlaceItems()
			m.game.ECS.Update()

		case gruid.MsgKeyDown:
			m.updateMsgKeyDown(msg)

		case gruid.MsgMouse:
			if msg.Action == gruid.MouseMove {
				m.mousePos = msg.P
			}
		}
	}

	// Handle action (if any provided).
	return m.handleAction()
}

// DRAW METHODS ------------------

// Draw implements gruid.Model.Draw. It draws a simple map that spans the whole
// grid.
func (m *model) Draw() gruid.Grid {
	ECS := m.game.ECS
	Map := m.game.Map

	if m.mode == modeMessageViewer {
		m.grid.Copy(m.viewer.Draw())
		return m.grid
	}

	if m.mode == modeInventoryActivate || m.mode == modeInventoryDrop {
		m.grid.Copy(m.inventory.Draw())
		return m.grid
	}

	m.grid.Fill(gruid.Cell{Rune: ' '}) // Clear the map.
	mapgrid := m.grid.Slice(m.grid.Range().Shift(1, 1, 0, 0))

	// Draw the map.
	it := Map.Grid.Iterator()
	for it.Next() {
		if !Map.Explored[it.P()] {
			continue
		}
		c := gruid.Cell{Rune: Map.Rune(it.Cell())}
		if m.game.InFOV(it.P()) {
			c.Style.Fg = ColorFOV
			c.Style.Bg = ColorFOV
		}
		mapgrid.Set(it.P(), c)
	}
	// Draw the entities.
	// TODO Refactor this ugly mess to use sorting.
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
	// Draw.
	for _, e := range corpsesToDraw {
		p := ECS.positions[e]
		r := ECS.renderables[e]
		mapgrid.Set(p.Point, gruid.Cell{
			Rune:  r.glyph,
			Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
		})
	}
	for _, e := range itemsToDraw {
		p := ECS.positions[e]
		r := ECS.renderables[e]
		mapgrid.Set(p.Point, gruid.Cell{
			Rune:  r.glyph,
			Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
		})
	}
	for _, e := range actorsToDraw {
		p := ECS.positions[e]
		r := ECS.renderables[e]
		mapgrid.Set(p.Point, gruid.Cell{
			Rune:  r.glyph,
			Style: gruid.Style{Fg: r.color, Bg: ColorFOV},
		})
	}
	m.DrawNames(mapgrid)
	m.DrawLog(m.grid.Slice(m.grid.Range().Lines(m.grid.Size().Y-4, m.grid.Size().Y-1)))
	m.DrawStatus(m.grid.Slice(m.grid.Range().Line(m.grid.Size().Y - 1)))
	return m.grid
}

// DrawLog draws the last two lines of the log.
func (m *model) DrawLog(gd gruid.Grid) {
	j := 1
	for i := len(m.game.Log) - 1; i >= 0; i-- {
		if j < 0 {
			break
		}
		e := m.game.Log[i]
		st := gruid.Style{}
		st.Fg = e.Color
		m.log.Content = ui.NewStyledText(e.String(), st)
		m.log.Draw(gd.Slice(gd.Range().Line(j)))
		j--
	}
}

// DrawStatus draws the status line.
func (m *model) DrawStatus(gd gruid.Grid) {
	// Write the HP on top of that.
	st := gruid.Style{Fg: ColorStatusHealthy}
	st.Bg = ColorLogMonsterAttack
	player_health := m.game.ECS.healths[0]
	if player_health.hp < player_health.maxhp/2 {
		st.Fg = ColorStatusWounded
	}
	m.log.Content = ui.Textf("HP: %d/%d", player_health.hp, player_health.maxhp).WithStyle(st)
	m.log.Draw(gd)
}

// DrawNames renders the names of the named entities at the current mouse location
// if it is in the map.
func (m *model) DrawNames(gd gruid.Grid) {
	maprg := gruid.NewRange(0, 2, UIWidth, UIHeight-1)
	if !m.mousePos.In(maprg) {
		return
	}
	// p := m.mousePos.Sub(gruid.Point{X: 0, Y: 2})
	p := m.mousePos.Shift(-1, -1)
	// We get the names of the entities at p.
	names := []string{}
	for _, e := range m.game.ECS.EntitiesWith(Position{}) {
		q := m.game.ECS.positions[e]
		if q.Point != p || !m.game.InFOV(q.Point) {
			continue
		}
		if name, ok := m.game.ECS.names[e]; ok {
			names = append(names, name.string)
		}
	}
	if len(names) == 0 {
		return
	}
	// We sort the names. This could be improved to sort by entity type
	// too, as well as to remove duplicates (for example showing "corpse
	// (3x)" if there are three corpses).
	sort.Strings(names)

	text := strings.Join(names, ", ")
	width := utf8.RuneCountInString(text) + 2
	rg := gruid.NewRange(p.X+1, p.Y-1, p.X+1+width, p.Y+2)
	// We adjust a bit the box's placement in case it's on an edge.
	if p.X+1+width >= UIWidth {
		rg = rg.Shift(-1-width, 0, -1-width, 0)
	}
	if p.Y+2 > MapHeight {
		rg = rg.Shift(0, -1, 0, -1)
	}
	if p.Y-1 < 0 {
		rg = rg.Shift(0, 1, 0, 1)
	}
	slice := gd.Slice(rg)
	m.desc.Content = ui.Text(text)
	m.desc.Draw(slice)
}
