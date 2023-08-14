package main

import (
	"github.com/anaseto/gruid"
)

type ECS struct {
	entities []int

	positions   map[int]*Position
	renderables map[int]*Renderable
	names       map[int]*Name
	inputs      map[int]*Input
	bumps       map[int]*Bump
	fovs        map[int]*FOV
	obstructs   map[int]*Obstruct

	systems []System

	drawgrid *gruid.Grid
	Map      *Map
}

// Note that we do not initialize the map here. The idea is that
// the callee is initializing that and will assign it right after this.
func NewECS() *ECS {
	ecs := &ECS{
		positions:   make(map[int]*Position),
		renderables: make(map[int]*Renderable),
		names:       make(map[int]*Name),
		inputs:      make(map[int]*Input),
		bumps:       make(map[int]*Bump),
		fovs:        make(map[int]*FOV),
		obstructs:   make(map[int]*Obstruct),
	}
	ecs.systems = append(ecs.systems, &MovementSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &FOVSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &RenderSystem{ecs: ecs})

	return ecs
}

func (ecs *ECS) Create(components ...any) int {
	idx := len(ecs.entities)
	ecs.entities = append(ecs.entities, idx)
	for _, component := range components {
		switch c := component.(type) {
		case Position:
			ecs.positions[idx] = &c
		case Renderable:
			ecs.renderables[idx] = &c
		case Name:
			ecs.names[idx] = &c
		case Input:
			ecs.inputs[idx] = &c
		case Bump:
			ecs.bumps[idx] = &c
		case FOV:
			ecs.fovs[idx] = &c
		case Obstruct:
			ecs.obstructs[idx] = &c
		}
	}
	return idx
}

func (ecs *ECS) Update() {
	for _, s := range ecs.systems {
		s.Update()
	}
}

func (ecs *ECS) Exists(entity int) bool {
	for _, e := range ecs.entities {
		if entity == e {
			return true
		}
	}
	return false
}

func (ecs *ECS) AddComponent(entity int, component any) {
	switch c := component.(type) {
	case Position:
		ecs.positions[entity] = &c
	case Renderable:
		ecs.renderables[entity] = &c
	case Name:
		ecs.names[entity] = &c
	case Input:
		ecs.inputs[entity] = &c
	case Bump:
		ecs.bumps[entity] = &c
	case FOV:
		ecs.fovs[entity] = &c
	case Obstruct:
		ecs.obstructs[entity] = &c
	}
}

func (ecs *ECS) AddComponents(entity int, components ...any) {
	for _, c := range components {
		ecs.AddComponent(entity, c)
	}
}

func (ecs *ECS) HasComponent(entity int, component any) bool {
	switch component.(type) {
	case Position:
		if c := ecs.positions[entity]; c != nil {
			return true
		}
	case Renderable:
		if c := ecs.renderables[entity]; c != nil {
			return true
		}
	case Name:
		if c := ecs.names[entity]; c != nil {
			return true
		}
	case Input:
		if c := ecs.inputs[entity]; c != nil {
			return true
		}
	case Bump:
		if c := ecs.bumps[entity]; c != nil {
			return true
		}
	case FOV:
		if c := ecs.fovs[entity]; c != nil {
			return true
		}
	case Obstruct:
		if c := ecs.obstructs[entity]; c != nil {
			return true
		}
	}
	return false
}

func (ecs *ECS) GetComponent(entity int, component any) any {
	switch component.(type) {
	case Position:
		return ecs.positions[entity]
	case Renderable:
		return ecs.renderables[entity]
	case Name:
		return ecs.names[entity]
	case Input:
		return ecs.inputs[entity]
	case Bump:
		return ecs.bumps[entity]
	case FOV:
		return ecs.fovs[entity]
	case Obstruct:
		return ecs.obstructs[entity]
	}
	return nil
}

func (ecs *ECS) GetEntityAt(p gruid.Point) (entity int, ok bool) {
	for i, pos := range ecs.positions {
		if pos != nil && p.X == pos.X && p.Y == pos.Y {
			return i, true
		}
	}
	return -1, false
}

func (ecs *ECS) RemoveComponent(entity int, component any) {
	switch component.(type) {
	case Position:
		ecs.positions[entity] = nil
	case Renderable:
		ecs.renderables[entity] = nil
	case Name:
		ecs.names[entity] = nil
	case Input:
		ecs.inputs[entity] = nil
	case Bump:
		ecs.bumps[entity] = nil
	case FOV:
		ecs.fovs[entity] = nil
	case Obstruct:
		ecs.obstructs[entity] = nil
	}
}

// Not valid Go?
// No.
// https://stackoverflow.com/questions/17934611/multiple-assignment-from-array-or-slice
// func (ecs *ECS) GetComponents(entity int, components ...any) (results []any, ok bool) {
// 	ok = true
// 	for _, c := range components {
// 		results = append(results, ecs.GetComponent(entity, c))
// 		if results[len(results)-1] == nil {
// 			ok = false
// 		}
// 	}
// 	results = append(results, ok)
// 	return results, ok
// }

// Draws all entities onto a passed grid.
func (ecs *ECS) Draw(grid *gruid.Grid) {
	for _, e := range ecs.entities {
		if ecs.HasComponent(e, Position{}) && ecs.HasComponent(e, Renderable{}) {
			p := ecs.positions[e]
			r := ecs.renderables[e]
			bg := gruid.ColorDefault
			if ecs.HasComponent(e, FOV{}) {
				bg = ColorFOV
			}
			grid.Set(p.Point, gruid.Cell{
				Rune:  r.glyph,
				Style: gruid.Style{Fg: r.color, Bg: bg},
			})
		}
	}
}
