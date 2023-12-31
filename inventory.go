// Some utility methods for dealing with and rendering the inventory.

package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

func (m *model) OpenInventory(title string) {
	// Build list of entries in player inventory.
	inv := m.game.ECS.inventories[0]
	entries := []ui.MenuEntry{}
	r := 'a'
	for _, it := range inv.items {
		name := m.game.ECS.names[it].string
		glyph := m.game.ECS.renderables[it].glyph
		fg := m.game.ECS.renderables[it].fg
		stt := ui.Text("").WithMarkup('k', gruid.Style{}.WithFg(fg))
		entries = append(entries, ui.MenuEntry{
			Text: stt.WithText(string(r) + " - @k" + string(glyph) + "@N " + name),
			Keys: []gruid.Key{gruid.Key(r)},
		})
		r++
	}
	// We create a new menu widget for the inventory window.
	m.inventory = ui.NewMenu(ui.MenuConfig{
		Grid:    gruid.NewGrid(40, MapHeight),
		Box:     &ui.Box{Title: ui.Text(title).WithStyle(gruid.Style{}.WithFg(ColorPlayer))},
		Entries: entries,
	})
}

// updateInventory handles input messages when the inventory window is open.
func (m *model) updateInventory(msg gruid.Msg) {
	// We call the Update function of the menu widget, so that we can
	// inspect information about user activity on the menu.
	m.inventory.Update(msg)
	switch m.inventory.Action() {
	case ui.MenuQuit:
		// The user requested to quit the menu.
		m.mode = modeNormal
		return
	case ui.MenuInvoke:
		// The user invoked a particular entry of the menu (either by
		// using enter or clicking on it).
		n := m.inventory.Active()
		var err error
		switch m.mode {
		case modeInventoryDrop:
			err = m.game.InventoryRemove(0, n)
		case modeInventoryActivate:
			err = m.game.InventoryActivate(0, n)
		}
		if err != nil {
			m.game.Logf(err.Error(), ColorLogSpecial)
		}
		m.game.ECS.Update()
		m.mode = modeNormal
	}
}

const ErrNoShow = "ErrNoShow"

func (g *game) InventoryActivate(entity, itemidx int) error {
	item := g.ECS.inventories[entity].items[itemidx]
	item_name := g.ECS.names[item].string
	entity_name := g.ECS.names[entity].string
	// Item can provide healing. Apply healing.
	if g.ECS.HasComponent(item, Healing{}) {
		g.ECS.healths[entity].hp += g.ECS.healings[item].amount
		g.Logf("%s uses %s.", ColorLogSpecial, entity_name, item_name)
	}
	// Item has ranged effect.
	// What we would like:
	// 	Move into modeTargeting, select an entity, and then activate spell at
	//    the target location.
	if g.ECS.HasComponent(item, Ranged{}) {
		g.Logf("%s uses %s.", ColorLogSpecial, entity_name, item_name)
	}
	// Item was consumable, so we delete from inventory.
	if g.ECS.HasComponent(item, Consumable{}) {
		g.ECS.inventories[entity].items = remove(g.ECS.inventories[entity].items, item)
		g.ECS.Delete(item)
	}
	return nil
}

// func (g *game) InventoryActivateWithTarget(entity, itemidx int) error {
// 	item := g.ECS.inventories[entity].items[itemidx]
// }

func (g *game) InventoryRemove(entity, itemidx int) error {
	item := g.ECS.inventories[entity].items[itemidx]
	item_name := g.ECS.names[item].string
	entity_name := g.ECS.names[entity].string
	// Add Position component back to the item.
	pos := g.ECS.positions[entity].Point
	g.ECS.AddComponent(item, Position{pos})
	// Remove item from inventory
	g.Logf("%s drops %s", ColorLogSpecial, entity_name, item_name)
	g.ECS.inventories[entity].items = remove(g.ECS.inventories[entity].items, item)
	return nil
}
