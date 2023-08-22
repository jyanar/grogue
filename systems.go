package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
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
			return
		}
		// Let's attempt to move to dest.
		if s.ecs.Map.Walkable(dest) {
			// There's another entity at the location.
			// If attack is defined between this pair, perform attack.
			if target, ok := s.ecs.GetEntityAt(dest); ok {
				if s.ecs.HasComponents(target, Health{}, Obstruct{}) {
					dmg_src := s.ecs.damages[e].int
					name_src := s.ecs.names[e].string
					name_target := s.ecs.names[target].string
					health_target := s.ecs.healths[target]
					fmt.Printf("%s hits the %s for %d damage!\n", name_src, name_target, dmg_src)
					health_target.hp -= dmg_src
					if health_target.hp <= 0 {
						health_target.hp = 0
						s.ecs.AddComponent(target, Death{}) // Entity marked for death.
					}
					return
				}
			}
			// Otherwise, move to destination.
			p.Point = dest
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
		name := s.ecs.names[e]
		s.ecs.obstructs[e] = nil   // No longer blocking.
		s.ecs.perceptions[e] = nil // No longer perceiving.
		s.ecs.ais[e] = nil         // No longer pathing.
		s.ecs.bumps[e] = nil
		s.ecs.inputs[e] = nil
		s.ecs.damages[e] = nil
		s.ecs.healths[e] = nil
		s.ecs.deaths[e] = nil // Consume the death component.
		s.ecs.AddComponent(e, Name{"Remains of " + name.string})
		s.ecs.AddComponent(e, Renderable{glyph: '%', color: ColorCorpse, order: ROCorpse})
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
		for _, other := range s.ecs.EntitiesWith(Position{}) {
			// Ignore self.
			if other == e {
				continue
			}
			// If other entity is within perceptive radius, add to perceived list.
			pos_other := s.ecs.positions[other]
			if paths.DistanceChebyshev(pos.Point, pos_other.Point) <= per.radius {
				per.perceived = append(per.perceived, other)
				// If the other entity is the player, switch creature state to Hunting.
				if other == 0 && s.ecs.HasComponent(e, AI{}) {
					s.ecs.ais[e].state = CSHunting
				}
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
	return aip.nb.Cardinal(q,
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
	return paths.DistanceManhattan(p, q)
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
