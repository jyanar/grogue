package main

import (
	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/ui"
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
	ActionExamine                 // Examine the map.
	ActionIAnimate                // Start an interruptible animation.
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
		if !m.game.ECS.PlayerDead() {
			m.OpenInventory("Use item")
			m.mode = modeInventoryActivate
			m.game.CollectMessages()
		}

	case ActionDrop:
		m.OpenInventory("Drop item")
		m.mode = modeInventoryDrop
		m.game.CollectMessages()

	case ActionPickup:
		ok := m.game.PickupItem()
		m.game.CollectMessages()
		if ok {
			m.game.ECS.Update()
		}

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

	case ActionExamine:
		m.mode = modeExamination

	case ActionIAnimate:
		p, hasPos := m.game.ECS.GetComponent(0, Position{})
		if hasPos {
			m.ianimation = NewExampleIAnimation(p.(Position).Point)
		}
	}
	return nil
}
