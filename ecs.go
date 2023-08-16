package main

import (
	"fmt"

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
	healths     map[int]*Health
	damages     map[int]*Damage
	deaths      map[int]*Death

	systems []System

	Map *Map
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
		healths:     make(map[int]*Health),
		damages:     make(map[int]*Damage),
		deaths:      make(map[int]*Death),
	}
	ecs.systems = append(ecs.systems, &BumpSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &FOVSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &DeathSystem{ecs: ecs})

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
		case Health:
			ecs.healths[idx] = &c
		case Damage:
			ecs.damages[idx] = &c
		case Death:
			ecs.deaths[idx] = &c
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

// Adds a component to an entity. If one of this type already exists,
// replaces it.
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
	case Health:
		ecs.healths[entity] = &c
	case Damage:
		ecs.damages[entity] = &c
	case Death:
		ecs.deaths[entity] = &c
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
	case Health:
		if c := ecs.healths[entity]; c != nil {
			return true
		}
	case Damage:
		if c := ecs.damages[entity]; c != nil {
			return true
		}
	case Death:
		if c := ecs.deaths[entity]; c != nil {
			return true
		}
	}
	return false
}

func (ecs *ECS) HasComponents(entity int, components ...any) bool {
	for _, c := range components {
		if !ecs.HasComponent(entity, c) {
			return false
		}
	}
	return true
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
	case Health:
		return ecs.healths[entity]
	case Damage:
		return ecs.damages[entity]
	case Death:
		return ecs.deaths[entity]
	}
	return nil
}

func (ecs *ECS) GetEntityAt(p gruid.Point) (entity int, ok bool) {
	for i, q := range ecs.positions {
		if q != nil && p == q.Point {
			return i, true
		}
	}
	return -1, false
}

// Returns true if there is no blocking entity at p.
func (ecs *ECS) NoBlockingEntityAt(p gruid.Point) bool {
	if e, ok := ecs.GetEntityAt(p); ok {
		if ecs.HasComponent(e, Obstruct{}) {
			return false
		}
	}
	return true
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
	case Health:
		ecs.obstructs[entity] = nil
	case Damage:
		ecs.damages[entity] = nil
	case Death:
		ecs.deaths[entity] = nil
	}
}

func (ecs *ECS) printDebug(e int) {
	fmt.Println("====================")
	fmt.Println("Entity: " + string(e))
	if ecs.bumps[e] != nil {
		fmt.Printf("%v, %T\n", ecs.bumps[e], ecs.bumps[e])
	}
	if ecs.damages[e] != nil {
		fmt.Printf("%v, %T\n", ecs.damages[e], ecs.damages[e])
	}
	if ecs.deaths[e] != nil {
		fmt.Printf("%v, %T\n", ecs.deaths[e], ecs.deaths[e])
	}
	if ecs.fovs[e] != nil {
		fmt.Printf("%v, %T\n", ecs.fovs[e], ecs.fovs[e])
	}
	if ecs.healths[e] != nil {
		fmt.Printf("%v, %T\n", ecs.healths[e], ecs.healths[e])
	}
	if ecs.inputs[e] != nil {
		fmt.Printf("%v, %T\n", ecs.inputs[e], ecs.inputs[e])
	}
	if ecs.names[e] != nil {
		fmt.Printf("%v, %T\n", ecs.names[e], ecs.names[e])
	}
	if ecs.obstructs[e] != nil {
		fmt.Printf("%v, %T\n", ecs.obstructs[e], ecs.obstructs[e])
	}
	if ecs.renderables[e] != nil {
		fmt.Printf("%v, %T\n", ecs.renderables[e], ecs.renderables[e])
	}
	if ecs.positions[e] != nil {
		fmt.Printf("%v, %T\n", ecs.positions[e], ecs.positions[e])
	}

}
