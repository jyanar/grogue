package main

import (
	"fmt"

	"codeberg.org/anaseto/gruid"
	"codeberg.org/anaseto/gruid/paths"
	"codeberg.org/anaseto/gruid/rl"
)

type System interface {
	Update()
}

type PerceptionSystem struct {
	ecs *ECS
}

// Perception - allows entities with Perception{} and Position{} to perceive
// other entities within their field of view. If the given entity has an AI
// component (is a mob) and the player is within its field of view, it will
// switch to the hunting state.
func (s *PerceptionSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Position{}, Perception{}) {
		return
	}
	pos := s.ecs.GetComponentUnchecked(e, Position{}).(Position)
	per := s.ecs.GetComponentUnchecked(e, Perception{}).(Perception)
	if per.FOV == nil {
		per.FOV = rl.NewFOV(gruid.NewRange(-per.LOS, -per.LOS, per.LOS+1, per.LOS+1))
	}
	rg := gruid.NewRange(-per.LOS, -per.LOS, per.LOS+1, per.LOS+1)
	per.FOV.SetRange(rg.Add(pos.Point).Intersect(s.ecs.Map.Grid.Range()))
	passable := func(p gruid.Point) bool {
		return s.ecs.Map.Grid.At(p) != Wall
	}
	for _, point := range per.FOV.SSCVisionMap(pos.Point, per.LOS, passable, true) {
		if paths.DistanceChebyshev(point, pos.Point) > per.LOS {
			continue
		}
	}
	for _, other := range s.ecs.EntitiesWith(Position{}, Visible{}) {
		// Ignore self.
		if other == e {
			continue
		}
		// If other entity is within perceptive radius, add to perceived list.
		pos_other := s.ecs.GetComponentUnchecked(other, Position{}).(Position)
		if per.FOV.Visible(pos_other.Point) {
			per.perceived = append(per.perceived, other)
		}
	}
	// If we're a mob and the player is perceived, switch to hunting state.
	if e != 0 && s.ecs.HasComponent(e, AI{}) {
		player_found := false
		for _, other := range per.perceived {
			if other == 0 {
				player_found = true
				break
			}
		}
		ai := s.ecs.GetComponentUnchecked(e, AI{}).(AI)
		if player_found {
			ai.state = CSHunting
		} else {
			ai.state = CSWandering
		}
		s.ecs.AddComponent(e, ai)
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
	return aip.nb.All(q,
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
	return paths.DistanceChebyshev(p, q)
}

func (s *AISystem) Update(e int) {
	if !s.ecs.HasComponents(e, Position{}, AI{}) {
		return
	}
	ai := s.ecs.GetComponentUnchecked(e, AI{}).(AI)
	pos := s.ecs.GetComponentUnchecked(e, Position{}).(Position)
	switch ai.state {
	case CSSleeping:
		// Do nothing, the entity is asleep!
		return
	case CSWandering:
		// Set a destination, if one is not yet set or we've reached it.
		if ai.dest == nil || *ai.dest == pos.Point {
			for {
				f := s.ecs.Map.RandomFloor()
				if f != pos.Point {
					ai.dest = &f
					s.ecs.AddComponent(e, ai)
					break
				}
			}
		}
	case CSHunting:
		// Set destination to be the player.
		pp := s.ecs.GetComponentUnchecked(0, Position{}).(Position)
		ai.dest = &pp.Point
		s.ecs.AddComponent(e, ai)
	}
	// Compute path to ai.dest.
	path := s.ecs.Map.PR.AstarPath(&aiPath{ecs: s.ecs}, pos.Point, *ai.dest)
	q := path[1]
	// Move entity to first position in the path.
	s.ecs.AddComponent(e, Bump{q.Sub(pos.Point)})
}

type BumpSystem struct {
	ecs *ECS
}

func (s *BumpSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Bump{}, Position{}) {
		return
	}
	b := s.ecs.GetComponentUnchecked(e, Bump{}).(Bump)
	p := s.ecs.GetComponentUnchecked(e, Position{}).(Position)
	dest := p.Point.Add(b.Point)
	s.ecs.RemoveComponent(e, Bump{})
	// Ignore movement to the same tile.
	if b.X == 0 && b.Y == 0 {
		return
	}
	// Let's attempt to move to dest.
	if s.ecs.Map.Walkable(dest) {
		// Check whether there are blocking entities at dest.
		attackable_entities := s.ecs.EntitiesAtPWith(dest, Health{}, Obstruct{})
		if len(attackable_entities) == 0 {
			p.Point = dest // No entity blocking the way, move to dest.
			s.ecs.AddComponent(e, p)
			return
		}
		if len(attackable_entities) > 1 {
			panic(fmt.Sprintf("More than one entity with obstruct at position: %v", dest))
		}
		// Attack entity at location.
		dmg := s.ecs.GetComponentUnchecked(e, Damage{}).(Damage)
		if !s.ecs.HasComponent(attackable_entities[0], DamageEffects{}) {
			s.ecs.AddComponent(attackable_entities[0], DamageEffects{effects: []DamageEffect{}})
		}
		dmgfx := s.ecs.GetComponentUnchecked(attackable_entities[0], DamageEffects{}).(DamageEffects)
		target_entity := attackable_entities[0]
		// Add damage effect to the target entity
		dmgfx.effects = append(dmgfx.effects, DamageEffect{source: e, amount: dmg.int})
		s.ecs.AddComponent(target_entity, dmgfx)
		s.ecs.AddComponent(e, p)

		// s.ecs.AddComponent(target_entity, DamageEffect{source: e, amount: attack_power})
		// s.ecs.DamageEffectSystem.Update(target_entity)
	} else {
		s.ecs.Create(LogEntry{Text: "The wall is firm and unyielding!", Color: ColorLogSpecial})
	}
}

// Processes damage effects on entities with health components.
type DamageEffectSystem struct {
	ecs *ECS
}

func (s *DamageEffectSystem) Update(e int) {
	if !s.ecs.HasComponents(e, DamageEffects{}, Health{}) {
		return
	}
	health := s.ecs.GetComponentUnchecked(e, Health{}).(Health)
	dmgfx := s.ecs.GetComponentUnchecked(e, DamageEffects{}).(DamageEffects)
	for _, de := range dmgfx.effects {
		health.hp -= de.amount
		name_attacker := s.ecs.GetComponentUnchecked(de.source, Name{}).(Name).string
		name_receiver := s.ecs.GetComponentUnchecked(e, Name{}).(Name).string
		var msg string
		if de.source == 0 {
			msg = fmt.Sprintf("You stab the %s with your sword!", name_receiver)
		} else {
			if e == 0 {
				msg = fmt.Sprintf("The %s mauls you!", name_attacker)
			} else {
				msg = fmt.Sprintf("The %s hits the %s.", name_attacker, name_receiver)
			}
		}
		msgcolor := ColorLogPlayerAttack
		if e == 0 {
			msgcolor = ColorLogMonsterAttack
		}
		if !s.ecs.BloodAt(s.ecs.GetComponentUnchecked(e, Position{}).(Position).Point) { // Add blood
			s.ecs.Create(
				Name{"blood"},
				Position{s.ecs.GetComponentUnchecked(e, Position{}).(Position).Point},
				NewRenderable('.', ColorCorpse, ColorBlood, ROFloor),
			)
		}
		s.ecs.Create(LogEntry{Text: msg, Color: msgcolor})
		if health.hp <= 0 {
			health.hp = 0
		}
	}
	s.ecs.RemoveComponent(e, DamageEffects{}) // Consume the damage effects.
	s.ecs.AddComponent(e, health)             // Update health.
	// s.printDebug(e) // Debugging output.
	// Uncomment the following lines to print debug information.
	// fmt.Printf("Entity: %d\n", e)
	// fmt.Printf("Health: %d\n", health.hp)
	// fmt.Printf("Damage Effects: %+v\n", dmgfx.effects)
	// fmt.Println("Components:", s.ecs.GetComponentsFor(e))
	// fmt.Println("====================")
	// Uncomment the following lines to print debug information.
	// s.printDebug(e) // Debugging output.
}

type FOVSystem struct {
	ecs *ECS
}

// Allows entities with FOV{} and Position{} to compute their field of view
// and mark cells within that FOV as explored. Typically only the player has
// this component, but other entities such as mobs can have them such as well,
// such as when the player drinks a potion of telepathy.
func (s *FOVSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Position{}, FOV{}) {
		return
	}
	p := s.ecs.GetComponentUnchecked(e, Position{}).(Position)
	f := s.ecs.GetComponentUnchecked(e, FOV{}).(FOV)
	if f.FOV == nil {
		f.FOV = rl.NewFOV(gruid.NewRange(-f.LOS, -f.LOS, f.LOS+1, f.LOS+1))
	}
	// We shift the FOV's range so that it will be centered on the position
	// of the entity.
	rg := gruid.NewRange(-f.LOS, -f.LOS, f.LOS+1, f.LOS+1)
	f.FOV.SetRange(rg.Add(p.Point).Intersect(s.ecs.Map.Grid.Range()))
	s.ecs.AddComponent(e, f)
	// We mark cells in field of view as explored. We use the symmetric shadow
	// casting algorithm provided by the rl package.
	isnotwall := func(p gruid.Point) bool {
		return s.ecs.Map.Grid.At(p) != Wall
	}
	for _, point := range f.FOV.SSCVisionMap(p.Point, f.LOS, isnotwall, true) {
		if paths.DistanceChebyshev(point, p.Point) > f.LOS {
			continue
		}
		if !s.ecs.Map.Explored[point] {
			s.ecs.Map.Explored[point] = true
		}
	}
}

// InFOV returns true if p is in the field of view of an entity with FOV. We only
// keep cells within maxLOS chebyshev distance from the source entity.
func (g *game) InFOV(p gruid.Point) bool {
	for _, e := range g.ECS.EntitiesWith(Position{}, FOV{}) {
		pos := g.ECS.GetComponentUnchecked(e, Position{}).(Position)
		fov := g.ECS.GetComponentUnchecked(e, FOV{}).(FOV)
		if fov.FOV.Visible(p) && paths.DistanceChebyshev(pos.Point, p) <= fov.LOS {
			return true
		}
	}
	return false
}

type DeathSystem struct {
	ecs *ECS
}

func (s *DeathSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Health{}) {
		return
	}
	health := s.ecs.GetComponentUnchecked(e, Health{}).(Health)
	if health.hp > 0 {
		return
	}
	name := s.ecs.GetComponentUnchecked(e, Name{}).(Name).string
	pos := s.ecs.GetComponentUnchecked(e, Position{}).(Position)
	var fov FOV
	if e == 0 {
		fov = s.ecs.GetComponentUnchecked(e, FOV{}).(FOV)
	}
	s.ecs.ClearAllComponents(e) // Clear all components of the entity.
	s.ecs.AddComponents(e,
		Name{name + " corpse"},
		NewRenderableNoBg('%', ColorCorpse, ROCorpse),
		Position{pos.Point},
		Collectible{},
		Consumable{},
		Healing{amount: 2},
		Dead{},
		// TODO drop inventory items
	)
	msg := fmt.Sprintf("%s has died!", name)
	if e == 0 {
		msg = "You have died!"
		s.ecs.AddComponent(e, fov)
	}
	s.ecs.Create(LogEntry{
		Text:  msg,
		Color: ColorLogMonsterAttack,
	})
}

type AnimationSystem struct {
	ecs *ECS
}

// Updates all Animation objects in the ECS forward a tick.
func (s *AnimationSystem) Update(e int) {

	if !s.ecs.HasComponent(e, Animation{}) {
		return
	}

	// Advance animation by a single tick.
	anim := s.ecs.GetComponentUnchecked(e, Animation{}).(Animation)
	anim.frames[anim.index].itick++

	// If the current frame has expired, move to the next frame.
	if anim.frames[anim.index].itick >= anim.frames[anim.index].nticks {
		anim.frames[anim.index].itick = 0
		anim.index++
	}

	// If the current animation has expired, remove it from the ECS or restart.
	if anim.index >= len(anim.frames) {
		anim.index = 0
		if anim.repeat == 0 {
			s.ecs.Delete(e)
		} else if anim.repeat > 0 {
			anim.repeat--
		}
	}
}

type DebugSystem struct {
	ecs *ECS
}

// Prints out component information for every entity.
func (s *DebugSystem) Update() {
	fmt.Println("+++++++++++++++++++++++++++++++ DEBUG ++++++++++++++++++++++++++++++++++++++++++")
	for _, e := range s.ecs.entities {
		s.ecs.printDebug(e)
	}
}
