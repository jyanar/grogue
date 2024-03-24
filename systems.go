package main

import (
	"fmt"

	"github.com/anaseto/gruid"
	"github.com/anaseto/gruid/paths"
	"github.com/anaseto/gruid/rl"
)

type System interface {
	Update()
}

type PerceptionSystem struct {
	ecs *ECS
}

func (s *PerceptionSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Position{}, Perception{}) {
		return
	}
	pos := s.ecs.positions[e]
	per := s.ecs.perceptions[e]
	per.perceived = []int{}
	if per.FOV == nil {
		per.FOV = rl.NewFOV(gruid.NewRange(-per.LOS, -per.LOS, per.LOS+1, per.LOS+1))
	}
	rg := gruid.NewRange(-per.LOS, -per.LOS, per.LOS+1, per.LOS+1)
	per.FOV.SetRange(rg.Add(pos.Point).Intersect(s.ecs.Map.Grid.Range()))
	passable := func(p gruid.Point) bool {
		return s.ecs.Map.Grid.At(p) != Wall
	}
	for _, point := range per.FOV.SSCVisionMap(pos.Point, per.LOS, passable, false) {
		if paths.DistanceManhattan(point, pos.Point) > per.LOS {
			continue
		}
	}
	for _, other := range s.ecs.EntitiesWith(Position{}, Visible{}) {
		// Ignore self.
		if other == e {
			continue
		}
		// If other entity is within perceptive radius, add to perceived list.
		pos_other := s.ecs.positions[other]
		if per.FOV.Visible(pos_other.Point) {
			per.perceived = append(per.perceived, other)
		}
	}
	// Check if player is within perceived entities, and switch to appropriate state.
	for _, other := range per.perceived {
		name_other := s.ecs.names[other].string
		if name_other == "you" && s.ecs.HasComponent(e, AI{}) {
			s.ecs.ais[e].state = CSHunting
			break
		} else {
			s.ecs.ais[e].state = CSWandering
		}
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
	ai := s.ecs.ais[e]
	pos := s.ecs.positions[e]
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
					break
				}
			}
		}
	case CSHunting:
		// Set destination to be the player.
		ai.dest = &s.ecs.positions[0].Point
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
	// Get entity's bump and position data.
	b := s.ecs.bumps[e]
	p := s.ecs.positions[e]
	dest := p.Point.Add(b.Point)
	s.ecs.bumps[e] = nil // Consume bump component.
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
			return
		}
		if len(attackable_entities) > 1 {
			panic(fmt.Sprintf("More than one entity with obstruct at position: %v", dest))
		}
		// Attack entity at location.
		target_entity := attackable_entities[0]
		attack_power := s.ecs.damages[e].int
		s.ecs.AddComponent(target_entity, DamageEffect{source: e, amount: attack_power})
		s.ecs.DamageEffectSystem.Update(target_entity)
	} else {
		s.ecs.Create(LogEntry{Text: "The wall is firm and unyielding!", Color: ColorLogSpecial})
	}
}

// Processes damage effects on entities with health components.
type DamageEffectSystem struct {
	ecs *ECS
}

func (s *DamageEffectSystem) Update(e int) {
	if !s.ecs.HasComponents(e, DamageEffect{}, Health{}) {
		s.ecs.damageeffects[e] = []DamageEffect{}
		return
	}
	health := s.ecs.healths[e]
	for _, de := range s.ecs.damageeffects[e] {
		health.hp -= de.amount
		name_attacker := s.ecs.names[de.source].string
		name_receiver := s.ecs.names[e].string
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
		if !s.ecs.BloodAt(s.ecs.positions[e].Point) { // Add blood tile here.
			// s.ecs.Create(
			// 	Name{"blood"},
			// 	Position{s.ecs.positions[e].Point},
			// 	Renderable{cell: gruid.Cell{Rune: '.', Style: gruid.Style{Fg: ColorBlood}}, order: ROFloor},
			// )
			s.ecs.Create(
				Name{"blood"},
				Position{s.ecs.positions[e].Point},
				NewRenderable('.', ColorBlood, ColorBlood, ROFloor),
			)
		}
		s.ecs.Create(LogEntry{Text: msg, Color: msgcolor})
		if health.hp <= 0 {
			health.hp = 0
			// Process entity through DeathSystem
			s.ecs.AddComponent(e, Death{})
			s.ecs.DeathSystem.Update(e)
		}
	}
	s.ecs.damageeffects[e] = []DamageEffect{}
}

type FOVSystem struct {
	ecs *ECS
}

func (s *FOVSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Position{}, FOV{}) {
		return
	}
	p := s.ecs.positions[e]
	f := s.ecs.fovs[e]
	if f.FOV == nil {
		f.FOV = rl.NewFOV(gruid.NewRange(-f.LOS, -f.LOS, f.LOS+1, f.LOS+1))
	}
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

type DeathSystem struct {
	ecs *ECS
}

func (s *DeathSystem) Update(e int) {
	if !s.ecs.HasComponents(e, Death{}) {
		return
	}
	name := s.ecs.names[e].string
	fg := s.ecs.renderables[e].cell.Style.Fg
	s.ecs.obstructs[e] = nil   // No longer blocking.
	s.ecs.perceptions[e] = nil // No longer perceiving.
	s.ecs.ais[e] = nil         // No longer pathing.
	s.ecs.bumps[e] = nil
	s.ecs.inputs[e] = nil
	s.ecs.damages[e] = nil
	s.ecs.rangeds[e] = nil
	s.ecs.damageeffects[e] = nil
	// s.ecs.healths[e] = nil
	s.ecs.deaths[e] = nil // Consume the death component.
	if e == 0 {
		s.ecs.AddComponent(e, Name{"your corpse"})
	} else {
		s.ecs.AddComponent(e, Name{name + " corpse"})
	}
	s.ecs.AddComponent(e, NewRenderableNoBg('%', fg, ROCorpse))
	s.ecs.AddComponent(e, Collectible{})
	s.ecs.AddComponent(e, Consumable{})
	s.ecs.AddComponent(e, Healing{amount: 2})
	// Drop everything in inventory
	if s.ecs.HasComponent(e, Inventory{}) {
		for _, item := range s.ecs.inventories[e].items {
			s.ecs.AddComponent(item, Position{s.ecs.positions[e].Point})
		}
		s.ecs.inventories[e] = nil
	}
	msg := fmt.Sprintf("%s has died!", name)
	if e == 0 {
		msg = "You have died!"
	}
	s.ecs.Create(LogEntry{
		Text:  msg,
		Color: ColorLogMonsterAttack,
	})
}

type AnimationSystem struct {
	ecs *ECS
}

// Updates all CAnimation objects in the ECS forward a tick.
func (s *AnimationSystem) Update(e int) {

	if !s.ecs.HasComponent(e, CAnimation{}) {
		return
	}

	// Advance animation by a single tick.
	anim := s.ecs.animations[e]
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
