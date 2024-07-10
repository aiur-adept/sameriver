package sameriver

import (
	"testing"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func TestSpriteSystemBasic(t *testing.T) {
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      100,
		Height:     100,
		Fullscreen: false}

	sdl.Init(sdl.INIT_VIDEO)

	SDLMainMediaThread(func() {

		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		defer func() {
			window.Destroy()
			renderer.Destroy()
		}()

		img.Init(img.INIT_PNG)

		w := NewWorld(nil)
		tm := NewTextureManager()
		tm.Init(renderer)

		ss := NewSpriteSystem(renderer, tm)

		w.RegisterSystems(ss)

		e := w.Spawn(map[string]any{
			"components": map[ComponentID]any{
				POSITION_:   Vec2D{X: 0, Y: 0},
				BASESPRITE_: ss.GetSprite("test", 614, 800, 4, 4),
			},
		})

		ss.Render(renderer, e, w.GetSprite(e, BASESPRITE_))

		sdl.Delay(1000)

		// fail the test - this is a TODO
		t.Fatal("test failed")

	})
}
