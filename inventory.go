// Some utility methods for dealing with and rendering the inventory.

package main

import (
	"log"
	"sort"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
)

// sortedInventoryKeys returns the inventory's assigned letters in sorted order.
func sortedInventoryKeys(inv Inventory) []rune {
	keys := make([]rune, 0, len(inv.items))
	for k := range inv.items {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (m *model) OpenInventory(title string) {
	// Build list of entries in player inventory.
	inv := GetComponent[Inventory](m.game.ECS, 0)
	entries := []ui.MenuEntry{}
	for _, k := range sortedInventoryKeys(inv) {
		it := inv.items[k]
		name := GetComponent[Name](m.game.ECS, it).string
		renderable := GetComponent[Renderable](m.game.ECS, it)
		glyph := renderable.cell.Rune
		fg := renderable.cell.Style.Fg
		stt := ui.Text("").WithMarkup('k', gruid.Style{}.WithFg(fg))
		entries = append(entries, ui.MenuEntry{
			Text: stt.WithText(string(k) + " - @k" + string(glyph) + "@N " + name),
			Keys: []gruid.Key{gruid.Key(k)},
		})
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
		inv := m.game.PlayerInventory()
		key := sortedInventoryKeys(inv)[m.inventory.Active()]
		itemid := inv.items[key]
		var err error
		switch m.mode {
		case modeInventoryDrop:
			err = m.game.InventoryDrop(0, key)
		case modeInventoryActivate:
			// Check whether the given item has a ranged component
			if m.game.ECS.HasComponent(itemid, Ranged{}) {
				p := m.game.PlayerPosition()
				m.target = &targeting{
					pos:    p.Shift(1, 1),
					radius: 1,
					itemid: itemid,
				}
				m.mode = modeTargeting
				return
			}
			err = m.game.InventoryActivate(0, key)
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
	itemid := m.target.itemid
	itemdmg := GetComponent[Damage](m.game.ECS, itemid).int
	if entities := m.game.ECS.EntitiesAtPWith(p, Health{}); len(entities) > 0 {
		for _, e := range entities {
			m.game.ECS.AddComponent(e, DamageEffect{0, itemdmg})
		}
	}
	m.target = nil
	m.mode = modeNormal
	// Remove item from inventory and world
	if m.game.ECS.HasComponent(itemid, Consumable{}) {
		inv := m.game.PlayerInventory()
		inv.removeItem(itemid)
		m.game.ECS.AddComponent(0, inv)
		m.game.ECS.Delete(itemid)
	}
	m.game.ECS.Update()
}

const ErrNoShow = "ErrNoShow"

// TODO Better log messages
func (g *game) InventoryActivate(entity int, key rune) error {
	inventory := GetComponent[Inventory](g.ECS, entity)
	item_id := inventory.items[key]
	item_name := GetComponent[Name](g.ECS, item_id).string
	g.Logf("You use the %s.", ColorLogSpecial, item_name)
	// Item can provide healing. Apply healing.
	if g.ECS.HasComponent(item_id, Healing{}) {
		health := GetComponent[Health](g.ECS, entity)
		healing := GetComponent[Healing](g.ECS, item_id)
		health.hp += healing.amount
		if health.hp > health.maxhp {
			health.hp = health.maxhp
		}
		g.ECS.AddComponent(entity, health)
	}
	// Item was consumable, so we delete from inventory.
	if g.ECS.HasComponent(item_id, Consumable{}) {
		delete(inventory.items, key)
		g.ECS.AddComponent(entity, inventory)
		g.ECS.Delete(item_id)
	}
	return nil
}

func (g *game) InventoryDrop(entity int, key rune) error {
	inventory := GetComponent[Inventory](g.ECS, entity)
	item_id := inventory.items[key]
	item_name := GetComponent[Name](g.ECS, item_id).string
	prefix := "You drop the "
	g.Logf("%s %s.", ColorLogSpecial, prefix, item_name)
	// Remove item from inventory.
	delete(inventory.items, key)
	g.ECS.AddComponent(entity, inventory)
	pos := GetComponent[Position](g.ECS, entity).Point
	// Add Position component back to the item.
	g.ECS.AddComponent(item_id, Position{pos})
	return nil
}
