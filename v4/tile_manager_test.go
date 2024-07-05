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

	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm := NewTileManager(renderer).SetDimension(32)
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

func TestTileManagerSaveLoad(t *testing.T) {
	// create test window
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      300,
		Height:     300,
		Fullscreen: false}

	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		defer window.Destroy()
		defer renderer.Destroy()

		tm := NewTileManager(renderer).SetDimension(32)
		tm.LoadTile("grass", "assets/tile_grass.bmp")
		tm.Save("test.json")

		tm2 := LoadTileManager(renderer, "test.json")
		if tm.Files["grass"] != tm2.Files["grass"] {
			t.Errorf("TileManager save/load failed")
		}
		if tm.TileDimension != tm2.TileDimension {
			t.Errorf("TileManager save/load failed")
		}
	})
}
