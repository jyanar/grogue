// Implementing the ECS!

package main

// import "github.com/anaseto/gruid"

// // type Entity int
// type Component struct{}
// type System interface {
// 	Update()
// }

// type MovementSystem struct {
// 	ecs *ECS
// }

// type ECS struct {
// 	Entities  []Entity            // List of entities in the world.
// 	Positions map[int]gruid.Point // Entity index: map position.
// 	PlayerID  int                 // Index of player entity (for convenience).
// }

// func NewECS() *ECS {
// 	return &ECS{
// 		Positions: map[int]gruid.Point{},
// 	}
// }

// func (es *ECS) AddEntity(e Entity, p gruid.Point) int {
// 	i := len(es.Entities)
// 	es.Entities = append(es.Entities, e)
// 	es.Positions[i] = p
// 	return i
// }

// func (es *ECS) MoveEntity(i int, p gruid.Point) {
// 	es.Positions[i] = p
// }

// func (es *ECS) MovePlayer(p gruid.Point) {
// 	es.MoveEntity(es.PlayerID, p)
// }

// func (es *ECS) Player() *Player {
// 	return es.Entities[es.PlayerID].(*Player) // index 0 for player entity (convention)
// }

// // // Represents an object or creature on the map.
// // type Entity interface {
// // 	Rune() rune         // The glyph representing this entity.
// // 	Color() gruid.Color // The color of this entity.
// // }

// type Player struct{}

// func (p *Player) Rune() rune {
// 	return '@'
// }

// func (p *Player) Color() gruid.Color {
// 	return gruid.ColorDefault
// }
