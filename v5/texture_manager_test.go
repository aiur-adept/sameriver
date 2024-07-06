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

	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm.Init(renderer)

		for i := 0; i < 3; i++ {
			tm.Render(renderer, "tile_grass", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
			tm.Render(renderer, "tile_water", 0, 0, 100, 100)
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

	SDLMainMediaThread(func() {
		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		tm.Init(renderer)

		jsonStr := tm.String()

		fmt.Println(jsonStr)

		tm2 := NewTextureManager()
		tm2.UnmarshalJSON([]byte(jsonStr))
		tm2.LoadFiles(renderer)

		for i := 0; i < 3; i++ {
			tm2.Render(renderer, "tile_grass", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
			tm2.Render(renderer, "tile_water", 0, 0, 100, 100)
			renderer.Present()
			time.Sleep(500 * time.Millisecond)
		}
		window.Destroy()
	})
}
