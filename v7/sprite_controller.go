package sameriver

type SpriteController struct {
	Update func(e *Entity, dt_ms float64)
}

func NewSpriteController(update func(e *Entity, dt_ms float64)) *SpriteController {
	return &SpriteController{Update: update}
}
