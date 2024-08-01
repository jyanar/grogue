// Some utility methods for dealing with and rendering the inventory.

package main

import (
	"log"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

func (m *model) OpenInventory(title string) {
	// Build list of entries in player inventory.
	iC, _ := m.game.ECS.GetComponent(0, Inventory{})
	inv := iC.(Inventory)
	entries := []ui.MenuEntry{}
	r := 'a'
	for _, it := range inv.items {
		nC, _ := m.game.ECS.GetComponent(it, Name{})
		rC, _ := m.game.ECS.GetComponent(it, Renderable{})
		name := nC.(Name).string
		renderable := rC.(Renderable)
		glyph := renderable.cell.Rune
		fg := renderable.cell.Style.Fg
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
			// Check whether the given item has a ranged component
			// item_idx := m.game.ECS.inventories[0].items[n]
			iC, _ := m.game.ECS.GetComponent(0, Inventory{})
			inv := iC.(Inventory)
			item_idx := inv.items[n]

			if m.game.ECS.HasComponent(item_idx, Ranged{}) {
				pC, _ := m.game.ECS.GetComponent(0, Position{})
				p := pC.(Position).Point
				m.target = &targeting{
					pos:    p.Shift(1, 1),
					radius: 2,
					item:   item_idx,
				}
				m.mode = modeTargeting
				return
			}
			err = m.game.InventoryActivate(0, n)
		}
		if err != nil {
			m.game.Logf(err.Error(), ColorLogSpecial)
		}
		m.game.ECS.Update()
		m.mode = modeNormal
	}
}

func (m *model) activateTarget(p gruid.Point) {
	log.Println("Activating target at point p!")
	log.Println(p)
	// Check if there is an entity here capable of taking damage
	item := m.target.item
	// item_dmg := m.game.ECS.damages[item].int
	dmg, _ := m.game.ECS.GetComponent(item, Damage{})
	item_dmg := dmg.(Damage).int
	if entities := m.game.ECS.EntitiesAtPWith(p, Health{}); len(entities) > 0 {
		for _, e := range entities {
			m.game.ECS.AddComponent(e, DamageEffect{0, item_dmg})
		}
	}
	m.target = nil
	m.mode = modeNormal
	// Get rid of item in inventory
	if m.game.ECS.HasComponent(item, Consumable{}) {
		iC, _ := m.game.ECS.GetComponent(0, Inventory{})
		inv := iC.(Inventory)
		inv.items = remove(inv.items, item)
		m.game.ECS.AddComponent(0, inv)
		m.game.ECS.Delete(item)
	}
	m.game.ECS.Update()
	// Can we force the game to re-render now?
}

const ErrNoShow = "ErrNoShow"

// TODO Better log messages
func (g *game) InventoryActivate(entity, itemidx int) error {

	item := g.ECS.inventories[entity].items[itemidx]
	item_name := g.ECS.names[item].string
	entity_name := g.ECS.names[entity].string
	var prefix string
	if entity == 0 {
		prefix = "You use the"
	} else {
		prefix = entity_name + " uses the"
	}
	g.Logf("%s %s.", ColorLogSpecial, prefix, item_name)
	// Item can provide healing. Apply healing.
	if g.ECS.HasComponent(item, Healing{}) {
		g.ECS.healths[entity].hp += g.ECS.healings[item].amount
	}
	// TODO Ranged effects
	// if g.ECS.HasComponent(item, Ranged{}) {
	// }
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
