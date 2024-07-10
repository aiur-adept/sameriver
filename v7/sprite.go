package sameriver

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	Texture string

	// what frame we are presenting
	FrameX uint8
	FrameY uint8

	// the dimensions in the source image of a frame
	FrameW uint8
	FrameH uint8

	// how many frames are in the sheet in X and Y
	DimX uint8
	DimY uint8

	Visible bool
	Flip    sdl.RendererFlip
}
