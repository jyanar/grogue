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
	for _, e := range s.ecs.entities {
		if s.ecs.HasComponent(e, Bump{}) && s.ecs.HasComponent(e, Position{}) {
			p := s.ecs.positions[e]
			b := s.ecs.bumps[e]
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
}

type FOVSystem struct {
	ecs *ECS
}

func (s *FOVSystem) Update() {
	for _, e := range s.ecs.entities {
		if s.ecs.HasComponent(e, FOV{}) && s.ecs.HasComponent(e, Position{}) {
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
}

type DeathSystem struct {
	ecs *ECS
}

func (s *DeathSystem) Update() {
	for _, e := range s.ecs.entities {
		if s.ecs.HasComponent(e, Death{}) {
			name := s.ecs.names[e]
			s.ecs.obstructs[e] = nil // No longer blocking.
			s.ecs.AddComponent(e, Name{"Remains of " + name.string})
			s.ecs.AddComponent(e, Renderable{glyph: '%', color: ColorCorpse, order: ROCorpse})
		}
		s.ecs.deaths[e] = nil // Consume the death component.
	}
}
