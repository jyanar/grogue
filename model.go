package main

import (
	"context"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"codeberg.org/anaseto/gruid/ui"
)

type model struct {
	grid       gruid.Grid       // The drawing grid.
	game       game             // The game state.
	action     action           // The current UI action.
	mode       mode             // The current UI mode.
	log        *ui.Label        // Label for the log.
	status     *ui.Label        // Label for the status.
	desc       *ui.Label        // Label for position description.
	viewer     *ui.Pager        // Message's history viewer.
	inventory  *ui.Menu         // Inventory menu.
	pr         *paths.PathRange // Pathing algorithm.
	target     *targeting       // Mouse position.
	ianimation *Animation       // Interruptible animation.
}

// targeting describes information related to examination or selection of
// particular positions in the map.
type targeting struct {
	pos    gruid.Point   // The current position of the cursor.
	path   []gruid.Point // The path to the current position.
	itemid int           // The entity ID of the item being used/thrown/activated.
	radius int           // Radius of the targeting area.
}

// mode describes distinct kinds of modes for the UI. It is used to send user
// input messages to different handlers (inventory window, map, message viewer,
// etc.), depending on the current mode.
type mode int

const (
	modeNormal            mode = iota // Controlling the player.
	modeEnd                           // Win or death (currently only death).
	modeMessageViewer                 // Currently viewing messages.
	modeInventoryActivate             // Browsing inventory, in order to use an item.
	modeInventoryDrop                 // Browsing inventory, in order to drop an item.
	modeExamination                   // Keyboard map examination mode.
	modeTargeting
)

func NewModel(gd gruid.Grid) *model {
	return &model{
		grid:   gd,
		log:    &ui.Label{},
		status: &ui.Label{},
		desc:   &ui.Label{Box: &ui.Box{}},
		viewer: ui.NewPager(ui.PagerConfig{
			Grid: gruid.NewGrid(UIWidth, UIHeight-1),
			Box:  &ui.Box{},
		}),
		pr: paths.NewPathRange(gd.Range()),
	}
}

type msgTick struct{}

func frameTicker() gruid.Sub {
	return func(ctx context.Context, ch chan<- gruid.Msg) {
		t := time.NewTicker(100 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				ch <- msgTick{}
			}
		}
	}
}

// Update implements gruid.Model.update. It handles keyboard and mouse input
// messages and updates the model in response to them.
func (m *model) Update(msg gruid.Msg) gruid.Effect {

	m.action = action{}

	switch m.mode {

	case modeNormal:
		switch msg := msg.(type) {

		case gruid.MsgInit:
			m.game.Initialize()
			return frameTicker()

		case gruid.MsgKeyDown:
			// Interrupt animation on any key press.
			if m.ianimation != nil {
				m.ianimation = nil
			}
			m.updateMsgKeyDown(msg)

		case gruid.MsgMouse:
			m.updateTargeting(msg)

		case msgTick:
			m.handleMsgTick()
		}

	case modeMessageViewer:
		m.viewer.Update(msg) // e.g., scrolling.
		if m.viewer.Action() == ui.PagerQuit {
			m.mode = modeNormal
		}
		return nil

	case modeInventoryActivate, modeInventoryDrop:
		m.updateInventory(msg)
		return nil

	case modeTargeting, modeExamination:
		m.updateTargeting(msg)
		return nil

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
	}

	// Handle action (if any provided).
	return m.handleAction()
}

/////////////////////////////////////
// DRAW METHODS ------------------ //
/////////////////////////////////////

// Draw implements gruid.Model.Draw.
// It draws a simple map that spans the whole grid.
func (m *model) Draw() gruid.Grid {
	ECS := m.game.ECS
	Map := m.game.Map

	// Render message viewer, if that's the mode we're in.
	if m.mode == modeMessageViewer {
		m.grid.Copy(m.viewer.Draw())
		return m.grid
	}

	// Render the inventory, if that's the mode we're in.
	if m.mode == modeInventoryActivate || m.mode == modeInventoryDrop {
		m.grid.Copy(m.inventory.Draw())
		return m.grid
	}

	///////////////////////////////////////////////////
	// Otherwise, render the map, entities, and log. //
	///////////////////////////////////////////////////

	m.grid.Fill(gruid.Cell{Rune: ' '}) // Clear the map.

	mapgrid := m.grid.Slice(m.grid.Range().Shift(1, 1, 0, 0))
	loggrid := m.grid.Slice(m.grid.Range().Lines(m.grid.Size().Y-4, m.grid.Size().Y-1))
	statusgrid := m.grid.Slice(m.grid.Range().Line(m.grid.Size().Y - 1))

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

	// Draw background animations
	m.DrawBackgroundAnimations(mapgrid)

	// Collect entities to draw.
	type tup struct {
		entity int
		order  renderOrder
	}
	entitiesToDraw := []tup{}
	for _, e := range ECS.EntitiesWith(Position{}, Renderable{}) {
		p := ECS.GetComponentUnchecked(e, Position{}).(Position)
		r := ECS.GetComponentUnchecked(e, Renderable{}).(Renderable)
		// If entities is not in FOV, do not add them to list.
		if !m.game.Map.Explored[p.Point] || !m.game.InFOV(p.Point) {
			continue
		}
		entitiesToDraw = append(entitiesToDraw, tup{e, r.order})
	}
	sort.SliceStable(entitiesToDraw, func(i, j int) bool {
		return entitiesToDraw[i].order > entitiesToDraw[j].order
	})
	// Draw entities.
	for _, e := range entitiesToDraw {
		p := ECS.GetComponentUnchecked(e.entity, Position{}).(Position)
		r := ECS.GetComponentUnchecked(e.entity, Renderable{}).(Renderable)
		c := gruid.Cell{Rune: r.cell.Rune, Style: r.cell.Style}
		if r.LacksBg() {
			c.Style.Bg = mapgrid.At(p.Point).Style.Bg
		}
		mapgrid.Set(p.Point, c)
	}

	// Draw target (if targeting), names, log, and status.
	m.DrawTarget(mapgrid)
	m.DrawNames(mapgrid)
	m.DrawLog(loggrid)
	m.DrawStatus(statusgrid)

	// Draw background and player-triggered animations
	m.DrawInterruptibleAnimation(mapgrid)

	return m.grid
}

const (
	AttrNone gruid.AttrMask = iota
	AttrReverse
)

// DrawTarget draws the current position of the mouse.
func (m *model) DrawTarget(gd gruid.Grid) {
	if m.target == nil {
		return
	}
	for _, p := range m.target.path {
		c := gd.At(p)
		gd.Set(p, c.WithStyle(c.Style.WithAttrs(AttrReverse)))
	}
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

// DrawStatus writes the HP on the bottom of the screen. If the player is dead,
// displays "DEAD" in red.
func (m *model) DrawStatus(gd gruid.Grid) {
	if m.game.ECS.PlayerDead() {
		m.log.Content = ui.Text("  DEAD  ").WithStyle(gruid.Style{Fg: ColorBlood})
	} else {
		st := gruid.Style{Fg: ColorStatusHealthy}
		st.Bg = ColorLogMonsterAttack
		player_health := m.game.ECS.GetComponentUnchecked(0, Health{}).(Health)
		if player_health.hp < player_health.maxhp/2 {
			st.Fg = ColorStatusWounded
		}
		m.log.Content = ui.Textf("HP: %d/%d", player_health.hp, player_health.maxhp).WithStyle(st)
	}
	m.log.Draw(gd)
}

// DrawNames renders the names of the named entities at the current mouse location
// if it is in the map.
func (m *model) DrawNames(gd gruid.Grid) {
	maprg := gruid.NewRange(0, 2, UIWidth, UIHeight-1)
	if m.target == nil || !m.target.pos.In(maprg) {
		return
	}
	p := m.target.pos
	// We get the names of the entities at p.
	names := []string{}
	for _, e := range m.game.ECS.EntitiesAtPWith(p, Name{}) {
		if !m.game.InFOV(p) {
			continue
		}
		name := m.game.ECS.GetComponentUnchecked(e, Name{}).(Name).string
		names = append(names, name)
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

// Animations such as lava, water, torches. This function iterates over all such
// entities and draws them. See AnimationSystem to see how these are updated.
func (m *model) DrawBackgroundAnimations(gd gruid.Grid) {
	for _, e := range m.game.ECS.EntitiesWith(Animation{}) {
		anim := m.game.ECS.GetComponentUnchecked(e, Animation{}).(Animation)
		for _, fc := range anim.frames[anim.index].framecells {
			p := fc.p
			r := fc.r
			if m.game.Map.Explored[p] && m.game.InFOV(p) {
				gd.Set(p, r.cell)
			}
		}
	}
}

// Interruptible animations are generally those which can be interrupted by
// player input. For example, a potion exploding, or the player throwing
// something, etc. Only one interruptible animation can be active at any point
// in time.
func (m *model) DrawInterruptibleAnimation(gd gruid.Grid) {
	// Iterate over all framecells of the current frame, and draw.
	if m.ianimation != nil {
		anim := m.ianimation
		for _, fc := range anim.frames[anim.index].framecells {
			p := fc.p
			r := fc.r
			gd.Set(p, r.cell)
		}
	}
}
