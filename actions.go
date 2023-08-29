package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

type action struct {
	Type  actionType  // Kind of action (bump, quit, open inventory, etc)
	Delta gruid.Point // direction for ActionBump
}

type actionType int

const (
	NoAction           actionType = iota
	ActionBump                    // Movement request.
	ActionWait                    // Step forward one tick.
	ActionQuit                    // Quit the game.
	ActionViewMessages            // View history messages.
	ActionInventory               // Open inventory.
	ActionPickup                  // Pick up an item.
	ActionDrop                    // Drop an item.
)

func (m *model) handleAction() gruid.Effect {
	switch m.action.Type {
	case ActionBump:
		// Add a bump component to all entities with an Input component
		for _, e := range m.game.ECS.EntitiesWith(Input{}) {
			m.game.ECS.AddComponent(e, Bump{m.action.Delta})
		}
		m.game.ECS.Update()
		m.game.CollectMessages()

	case ActionWait:
		m.game.ECS.Update()
		m.game.CollectMessages()

	case ActionInventory:
		fmt.Println("OPEN INVENTORY++++++++++++++++++++++++++=")

	case ActionPickup:
		fmt.Println("PICKING UP ITEM++++++++++++++++++++++++++")
		m.game.PickupItem()

	case ActionDrop:
		fmt.Println("DROP AN ITEM++++++++++++++++++++++++++++++")

	case ActionViewMessages:
		m.mode = modeMessageViewer
		lines := []ui.StyledText{}
		for _, e := range m.game.Log {
			st := gruid.Style{}
			st.Fg = e.Color
			lines = append(lines, ui.NewStyledText(e.String(), st))
		}
		m.viewer.SetLines(lines)

	case ActionQuit:
		return gruid.End()

	}
	return nil
}
