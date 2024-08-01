package main

import (
	"fmt"

	"github.com/anaseto/gruid"
)

type ECS struct {
	entities []int
	nextID   int
	Map      *Map
	// Components
	components map[int]map[string]Component
	// positions     map[int]*Position
	// renderables   map[int]*Renderable
	// names         map[int]*Name
	// inputs        map[int]*Input
	// bumps         map[int]*Bump
	// fovs          map[int]*FOV
	// obstructs     map[int]*Obstruct
	// healths       map[int]*Health
	// damages       map[int]*Damage
	// deaths        map[int]*Death
	// perceptions   map[int]*Perception
	// visibles      map[int]*Visible
	// ais           map[int]*AI
	// logentries    map[int]*LogEntry
	// consumables   map[int]*Consumable
	// healings      map[int]*Healing
	// collectibles  map[int]*Collectible
	// inventories   map[int]*Inventory
	// rangeds       map[int]*Ranged
	// damageeffects map[int][]DamageEffect
	// animations    map[int]*Animation
	// Systems
	PerceptionSystem
	AISystem
	BumpSystem
	FOVSystem
	DeathSystem
	DamageEffectSystem
	DebugSystem
	AnimationSystem
}

// Note that we do not initialize the map here. The idea is that
// the callee is initializing that and will assign it right after this.
func NewECS() *ECS {
	ecs := &ECS{
		nextID:     0,
		components: make(map[int]map[string]Component),
		// positions:     make(map[int]*Position),
		// renderables:   make(map[int]*Renderable),
		// names:         make(map[int]*Name),
		// inputs:        make(map[int]*Input),
		// bumps:         make(map[int]*Bump),
		// fovs:          make(map[int]*FOV),
		// obstructs:     make(map[int]*Obstruct),
		// healths:       make(map[int]*Health),
		// damages:       make(map[int]*Damage),
		// deaths:        make(map[int]*Death),
		// perceptions:   make(map[int]*Perception),
		// visibles:      make(map[int]*Visible),
		// ais:           make(map[int]*AI),
		// logentries:    make(map[int]*LogEntry),
		// consumables:   make(map[int]*Consumable),
		// healings:      make(map[int]*Healing),
		// collectibles:  make(map[int]*Collectible),
		// inventories:   make(map[int]*Inventory),
		// rangeds:       make(map[int]*Ranged),
		// damageeffects: make(map[int][]DamageEffect),
		// animations:    make(map[int]*Animation),
	}
	ecs.PerceptionSystem = PerceptionSystem{ecs: ecs}
	ecs.AISystem = AISystem{ecs: ecs, aip: &aiPath{ecs: ecs}}
	ecs.BumpSystem = BumpSystem{ecs: ecs}
	ecs.FOVSystem = FOVSystem{ecs: ecs}
	ecs.DeathSystem = DeathSystem{ecs: ecs}
	ecs.DamageEffectSystem = DamageEffectSystem{ecs: ecs}
	ecs.AnimationSystem = AnimationSystem{ecs: ecs}
	ecs.DebugSystem = DebugSystem{ecs: ecs}
	return ecs
}

func (ecs *ECS) Initialize() {
	for _, e := range ecs.entities {
		ecs.PerceptionSystem.Update(e)
		ecs.AISystem.Update(e)
		ecs.FOVSystem.Update(e)
	}
}

// Iterates through each entity
func (ecs *ECS) Update() {
	for _, e := range ecs.entities {
		ecs.DamageEffectSystem.Update(e)
		ecs.PerceptionSystem.Update(e)
		ecs.AISystem.Update(e)
		ecs.BumpSystem.Update(e)
		ecs.FOVSystem.Update(e)
		ecs.DamageEffectSystem.Update(e)
		ecs.DeathSystem.Update(e)
	}
	// ecs.DebugSystem.Update()
}

func (ecs *ECS) UpdateAnimation() {
	for _, e := range ecs.EntitiesWith(Animation{}) {
		ecs.AnimationSystem.Update(e)
	}
}

func (ecs *ECS) Create(components ...any) int {
	idx := ecs.nextID
	ecs.entities = append(ecs.entities, idx)
	for _, component := range components {
		s.ecs.AddComponent(idx, component)
	}
	ecs.nextID++
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
	delete(ecs.components, entity)
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
	if componentMap, ok := ecs.components[entity]; ok {
		componentString := fmt.Sprintf("%T", component)
		componentMap[componentString] = component
	} else {
		ecs.components[entity] = make(map[string]Component)
		ecs.AddComponent(entity, component)
	}
}

func (ecs *ECS) AddComponents(entity int, components ...any) {
	for _, c := range components {
		ecs.AddComponent(entity, c)
	}
}

func (ecs *ECS) GetComponent(entity int, component Component) (Component, bool) {
	if _, ok := ecs.components[entity]; ok {
		componentString := fmt.Sprintf("%T", component)
		if component, exists := ecs.components[entity][componentString]; exists {
			return component, true
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (ecs *ECS) HasComponent(entity int, component any) bool {
	if _, ok := ecs.GetComponent(entity, component); ok {
		return true
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
		p, _ := ecs.GetComponent(e, Position{})
		pos := p.(Position)
		if pos == p {
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
		// name := ecs.names[e].string
		// if name == "blood" {
		// 	return true
		// }
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
	if ecs.visibles[e] != nil {
		fmt.Printf("%v, %T\n", ecs.visibles[e], ecs.visibles[e])
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
	if ecs.healings[e] != nil {
		fmt.Printf("%v, %T\n", ecs.healings[e], ecs.healings[e])
	}
	if ecs.collectibles[e] != nil {
		fmt.Printf("%v, %T\n", ecs.collectibles[e], ecs.collectibles[e])
	}
	if ecs.inventories[e] != nil {
		fmt.Printf("%v, %T\n", ecs.inventories[e], ecs.inventories[e])
	}
	if ecs.rangeds[e] != nil {
		fmt.Printf("%v, %T\n", ecs.rangeds[e], ecs.rangeds[e])
	}
	if len(ecs.damageeffects[e]) > 0 {
		for _, de := range ecs.damageeffects[e] {
			fmt.Printf("%v, %T\n", de, de)
		}
	}
	if ecs.animations[e] != nil {
		fmt.Printf("%v, %T\n", ecs.animations[e], ecs.animations[e])
	}
}
