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
				BOX_:        Vec2D{X: 32, Y: 48},
				BASESPRITE_: ss.GetSprite("test", 4, 4),
			},
		})

		// for loop 100 times
		animationFPS := 0.2
		animation_accum := NewTimeAccumulator(animationFPS * 1000)
		animation_logic := func(dt_ms float64) {
			if animation_accum.Tick(dt_ms) {
				sprite := w.GetSprite(e, BASESPRITE_)
				sprite.FrameX += 1
				sprite.FrameX %= sprite.DimX
			}
		}
		w.AddEntityLogic(e, "animation", animation_logic)
		for i := 0; i < 100; i++ {
			w.Update(FRAME_MS)
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(nil)
			ss.Render(renderer, e, w.GetSprite(e, BASESPRITE_))
			renderer.Present()

			sdl.Delay(FRAME_MS)

		}

		// fail the test - this is a TODO
		t.Fatal("test failed")

	})
}
