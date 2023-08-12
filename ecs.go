package main

import (
	"fmt"
)

type ECS struct {
	entities []int

	positions   map[int]*Position
	renderables map[int]*Renderable
	names       map[int]*Name
	inputs      map[int]*Input
	bumps       map[int]*Bump
	fovs        map[int]*FOV

	systems []System

	Map *Map
}

func NewECS() *ECS {
	ecs := &ECS{}
	ecs.positions = make(map[int]*Position)
	ecs.renderables = make(map[int]*Renderable)
	ecs.names = make(map[int]*Name)
	ecs.inputs = make(map[int]*Input)
	ecs.bumps = make(map[int]*Bump)
	ecs.fovs = make(map[int]*FOV)
	ecs.systems = append(ecs.systems, &MovementSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &FOVSystem{ecs: ecs})
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
	}
	return false
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
	}
}

func (ecs *ECS) printData(entity int) {
	fmt.Printf("Information for %v...\n", entity)
	if v, ok := ecs.positions[entity]; ok {
		fmt.Printf("%v, %T\n", v, v)
	}
	if v, ok := ecs.renderables[entity]; ok {
		fmt.Printf("%v, %T\n", v, v)
	}
	if v, ok := ecs.names[entity]; ok {
		fmt.Printf("%v, %T\n", v, v)
	}
}

// func main() {
// 	// ecs := ECS{}
// 	ecs := NewECS()
// 	ecs.Create(Position{10, 10}, Renderable{glyph: '@'})
// 	ecs.Create(Renderable{glyph: 'g'})

// 	fmt.Println("\nPrinting data...")
// 	ecs.printData(0)
// 	ecs.printData(1)

// 	ecs.AddComponent(1, Position{15, 10})
// 	fmt.Println("\nPrinting again!")
// 	ecs.printData(0)
// 	ecs.printData(1)

// 	ecs.RemoveComponent(0, Renderable{})
// 	fmt.Println("\nPrinting again!")
// 	ecs.printData(0)
// 	ecs.printData(1)

// }
