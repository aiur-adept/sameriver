package sameriver

import (
	"testing"
	"time"
)

func TestPhysicsSystemWithGranularity(t *testing.T) {
	// normal setup
	w := testingWorld()
	p := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(p, cs)
	e := testingSpawnPhysics(w)
	*w.GetVec2D(e, VELOCITY_) = Vec2D{1, 1}
	pos := w.GetVec2D(e, POSITION_)
	pos0 := *pos
	// granular setup
	wg := testingWorld()
	pg := NewPhysicsSystemWithGranularity(4)
	csg := NewCollisionSystem(FRAME_DURATION / 2)
	wg.RegisterSystems(pg, csg)
	eg := testingSpawnPhysics(wg)
	*wg.GetVec2D(eg, VELOCITY_) = Vec2D{1, 1}
	posg := wg.GetVec2D(eg, POSITION_)
	posg0 := *posg

	// simulate constant load of other logics with a ratio
	// value of 8 means 1/8th to physics
	physicsTimeShareReciprocal := 8.0

	// run a frame of same allowance ms for both normal and granular
	runFrame := func() {
		w.Update(FRAME_MS / physicsTimeShareReciprocal)
		wg.Update(FRAME_MS / physicsTimeShareReciprocal)
		time.Sleep(FRAME_DURATION)
	}

	observePos := func() {
		Logger.Printf("normal pos: %v", *pos)
		Logger.Printf("granular pos: %v", *posg)
		Logger.Printf("position.x ratio: %f", pos.X/posg.X)
	}

	// Frame 0
	runFrame()
	Logger.Println("after Update at t=0, hotness of physics update:")
	// observe, in the below, the hotness are basically the same for
	// granularity 1 as granularity 4, insane
	for _, l := range w.RuntimeSharer.RunnerMap["systems"].logicUnits {
		Logger.Printf("normal %s: h%d", l.name, l.hotness)
	}
	for _, l := range wg.RuntimeSharer.RunnerMap["systems"].logicUnits {
		Logger.Printf("granular %s: h%d", l.name, l.hotness)
	}
	observePos()

	// Frame 1
	Logger.Println("TEST FRAME 2")
	runFrame()
	if *pos == pos0 {
		t.Fatal("failed to update position")
	}
	if *posg == posg0 {
		t.Fatal("failed to update position in granular")
	}
	// as of this comment, 2023-03-18, observe that the numeric result is different;
	// this is *at least* because
	// the physics update is getting slightly different unstable dts, due ultimately to
	// the runtimelimiter passing different dt_ms to the logicunit each time it runs
	// based on wall time since it last scheduled it, which can vary over a single
	// frame as it tries to pack in the time and repeatedly polls time since last
	// run.
	observePos()

	// let's observe the behaviour over a longer term
	for i := 0; i < 30; i++ {
		runFrame()
	}
	Logger.Printf("+ 30 frames:")
	// observe that they're actually closer than after the 1st or 2nd frame
	observePos()
	for i := 0; i < 60; i++ {
		runFrame()
	}
	Logger.Printf("+ 60 frames:")
	// observe that they're actually closer than after the 1st or 2nd frame
	observePos()
}

func TestPhysicsSystemMotion(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(ps, cs)
	e := testingSpawnPhysics(w)
	*w.GetVec2D(e, VELOCITY_) = Vec2D{1, 1}
	pos := *w.GetVec2D(e, POSITION_)
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_MS / 2)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_MS / 2)
	if *w.GetVec2D(e, POSITION_) == pos {
		t.Fatal("failed to update position")
	}
}

func TestPhysicsSystemMany(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(ps, cs)
	for i := 0; i < 500; i++ {
		testingSpawnPhysics(w)
	}
	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_MS / 2)
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_MS / 2)
}

func TestPhysicsSystemBounds(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(ps, cs)
	e := testingSpawnPhysics(w)
	directions := []Vec2D{
		Vec2D{100, 0},
		Vec2D{-100, 0},
		Vec2D{0, 100},
		Vec2D{0, -100},
	}
	worldCenter := Vec2D{w.Width / 2, w.Height / 2}
	worldTopRight := Vec2D{w.Width, w.Height}
	pos := w.GetVec2D(e, POSITION_)
	box := w.GetVec2D(e, BOX_)
	vel := w.GetVec2D(e, VELOCITY_)
	for _, d := range directions {
		*pos = Vec2D{512, 512}
		*vel = d
		for i := 0; i < 64; i++ {
			w.Update(FRAME_MS / 2)
			time.Sleep(1 * time.Millisecond)
		}
		if !RectWithinRect(*pos, *box, worldCenter, worldTopRight) {
			t.Fatalf("traveling with velocity %v placed entity "+
				"outside world (at position %v, box %v)", *vel, *pos, *box)
		}
	}
}

func TestPhysicsSystemRigidBody(t *testing.T) {
	w := testingWorld()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	ps := NewPhysicsSystem()
	w.RegisterSystems(ps, cs)

	ec := cs.w.Events.Subscribe(SimpleEventFilter("collision"))

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{0, 0},
			// velocity slow enough to avoid tunneling
			VELOCITY_:     Vec2D{0.1, 0.1},
			ACCELERATION_: Vec2D{0, 0},
			BOX_:          Vec2D{1, 1},
			MASS_:         3.0,
			RIGIDBODY_:    true,
		}})

	e2 := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_:     Vec2D{10, 10},
			VELOCITY_:     Vec2D{0, 0},
			ACCELERATION_: Vec2D{0, 0},
			BOX_:          Vec2D{1, 1},
			MASS_:         3.0,
			RIGIDBODY_:    true,
		}})

	// Update twice since physics system won't run the first time(needs a dt)
	w.Update(FRAME_MS / 2)
	for i := 0; i < 100; i++ {
		time.Sleep(FRAME_DURATION)
		w.Update(FRAME_MS / 2)
	}
	Logger.Printf("e: %v", *w.GetVec2D(e, POSITION_))
	Logger.Printf("e2: %v", *w.GetVec2D(e2, POSITION_))

	// should have collision events
	select {
	case <-ec.C:
		break
	default:
		t.Fatal("collision event wasn't received within 1 frame")
	}
}
