
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

func (e *Entity) GetBox() *Vec2D {
	return &e.World.Components.Box[e.ID]
}
func (e *Entity) GetLogic() *LogicUnit {
	return &e.World.Components.Logic[e.ID]
}
func (e *Entity) GetMass() *float64 {
	return &e.World.Components.Mass[e.ID]
}
func (e *Entity) GetMaxVelocity() *float64 {
	return &e.World.Components.MaxVelocity[e.ID]
}
func (e *Entity) GetMovementTarget() *Vec2D {
	return &e.World.Components.MovementTarget[e.ID]
}
func (e *Entity) GetPosition() *Vec2D {
	return &e.World.Components.Position[e.ID]
}
func (e *Entity) GetSprite() *Sprite {
	return &e.World.Components.Sprite[e.ID]
}
func (e *Entity) GetSteer() *Vec2D {
	return &e.World.Components.Steer[e.ID]
}
func (e *Entity) GetTagList() *TagList {
	return &e.World.Components.TagList[e.ID]
}
func (e *Entity) GetVelocity() *Vec2D {
	return &e.World.Components.Velocity[e.ID]
}
