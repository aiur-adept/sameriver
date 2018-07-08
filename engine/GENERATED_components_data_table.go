
//
//
// THIS FILE HAS BEEN GENERATED BY sameriver-generate
//
//
// DO NOT MODIFY BY HAND UNLESS YOU WANNA HAVE A GOOD TIME WHEN THE NEXT
// GENERATION DESTROYS WHAT YOU WROTE. UNLESS YOU KNOW HOW TO HAVE A GOOD TIME
//
//

package engine

type ComponentsTable struct {
	em             *EntityManager
	Box            [MAX_ENTITIES]Vec2D
	Logic          [MAX_ENTITIES]LogicUnit
	MovementTarget [MAX_ENTITIES]Vec2D
	Position       [MAX_ENTITIES]Vec2D
	Sprite         [MAX_ENTITIES]Sprite
	Steer          [MAX_ENTITIES]float64
	TagList        [MAX_ENTITIES]TagList
	Velocity       [MAX_ENTITIES]Vec2D
}

func NewComponentsTable(em *EntityManager) *ComponentsTable {
	return &ComponentsTable{em: em}
}
