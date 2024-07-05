package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

// TODO: store only the string name of the texture, for lookup in a SpriteManager
// which will load the texture on demand
type Sprite struct {
	Texture string
	Frame   uint8
	Visible bool
	Flip    sdl.RendererFlip
}
