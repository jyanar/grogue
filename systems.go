package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
)

type System interface {
	Update()
}

type MovementSystem struct {
	ecs *ECS
}

func (s *MovementSystem) Update() {
	for _, e := range s.ecs.entities {
		if s.ecs.HasComponent(e, Bump{}) && s.ecs.HasComponent(e, Position{}) {
			p := s.ecs.positions[e]
			b := s.ecs.bumps[e]
			if s.ecs.Map.Walkable(p.Point.Add(b.Point)) {
				if _, ok := s.ecs.GetEntityAt(p.Point.Add(b.Point)); ok {
					fmt.Println("Entity there! Can't move!")
				} else {
					p.Point = p.Point.Add(b.Point)
				}
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

type RenderSystem struct {
	ecs *ECS
}

// Note, if we want the RenderSystem here to handle everything (map drawing, entity
// drawing, etc) then we need to pass a pointer of the model to the ECS system. Which
// ends up making the code a bit less clean and more ugly. However, it does end up
// with the benefit of having all the logic being isolated to one place -- here, in the
// systems. The only exception really is the Map, which has its own set of logics for
// generating the map and everything.
//
// It might just be better to keep the rendersystem separate, in model. That way we can
// adhere to the Model View Controller architecture cleanly.
func (s *RenderSystem) Update() {
	s.ecs.drawgrid.Fill(gruid.Cell{Rune: ' '})
	// Draw the map first.
	s.ecs.Map.Draw(s.ecs.drawgrid)
	// Draw the entities.
	for _, e := range s.ecs.entities {
		if s.ecs.HasComponent(e, Position{}) && s.ecs.HasComponent(e, Renderable{}) {
			p := s.ecs.positions[e]
			r := s.ecs.renderables[e]
			bg := gruid.ColorDefault
			if s.ecs.HasComponent(e, FOV{}) {
				bg = ColorFOV
			}
			s.ecs.drawgrid.Set(p.Point, gruid.Cell{
				Rune:  r.glyph,
				Style: gruid.Style{Fg: r.color, Bg: bg},
			})
		}
	}
}
