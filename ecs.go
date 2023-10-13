package main

import (
	"fmt"

	"github.com/anaseto/gruid"
)

type ECS struct {
	entities []int
	nextID   int

	positions    map[int]*Position
	renderables  map[int]*Renderable
	names        map[int]*Name
	inputs       map[int]*Input
	bumps        map[int]*Bump
	fovs         map[int]*FOV
	obstructs    map[int]*Obstruct
	healths      map[int]*Health
	damages      map[int]*Damage
	deaths       map[int]*Death
	perceptions  map[int]*Perception
	ais          map[int]*AI
	logentries   map[int]*LogEntry
	consumables  map[int]*Consumable
	collectibles map[int]*Collectible
	inventories  map[int]*Inventory

	systems []System

	perceptionSystem PerceptionSystem
	aISystem         AISystem
	bumpSystem       BumpSystem
	fOVSystem        FOVSystem
	deathSystem      DeathSystem

	Map *Map
}

// Note that we do not initialize the map here. The idea is that
// the callee is initializing that and will assign it right after this.
func NewECS() *ECS {
	ecs := &ECS{
		nextID:       0,
		positions:    make(map[int]*Position),
		renderables:  make(map[int]*Renderable),
		names:        make(map[int]*Name),
		inputs:       make(map[int]*Input),
		bumps:        make(map[int]*Bump),
		fovs:         make(map[int]*FOV),
		obstructs:    make(map[int]*Obstruct),
		healths:      make(map[int]*Health),
		damages:      make(map[int]*Damage),
		deaths:       make(map[int]*Death),
		perceptions:  make(map[int]*Perception),
		ais:          make(map[int]*AI),
		logentries:   make(map[int]*LogEntry),
		consumables:  make(map[int]*Consumable),
		collectibles: make(map[int]*Collectible),
		inventories:  make(map[int]*Inventory),
	}
	ecs.perceptionSystem = PerceptionSystem{ecs: ecs}
	ecs.aISystem = AISystem{ecs: ecs, aip: &aiPath{ecs: ecs}}
	ecs.bumpSystem = BumpSystem{ecs: ecs}
	ecs.fOVSystem = FOVSystem{ecs: ecs}
	ecs.deathSystem = DeathSystem{ecs: ecs}
	ecs.systems = append(ecs.systems, &PerceptionSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &AISystem{ecs: ecs, aip: &aiPath{ecs: ecs}})
	ecs.systems = append(ecs.systems, &BumpSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &FOVSystem{ecs: ecs})
	ecs.systems = append(ecs.systems, &DeathSystem{ecs: ecs})
	// ecs.systems = append(ecs.systems, &DebugSystem{ecs: ecs})

	return ecs
}

// Entity-first updating.
func (ecs *ECS) Update2() {
	// Iterate over each entity, and update them based on what components they have.
	for _, e := range ecs.entities {
		// How do we check which systems this entity applies to? It has to be done
		// dynamically, to allow us to remove and add components to entities.
		if ecs.HasComponents(e, Position{}, Perception{}) {
			ecs.perceptionSystem.Update2(e)
		}
		if ecs.HasComponents(e, Position{}, AI{}) {
			ecs.aISystem.Update2(e)
		}
		if ecs.HasComponents(e, Position{}, Bump{}) {
			ecs.bumpSystem.Update2(e)
		}
		if ecs.HasComponents(e, Position{}, FOV{}) {
			ecs.fOVSystem.Update2(e)
		}
		if ecs.HasComponents(e, Death{}) {
			ecs.deathSystem.Update2(e)
		}
	}
}

// Systems-first updating.
func (ecs *ECS) Update() {
	for _, s := range ecs.systems {
		s.Update()
	}
}

func (ecs *ECS) Create(components ...any) int {
	idx := ecs.nextID
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
		case Perception:
			ecs.perceptions[idx] = &c
		case AI:
			ecs.ais[idx] = &c
		case LogEntry:
			ecs.logentries[idx] = &c
		case Consumable:
			ecs.consumables[idx] = &c
		case Collectible:
			ecs.collectibles[idx] = &c
		case Inventory:
			ecs.inventories[idx] = &c
		}
	}
	ecs.nextID += 1
	return idx
}

func remove(slice []int, s int) []int {
	idx := -1
	for i := 0; i < len(slice); i++ {
		if slice[i] == s {
			idx = i
			break
		}
	}
	if idx != -1 {
		return removeAt(slice, idx)
	}
	return slice
}

func removeAt(slice []int, idx int) []int {
	return append(slice[:idx], slice[idx+1:]...)
}

func (ecs *ECS) Delete(entity int) {
	// Remove from entity list
	ecs.entities = remove(ecs.entities, entity)
	// Delete associated data from maps
	delete(ecs.positions, entity)
	delete(ecs.renderables, entity)
	delete(ecs.names, entity)
	delete(ecs.inputs, entity)
	delete(ecs.bumps, entity)
	delete(ecs.fovs, entity)
	delete(ecs.obstructs, entity)
	delete(ecs.healths, entity)
	delete(ecs.damages, entity)
	delete(ecs.deaths, entity)
	delete(ecs.perceptions, entity)
	delete(ecs.ais, entity)
	delete(ecs.logentries, entity)
	delete(ecs.consumables, entity)
	delete(ecs.collectibles, entity)
	delete(ecs.inventories, entity)
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
	case Perception:
		ecs.perceptions[entity] = &c
	case AI:
		ecs.ais[entity] = &c
	case LogEntry:
		ecs.logentries[entity] = &c
	case Consumable:
		ecs.consumables[entity] = &c
	case Collectible:
		ecs.collectibles[entity] = &c
	case Inventory:
		ecs.inventories[entity] = &c
	}
}

func (ecs *ECS) AddComponents(entity int, components ...any) {
	for _, c := range components {
		ecs.AddComponent(entity, c)
	}
}

func (ecs *ECS) HasComponent(entity int, component any) bool {
	switch component.(type) {
	case Name:
		if ecs.names[entity] != nil {
			return true
		}
	case Position:
		if ecs.positions[entity] != nil {
			return true
		}
	case Renderable:
		if ecs.renderables[entity] != nil {
			return true
		}
	case Input:
		if ecs.inputs[entity] != nil {
			return true
		}
	case Bump:
		if ecs.bumps[entity] != nil {
			return true
		}
	case FOV:
		if ecs.fovs[entity] != nil {
			return true
		}
	case Obstruct:
		if ecs.obstructs[entity] != nil {
			return true
		}
	case Health:
		if ecs.healths[entity] != nil {
			return true
		}
	case Damage:
		if ecs.damages[entity] != nil {
			return true
		}
	case Death:
		if ecs.deaths[entity] != nil {
			return true
		}
	case Perception:
		if ecs.perceptions[entity] != nil {
			return true
		}
	case AI:
		if ecs.ais[entity] != nil {
			return true
		}
	case LogEntry:
		if ecs.logentries[entity] != nil {
			return true
		}
	case Consumable:
		if ecs.consumables[entity] != nil {
			return true
		}
	case Collectible:
		if ecs.collectibles[entity] != nil {
			return true
		}
	case Inventory:
		if ecs.inventories[entity] != nil {
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

func (ecs *ECS) EntitiesWith(components ...any) (entities []int) {
	for _, e := range ecs.entities {
		if ecs.HasComponents(e, components...) {
			entities = append(entities, e)
		}
	}
	return entities
}

func (ecs *ECS) EntitiesAt(p gruid.Point) (entities []int) {
	for _, e := range ecs.EntitiesWith(Position{}) {
		q := ecs.positions[e].Point
		if p == q {
			entities = append(entities, e)
		}
	}
	return entities
}

func (ecs *ECS) EntitiesAtPWith(p gruid.Point, components ...any) (entities []int) {
	for _, e := range ecs.EntitiesAt(p) {
		if ecs.HasComponents(e, components...) {
			entities = append(entities, e)
		}
	}
	return entities
}

// Returns true if there is no blocking entity at p.
func (ecs *ECS) NoBlockingEntityAt(p gruid.Point) bool {
	return len(ecs.EntitiesAtPWith(p, Obstruct{})) == 0
}

func (ecs *ECS) BloodAt(p gruid.Point) bool {
	entities := ecs.EntitiesAtPWith(p, Name{})
	if len(entities) == 0 {
		return false
	}
	for _, e := range entities {
		name := ecs.names[e].string
		if name == "blood" {
			return true
		}
	}
	return false
}

func (ecs *ECS) PlayerDead() bool {
	return ecs.names[0].string == "player corpse"
}

func (ecs *ECS) printDebug(e int) {
	fmt.Println("====================")
	fmt.Printf("Entity: %d\n", e)
	if ecs.bumps[e] != nil {
		fmt.Printf("%v, %T\n", ecs.bumps[e], ecs.bumps[e])
	}
	if ecs.positions[e] != nil {
		fmt.Printf("%v, %T\n", ecs.positions[e], ecs.positions[e])
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
	if ecs.deaths[e] != nil {
		fmt.Printf("%v, %T\n", ecs.deaths[e], ecs.deaths[e])
	}
	if ecs.perceptions[e] != nil {
		fmt.Printf("%v, %T\n", ecs.perceptions[e], ecs.perceptions[e])
	}
	if ecs.ais[e] != nil {
		fmt.Printf("%v, %T\n", ecs.ais[e], ecs.ais[e])
	}
	if ecs.logentries[e] != nil {
		fmt.Printf("%v, %T\n", ecs.logentries[e], ecs.logentries[e])
	}
	if ecs.consumables[e] != nil {
		fmt.Printf("%v, %T\n", ecs.consumables[e], ecs.consumables[e])
	}
	if ecs.collectibles[e] != nil {
		fmt.Printf("%v, %T\n", ecs.collectibles[e], ecs.collectibles[e])
	}
	if ecs.inventories[e] != nil {
		fmt.Printf("%v, %T\n", ecs.inventories[e], ecs.inventories[e])
	}
}
