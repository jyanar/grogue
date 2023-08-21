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
		b := s.ecs.bumps[e]
		p := s.ecs.positions[e]
		if s.ecs.Map.Walkable(p.Point.Add(b.Point)) {
			// There's an entitity at the target location.
			if target, ok := s.ecs.GetEntityAt(p.Point.Add(b.Point)); ok {
				// Attack is defined, if target has health and obstruct components.
				if s.ecs.HasComponents(target, Health{}, Obstruct{}) {
					dmg := s.ecs.damages[e].int
					name := s.ecs.names[e].string
					name_target := s.ecs.names[target].string
					health_target := s.ecs.healths[target]
					fmt.Printf("%s hits the %s for %d damage!\n", name, name_target, dmg)
					health_target.hp -= dmg
					if health_target.hp <= 0 {
						health_target.hp = 0
						s.ecs.AddComponent(target, Death{}) // Entity marked for death.
					}
					return
				}
			}
			// Otherwise, move to the location.
			p.Point = p.Point.Add(b.Point)
		}
		b = nil // Consume the bump component.
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
			if paths.DistanceChebyshev(pos.Point, pos_other.Point) < per.radius {
				per.perceived = append(per.perceived, other)
				// If the other entity is the player, switch creature state to Hunting.
				// if other == 0 && s.ecs.HasComponent(e, AI{}) {
				// 	s.ecs.ais[e].state = CSHunting
				// }
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
	fmt.Println("===================")
	for _, e := range s.ecs.EntitiesWith(AI{}, Position{}) {
		ai := s.ecs.ais[e]
		pos := s.ecs.positions[e]
		switch ai.state {

		case CSSleeping:
			return // Do nothing, creature is asleep!

		case CSWandering:
			// Set a destination, if one is not yet set.
			if ai.dest == nil || *ai.dest == pos.Point {
				for {
					f := s.ecs.Map.RandomFloor()
					if f != pos.Point {
						ai.dest = &f
						break
					}
				}
				fmt.Printf("New path established for entity: %d\n", e)
				fmt.Printf("Waypoint: %v\n\n", ai.dest)
			}
			// Compute path to it.
			path := s.ecs.Map.PR.AstarPath(&aiPath{ecs: s.ecs}, pos.Point, *ai.dest)
			// Move entity to first position in the path.
			s.ecs.positions[e] = &Position{path[1]}

			// var bump gruid.Point
			// bump = pos.Point.Sub(path[1])
			// s.ecs.AddComponent(e, Bump{bump})

			// // If current path is run out, start pathing towards new direction.
			// // if len(ai.path) < 1 {
			// // 	ai.path = s.ecs.Map.PR.AstarPath(&aiPath{ecs: s.ecs}, pos.Point, s.ecs.Map.RandomFloor())
			// // }
			// ai.path = s.ecs.Map.PR.AstarPath(&aiPath{ecs: s.ecs}, pos.Point, s.ecs.Map.RandomFloor())
			// bump := pos.Point.Sub(ai.path[1])
			// fmt.Printf("Current LOC: %v\n", pos.Point)
			// fmt.Printf("ai.path[1]:  %v\n", ai.path[1])
			// fmt.Printf("BUMP: %v, %T\n", bump, bump)
			// // fmt.Printf("path[2:5]: %v\n", ai.path[2:5])
			// s.ecs.positions[e] = &Position{ai.path[1]}
			// // s.ecs.AddComponent(e, Bump{pos.Point.Sub(ai.path[1])})
		case CSHunting:
			// Compute path to player.
			// player_pos := s.ecs.positions[0]
			// ai.path = s.ecs.Map.PR.AstarPath(&aiPath{ecs: s.ecs}, pos.Point, player_pos.Point)
			// s.ecs.AddComponent(e, Bump{pos.Point.Sub(ai.path[1])})
		}
	}
}

type DebugSystem struct {
	ecs *ECS
}

// Prints out component information for every entity.
func (s *DebugSystem) Update() {
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	for _, e := range s.ecs.entities {
		s.ecs.printDebug(e)
	}
}
