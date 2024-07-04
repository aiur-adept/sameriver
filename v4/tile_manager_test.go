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
		tm := NewTileManager(renderer)
		tm.LoadTile("grass", "assets/tile_grass.bmp")

		tm.DrawTile("grass", 0, 0, 32, 32)
		tm.DrawTile("grass", 32, 32, 32, 32)
		tm.DrawTile("grass", 0, 32, 32, 32)
		tm.DrawTile("grass", 32, 0, 32, 32)
		renderer.Present()
		time.Sleep(200 * time.Millisecond)
		window.Destroy()
	})
}
