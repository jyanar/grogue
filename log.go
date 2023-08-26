package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/ui"
)

func (e LogEntry) String() string {
	if e.Dups == 0 {
		return e.Text
	}
	return fmt.Sprintf("%s (%dx)", e.Text, e.Dups)
}

// log adds an entry to the player's log.
func (g *game) log(e LogEntry) {
	if len(g.Log) > 0 {
		if g.Log[len(g.Log)-1].Text == e.Text {
			g.Log[len(g.Log)-1].Dups++
			return
		}
	}
	g.Log = append(g.Log, e)
}

// Logf adds a formatted entry to the game log.
func (g *game) Logf(format string, color gruid.Color, a ...any) {
	e := LogEntry{Text: fmt.Sprintf(format, a...), Color: color}
	g.log(e)
}

// InitializeMessageViewer creates a new pager for viewing message's history
func (m *model) InitializeMessageViewer() {
	m.viewer = ui.NewPager(ui.PagerConfig{
		Grid: gruid.NewGrid(UIWidth, UIHeight-1),
		Box:  &ui.Box{},
	})
}

// CollectMessages iterates through the ECS and collects any entities with MessageLog components.
func (g *game) CollectMessages() {
	for _, e := range g.ECS.EntitiesWith(LogEntry{}) {
		msg := g.ECS.logentries[e]
		g.Logf(msg.Text, msg.Color)
		g.ECS.Delete(e)
	}
}
