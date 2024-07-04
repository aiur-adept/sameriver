package sameriver

import "github.com/veandco/go-sdl2/sdl"

type Viewport struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

func (vp *Viewport) DestRect(window *sdl.Window, x, y, w, h int32) sdl.Rect {
	x -= int32(vp.X)
	y -= int32(vp.Y)

	// get the scale of the viewport relative to window size
	ww, wh := window.GetSize()
	scaleX := 1.0 / (float32(vp.Width) / float32(ww))
	scaleY := 1.0 / (float32(vp.Height) / float32(wh))

	return sdl.Rect{
		int32(float32(x) * scaleX),
		int32(float32(y) * scaleY),
		int32(float32(w) * scaleX),
		int32(float32(h) * scaleY),
	}
}
