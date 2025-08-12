package main

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"github.com/k0kubun/pp/v3"
)

type ECS struct {
	entities   []int
	nextID     int
	Map        *Map
	components map[int]map[string]Component
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
		ecs.PerceptionSystem.Update(e)
		ecs.AISystem.Update(e)
		ecs.BumpSystem.Update(e)
		ecs.FOVSystem.Update(e)
	}
	for _, e := range ecs.entities {
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
		ecs.AddComponent(idx, component)
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

func (ecs *ECS) GetComponentUnchecked(entity int, component Component) Component {
	componentString := fmt.Sprintf("%T", component)
	return ecs.components[entity][componentString]
}

func (ecs *ECS) GetComponentsFor(entity int) map[string]Component {
	return ecs.components[entity]
}

func (ecs *ECS) RemoveComponent(entity int, component Component) {
	if _, ok := ecs.components[entity]; ok {
		componentString := fmt.Sprintf("%T", component)
		delete(ecs.components[entity], componentString)
	}
}

func (ecs *ECS) ClearAllComponents(entity int) {
	if _, ok := ecs.components[entity]; ok {
		ecs.components[entity] = make(map[string]Component)
	}
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
		ep, _ := ecs.GetComponent(e, Position{})
		pos := ep.(Position).Point
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
		n, hasName := ecs.GetComponent(e, Name{})
		if hasName && n.(Name).string == "blood" {
			return true
		}
	}
	return false
}

func (ecs *ECS) PlayerDead() bool {
	if _, ok := ecs.GetComponent(0, Death{}); ok {
		return true
	}
	return false
}

func (ecs *ECS) printDebug(e int) {
	fmt.Printf("Entity: %d\n", e)
	pp.Print(ecs.GetComponentsFor(e))
}
