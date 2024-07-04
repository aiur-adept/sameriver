package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	Texture *sdl.Texture
	Frame   uint8
	Visible bool
	Flip    sdl.RendererFlip
}

func NewSprite(renderer *sdl.Renderer, filename string) *Sprite {
	surface, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(err)
	}
	// create texture from surface
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	return &Sprite{Texture: texture}
}
