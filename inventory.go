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
			err = m.game.InventoryDrop(0, n)
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
	inventory := g.ECS.GetComponentUnchecked(entity, Inventory{}).(Inventory)
	item_id := inventory.items[itemidx]
	item_name := g.ECS.GetComponentUnchecked(item_id, Name{}).(Name).string
	g.Logf("You use the %s.", ColorLogSpecial, item_name)
	// Item can provide healing. Apply healing.
	if g.ECS.HasComponent(item_id, Healing{}) {
		health := g.ECS.GetComponentUnchecked(entity, Health{}).(Health)
		healing := g.ECS.GetComponentUnchecked(item_id, Healing{}).(Healing)
		health.hp += healing.amount
		if health.hp > health.maxhp {
			health.hp = health.maxhp
		}
		g.ECS.AddComponent(entity, health)
	}
	// TODO Ranged effects
	// Item was consumable, so we delete from inventory.
	if g.ECS.HasComponent(item_id, Consumable{}) {
		inventory.items = remove(inventory.items, item_id)
		g.ECS.AddComponent(entity, inventory)
		g.ECS.Delete(item_id)
	}
	return nil
}

// func (g *game) InventoryActivateWithTarget(entity, itemidx int) error {
// 	item := g.ECS.inventories[entity].items[itemidx]
// }

func (g *game) InventoryDrop(entity, itemidx int) error {
	inventory := g.ECS.GetComponentUnchecked(entity, Inventory{}).(Inventory)
	item_id := inventory.items[itemidx]
	item_name := g.ECS.GetComponentUnchecked(item_id, Name{}).(Name).string
	prefix := "You drop the "
	g.Logf("%s %s.", ColorLogSpecial, prefix, item_name)
	// Remove item from inventory.
	inventory.items = remove(inventory.items, item_id)
	g.ECS.AddComponent(entity, inventory)
	pos := g.ECS.GetComponentUnchecked(entity, Position{}).(Position).Point
	// Add Position component back to the item.
	g.ECS.AddComponent(item_id, Position{pos})
	return nil
}
