package sameriver

import (
	"testing"
	"time"
)

// test the tile manager
func TestTileManager(t *testing.T) {
	// create test window
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      300,
		Height:     300,
		Fullscreen: false}
	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm := NewTileManager(renderer, 32, 32)
		tm.LoadTile("grass", "assets/tile_grass.bmp")

		vp := &Viewport{20, 20, 50, 50}
		for i := 0; i < 10; i++ {
			vp.Width += 10
			vp.Height += 10
			time.Sleep(100 * time.Millisecond)
			renderer.Clear()
			tm.DrawTile("grass", 0, 0, vp, window)
			tm.DrawTile("grass", 32, 32, vp, window)
			tm.DrawTile("grass", 0, 32, vp, window)
			tm.DrawTile("grass", 32, 0, vp, window)
			renderer.Present()
		}
		time.Sleep(2000 * time.Millisecond)
		window.Destroy()
	})
}
