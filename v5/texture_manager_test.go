package sameriver

import (
	"fmt"
	"testing"
	"time"
)

func TestTextureManagerBasic(t *testing.T) {
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      100,
		Height:     100,
		Fullscreen: false}
	tm := NewTextureManager()

	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm.LoadTexture(renderer, "assets/tile_grass.bmp", "grass")
		tm.LoadTexture(renderer, "assets/tile_water.bmp", "water")

		for i := 0; i < 3; i++ {
			tm.Render(renderer, "grass", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
			tm.Render(renderer, "water", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
		}
		window.Destroy()
	})
}

func TestTextureManagerSaveLoad(t *testing.T) {
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      100,
		Height:     100,
		Fullscreen: false}
	tm := NewTextureManager()

	// in a real game, the scene Init() gets a Game object and creates a new
	// sprite system by passing game.Renderer
	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm.LoadTexture(renderer, "assets/tile_grass.bmp", "grass")
		tm.LoadTexture(renderer, "assets/tile_water.bmp", "water")

		jsonStr := tm.String()

		fmt.Println(jsonStr)

		tm2 := NewTextureManager()
		tm2.UnmarshalJSON([]byte(jsonStr))
		tm2.LoadFiles(renderer)

		for i := 0; i < 3; i++ {
			tm2.Render(renderer, "grass", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
			tm2.Render(renderer, "water", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
		}
		window.Destroy()
	})
}
