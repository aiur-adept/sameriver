package sameriver

import (
	"math/rand"
	"testing"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func TestSpriteSystemBasic(t *testing.T) {
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      500,
		Height:     500,
		Fullscreen: false}

	sdl.Init(sdl.INIT_VIDEO)

	SDLMainMediaThread(func() {

		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		defer func() {
			window.Destroy()
			renderer.Destroy()
		}()

		img.Init(img.INIT_PNG)

		w := NewWorld(map[string]any{
			"height": 500,
			"width":  500,
		})
		tm := NewTextureManager()
		tm.Init(renderer)

		p := NewPhysicsSystem()
		c := NewCollisionSystem(FRAME_DURATION)
		ss := NewSpriteSystem(renderer, tm)

		w.RegisterSystems(p, c, ss)

		e := w.Spawn(map[string]any{
			"components": map[ComponentID]any{
				POSITION_:     Vec2D{X: 16, Y: 24},
				VELOCITY_:     Vec2D{X: 0, Y: 0.1},
				ACCELERATION_: Vec2D{X: 0, Y: 0},
				RIGIDBODY_:    false,
				MASS_:         3.0,
				BOX_:          Vec2D{X: 32, Y: 48},
				BASESPRITE_:   ss.GetSprite("test", 4, 4),
			},
		})

		other := testingSpawnSimple(w)

		animationFPS := 0.2
		animation_accum := NewTimeAccumulator(animationFPS * 1000)
		animation_controller := func(e *Entity, dt_ms float64) {
			if animation_accum.Tick(dt_ms) {
				sprite := w.GetSprite(e, BASESPRITE_)
				sprite.FrameX += 1
				sprite.FrameX %= sprite.DimX
			}
		}
		ss.AddSpriteController(e, NewSpriteController(animation_controller))
		for i := 0; i < 100; i++ {
			w.Update(FRAME_MS)
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(nil)
			ss.Render(renderer, e, w.GetSprite(e, BASESPRITE_))
			renderer.Present()

			time.Sleep(time.Millisecond * FRAME_MS)
		}

		// make sure the despawn callback doesn't error for an entity without a spritecontroller
		w.Despawn(other)

	})

}

func TestSpriteSystemWander(t *testing.T) {
	skipCI(t)

	windowSpec := WindowSpec{
		Title:      "testing game",
		Width:      500,
		Height:     500,
		Fullscreen: false}

	sdl.Init(sdl.INIT_VIDEO)

	SDLMainMediaThread(func() {

		window, renderer := SDLCreateWindowAndRenderer(windowSpec)
		defer func() {
			window.Destroy()
			renderer.Destroy()
		}()

		img.Init(img.INIT_PNG)

		w := NewWorld(map[string]any{
			"height": 500,
			"width":  500,
		})
		tm := NewTextureManager()
		tm.Init(renderer)

		p := NewPhysicsSystem()
		c := NewCollisionSystem(FRAME_DURATION)
		ss := NewSpriteSystem(renderer, tm)

		w.RegisterSystems(p, c, ss)

		e := w.Spawn(map[string]any{
			"components": map[ComponentID]any{
				POSITION_:     Vec2D{X: 100, Y: 100},
				VELOCITY_:     Vec2D{X: 0, Y: 0},
				ACCELERATION_: Vec2D{X: 0, Y: 0},
				RIGIDBODY_:    false,
				MASS_:         3.0,
				BOX_:          Vec2D{X: 32, Y: 48},
				BASESPRITE_:   ss.GetSprite("test", 4, 4),
			},
		})

		wanderFPS := 0.7
		wander_accum := NewTimeAccumulator(wanderFPS * 1000)
		w.AddEntityLogic(e, "wander", func(dt_ms float64) {
			if wander_accum.Tick(dt_ms) {
				vel := w.GetVec2D(e, VELOCITY_)
				sprite := w.GetSprite(e, BASESPRITE_)
				if rand.Float64() < 0.5 {
					*vel = Vec2D{X: 0, Y: 0}
					return
				}
				// choose to either wander left or right
				if rand.Float64() < 0.5 {
					if rand.Float64() < 0.5 {
						*vel = Vec2D{X: -0.1, Y: 0}
						sprite.FrameY = 2
					} else {
						*vel = Vec2D{X: 0.1, Y: 0}
						sprite.FrameY = 1
					}
				} else {
					// choose to either wander up or down
					if rand.Float64() < 0.5 {
						*vel = Vec2D{X: 0, Y: -0.1}
						sprite.FrameY = 3
					} else {
						*vel = Vec2D{X: 0, Y: 0.1}
						sprite.FrameY = 0
					}
				}
			}
		})

		animationFPS := 0.2
		animation_accum := NewTimeAccumulator(animationFPS * 1000)
		animation_controller := func(e *Entity, dt_ms float64) {
			sprite := w.GetSprite(e, BASESPRITE_)

			if animation_accum.Tick(dt_ms) {
				sprite.FrameX += 1
				sprite.FrameX %= sprite.DimX
			}

			vel := w.GetVec2D(e, VELOCITY_)
			if vel.X == 0 && vel.Y == 0 {
				if sprite.FrameY == 2 {
					sprite.FrameX = 0
				} else {
					sprite.FrameX = 1
				}
			}
		}
		ss.AddSpriteController(e, NewSpriteController(animation_controller))

		for i := 0; i < 200; i++ {
			w.Update(FRAME_MS)
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(nil)
			ss.Render(renderer, e, w.GetSprite(e, BASESPRITE_))
			renderer.Present()

			time.Sleep(time.Millisecond * FRAME_MS)
		}

	})

}
