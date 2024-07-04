package sameriver

import (
	"testing"
	"time"
)

// test the tile manager
func TestTileMap(t *testing.T) {
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
		tm.LoadTile("water", "assets/tile_water.bmp")

		tmap := NewTileMap(100, 100)
		tmap.SetTile(3, 3, "grass")
		tmap.SetTile(3, 4, "grass")
		tmap.SetTile(3, 5, "grass")
		tmap.SetTile(4, 3, "grass")
		tmap.SetTile(4, 4, "water")
		tmap.SetTile(4, 5, "grass")
		tmap.SetTile(5, 3, "grass")
		tmap.SetTile(5, 4, "grass")
		tmap.SetTile(5, 5, "grass")
		tmap.DrawTiles(tm, &Viewport{0, 0, 300, 300}, window)

		renderer.Present()
		time.Sleep(5000 * time.Millisecond)
		window.Destroy()
	})
}
