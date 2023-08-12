package main

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
