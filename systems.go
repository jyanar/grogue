package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

type System interface {
	Update()
}

type BumpSystem struct {
	ecs *ECS
}

func (s *BumpSystem) Update() {
	for _, e := range s.ecs.EntitiesWith(Bump{}, Position{}) {
		// Get entity's bump and position data.
		b := s.ecs.bumps[e]
		p := s.ecs.positions[e]
		dest := p.Point.Add(b.Point)
		s.ecs.bumps[e] = nil // Consume bump component.
		// Ignore movement to the same tile.
		if b.X == 0 && b.Y == 0 {
			continue
		}
		// Let's attempt to move to dest.
		if s.ecs.Map.Walkable(dest) {
			// Check whether there are blocking entities at dest.
			attackable := s.ecs.EntitiesAtPWith(dest, Health{}, Obstruct{})
			if len(attackable) == 0 {
				p.Point = dest // No entity blocking the way, move to dest.
				continue
			}
			if len(attackable) > 1 {
				panic(fmt.Sprintf("More than one entity with obstruct at position: %v", dest))
			}
			// Attack entity at location.
			target := attackable[0]
			dmg_src := s.ecs.damages[e].int
			name_src := s.ecs.names[e].string
			name_target := s.ecs.names[target].string
			health_target := s.ecs.healths[target]
			msg := fmt.Sprintf("%s hits the %s for %d damage!", name_src, name_target, dmg_src)
			msgcolor := ColorLogMonsterAttack
			if e == 0 {
				msgcolor = ColorLogPlayerAttack
			}
			s.ecs.Create(LogEntry{Text: msg, Color: msgcolor})
			health_target.hp -= dmg_src
			if !s.ecs.BloodAt(dest) { // Add blood tile here.
				s.ecs.Create(
					Name{"blood"},
					Position{dest},
					Renderable{glyph: '.', fg: ColorBlood, order: ROFloor},
				)
			}
			if health_target.hp <= 0 {
				health_target.hp = 0
				s.ecs.AddComponent(target, Death{}) // Entity marked for death.
			}
		} else {
			s.ecs.Create(LogEntry{Text: "The wall is firm and unyielding!", Color: ColorLogSpecial})
		}
	}
}

type FOVSystem struct {
	ecs *ECS
}

func (s *FOVSystem) Update() {
	for _, e := range s.ecs.EntitiesWith(FOV{}, Position{}) {
		p := s.ecs.positions[e]
		f := s.ecs.fovs[e]
		if f.FOV == nil {
			f.FOV = rl.NewFOV(gruid.NewRange(-f.LOS, -f.LOS, f.LOS+1, f.LOS+1))
		}
		// We shift the FOV's range so that it will be centered on the position
		// of the entity.
		rg := gruid.NewRange(-f.LOS, -f.LOS, f.LOS+1, f.LOS+1)
		f.FOV.SetRange(rg.Add(p.Point).Intersect(s.ecs.Map.Grid.Range()))
		// We mark cells in field of view as explored. We use the symmetric shadow
		// casting algorithm provided by the rl package.
		passable := func(p gruid.Point) bool {
			return s.ecs.Map.Grid.At(p) != Wall
		}
		for _, point := range f.FOV.SSCVisionMap(p.Point, f.LOS, passable, false) {
			if paths.DistanceManhattan(point, p.Point) > f.LOS {
				continue
			}
			if !s.ecs.Map.Explored[point] {
				s.ecs.Map.Explored[point] = true
			}
		}
	}
}

type DeathSystem struct {
	ecs *ECS
}

func (s *DeathSystem) Update() {
	for _, e := range s.ecs.EntitiesWith(Death{}) {
		name := s.ecs.names[e].string
		s.ecs.obstructs[e] = nil   // No longer blocking.
		s.ecs.perceptions[e] = nil // No longer perceiving.
		s.ecs.ais[e] = nil         // No longer pathing.
		s.ecs.bumps[e] = nil
		s.ecs.inputs[e] = nil
		s.ecs.damages[e] = nil
		// s.ecs.healths[e] = nil
		s.ecs.deaths[e] = nil // Consume the death component.
		s.ecs.AddComponent(e, Name{name + " corpse"})
		s.ecs.AddComponent(e, Renderable{glyph: '%', fg: ColorCorpse, order: ROCorpse})
		s.ecs.AddComponent(e, Collectible{})
		s.ecs.AddComponent(e, Consumable{hp: 2})
		// Drop everything in inventory
		if s.ecs.HasComponent(e, Inventory{}) {
			for _, item := range s.ecs.inventories[e].items {
				s.ecs.AddComponent(item, Position{s.ecs.positions[e].Point})
				// r := s.ecs.renderables[item]
				// s.ecs.AddComponent(item, Renderable{r.glyph, r.fg, r.bg, ROActor})
			}
			s.ecs.inventories[e] = nil
		}
		s.ecs.Create(LogEntry{
			Text:  fmt.Sprintf("%s has died!", name),
			Color: ColorLogMonsterAttack,
		})
	}
}

type PerceptionSystem struct {
	ecs *ECS
}

func (s *PerceptionSystem) Update() {
	for _, e := range s.ecs.EntitiesWith(Position{}, Perception{}) {
		pos := s.ecs.positions[e]
		per := s.ecs.perceptions[e]
		per.perceived = []int{}
		if per.FOV == nil {
			per.FOV = rl.NewFOV(gruid.NewRange(-per.LOS, -per.LOS, per.LOS+1, per.LOS+1))
		}
		rg := gruid.NewRange(-per.LOS, -per.LOS, per.LOS+1, per.LOS+1)
		per.FOV.SetRange(rg.Add(pos.Point).Intersect(s.ecs.Map.Grid.Range()))
		passable := func(p gruid.Point) bool {
			return s.ecs.Map.Grid.At(p) != Wall
		}
		for _, point := range per.FOV.SSCVisionMap(pos.Point, per.LOS, passable, false) {
			if paths.DistanceManhattan(point, pos.Point) > per.LOS {
				continue
			}
		}
		for _, other := range s.ecs.EntitiesWith(Position{}) {
			// Ignore self.
			if other == e {
				continue
			}
			// If other entity is within perceptive radius, add to perceived list.
			pos_other := s.ecs.positions[other]
			if per.FOV.Visible(pos_other.Point) {
				per.perceived = append(per.perceived, other)
			}
		}
		// Check if player is within perceived entities, and switch to appropriate state.
		for _, other := range per.perceived {
			name_other := s.ecs.names[other].string
			if name_other == "player" && s.ecs.HasComponent(e, AI{}) {
				s.ecs.ais[e].state = CSHunting
				break
			} else {
				s.ecs.ais[e].state = CSWandering
			}
		}
	}
}

type AISystem struct {
	ecs *ECS
	aip *aiPath
}

type aiPath struct {
	ecs *ECS
	nb  paths.Neighbors
}

func (aip *aiPath) Neighbors(q gruid.Point) []gruid.Point {
	return aip.nb.All(q,
		func(r gruid.Point) bool {
			return aip.ecs.Map.Walkable(r)
		})
}

func (aip *aiPath) Cost(p, q gruid.Point) int {
	if !aip.ecs.NoBlockingEntityAt(q) {
		// Extra cost for blocked positions: this encourages the pathfinding
		// algorithm to take another path to reach the their destination.
		return 8
	}
	return 1
}

func (aip *aiPath) Estimation(p, q gruid.Point) int {
	return paths.DistanceChebyshev(p, q)
}

func (s *AISystem) Update() {
	for _, e := range s.ecs.EntitiesWith(AI{}, Position{}) {
		ai := s.ecs.ais[e]
		pos := s.ecs.positions[e]
		switch ai.state {

		case CSSleeping:
			// Do nothing, the entity is asleep!
			return

		case CSWandering:
			// Set a destination, if one is not yet set or we've reached it.
			if ai.dest == nil || *ai.dest == pos.Point {
				for {
					f := s.ecs.Map.RandomFloor()
					if f != pos.Point {
						ai.dest = &f
						break
					}
				}
			}

		case CSHunting:
			// Set destination to be the player.
			ai.dest = &s.ecs.positions[0].Point
		}
		// Compute path to ai.dest.
		path := s.ecs.Map.PR.AstarPath(&aiPath{ecs: s.ecs}, pos.Point, *ai.dest)
		q := path[1]
		// Move entity to first position in the path.
		s.ecs.AddComponent(e, Bump{q.Sub(pos.Point)})
	}
}

type DebugSystem struct {
	ecs *ECS
}

// Prints out component information for every entity.
func (s *DebugSystem) Update() {
	fmt.Println("+++++++++++++++++++++++++++++++ DEBUG ++++++++++++++++++++++++++++++++++++++++++")
	for _, e := range s.ecs.entities {
		s.ecs.printDebug(e)
	}
}
