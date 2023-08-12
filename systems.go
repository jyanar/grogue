package main

import (
	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type System interface {
	Update()
}

type MovementSystem struct {
	ecs *ECS
}

func (ms *MovementSystem) Update() {
	for _, e := range ms.ecs.entities {
		if ms.ecs.HasComponent(e, Bump{}) && ms.ecs.HasComponent(e, Position{}) {
			p := ms.ecs.positions[e]
			b := ms.ecs.bumps[e]
			if ms.ecs.Map.Walkable(p.Point.Add(b.Point)) {
				p.Point = p.Point.Add(b.Point)
			}
			b = nil
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
