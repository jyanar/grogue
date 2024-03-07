package main

import (
	"time"

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
	ActionExamine                 // Examine the map.
	ActionAnimate                 // Execute an animation.
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
		m.OpenInventory("Use item")
		m.mode = modeInventoryActivate
		m.game.CollectMessages()

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
		m.target.pos = m.game.ECS.positions[0].Point.Shift(2, 2)

	case ActionAnimate:
		m.mode = modeAnimation
		return gruid.Cmd(func() gruid.Msg {
			t := time.NewTimer(m.animation.frames[0].duration)
			<-t.C
			return msgAnimation(true)
		})

	case ActionIAnimate:
		m.ianimation = NewExampleIAnimation(m.game.ECS.positions[0].Point)
	}
	return nil
}
