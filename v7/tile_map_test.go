package sameriver

import (
	"os"
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

	tm := NewTextureManager()

	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm := NewTileManager(renderer, tm).SetDimension(32)
		tm.LoadTile("grass", "assets/tile_grass.bmp")
		tm.LoadTile("water", "assets/tile_water.bmp")

		tmap := NewTileMap(tm, 100, 100)
		tmap.SetTile(3, 3, "grass")
		tmap.SetTile(3, 4, "grass")
		tmap.SetTile(3, 5, "grass")
		tmap.SetTile(4, 3, "grass")
		tmap.SetTile(4, 4, "water")
		tmap.SetTile(4, 5, "grass")
		tmap.SetTile(5, 3, "grass")
		tmap.SetTile(5, 4, "grass")
		tmap.SetTile(5, 5, "grass")
		tmap.DrawTiles(&Viewport{100, 100, 200, 200}, window)

		renderer.Present()
		time.Sleep(1000 * time.Millisecond)
		window.Destroy()
	})
}

func TestTileMapSaveLoad(t *testing.T) {
	// create test window
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      300,
		Height:     300,
		Fullscreen: false}

	tm := NewTextureManager()

	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		defer window.Destroy()
		defer renderer.Destroy()

		tm := NewTileManager(renderer, tm).SetDimension(32)
		tm.LoadTile("grass", "assets/tile_grass.bmp")
		tm.LoadTile("water", "assets/tile_water.bmp")

		tmap := NewTileMap(tm, 10, 10)
		tmap.SetTile(3, 3, "grass")
		tmap.SetTile(4, 4, "water")

		tmap.Save("test.json")

		tmap2 := LoadTileMap(renderer, "test.json")
		if tmap.tm.Files["grass"] != tmap2.tm.Files["grass"] {
			t.Errorf("TileManager save/load failed")
		}
		if tmap.tm.Dimension != tmap2.tm.Dimension {
			t.Errorf("TileManager save/load failed")
		}

		if tmap.Tiles[3][3] != tmap2.Tiles[3][3] {
			t.Errorf("TileMap save/load failed")
		}
		if tmap.Tiles[4][4] != tmap2.Tiles[4][4] {
			t.Errorf("TileMap save/load failed")
		}

		os.Remove("test.json")
	})
}

func TestTileMapPerlinTerrain(t *testing.T) {
	// create test window
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      300,
		Height:     300,
		Fullscreen: false}

	tm := NewTextureManager()

	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		defer window.Destroy()
		defer renderer.Destroy()

		tm := NewTileManager(renderer, tm).SetDimension(32)
		tm.LoadTile("grass", "assets/tile_grass.bmp")
		tm.LoadTile("water", "assets/tile_water.bmp")

		tmap := NewTileMap(tm, 100, 100)

		seed := 108
		terrain := GeneratePerlinTerrain(seed, 100, 100)
		for x := 0; x < 100; x++ {
			for y := 0; y < 100; y++ {
				if terrain[x][y] > 0 {
					tmap.SetTile(int32(x), int32(y), "grass")
				} else {
					tmap.SetTile(int32(x), int32(y), "water")
				}
			}
		}

		tmap.DrawTiles(&Viewport{0, 0, 800, 800}, window)

		renderer.Present()
		time.Sleep(5000 * time.Millisecond)
	})
}
