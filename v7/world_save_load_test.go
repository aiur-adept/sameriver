package sameriver

import (
	"testing"
)

func TestWorldSaveLoad(t *testing.T) {
	w := testingWorld()
	p := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(p, cs)

	e := testingSpawnPhysics(w)

	bb := w.CreateBlackboard("testbb")
	bb.Set("test", e.ID)

	w.Save("test.json")

	w2 := LoadWorld("test.json")
	p2 := NewPhysicsSystem()
	cs2 := NewCollisionSystem(FRAME_DURATION / 2)
	w2.RegisterSystems(p2, cs2)

	// check if e is in w2
	e2 := w2.GetEntity(e.ID)
	if e2 == nil {
		t.Fatalf("entity %d not found in world", e.ID)
	}

	if e2.ID != e.ID {
		t.Fatalf("entity %d not found in world", e.ID)
	}

	// check if e2 has the same components as e
	for _, c := range e.Components {
		if !w2.EntityHasComponentString(e2, c) {
			t.Fatalf("entity %d does not have component %s", e.ID, c)
		}
	}

	// try to get "test" key from blackboard "testbb"
	test := w2.Blackboards["testbb"].GetInt("test")
	if test != e.ID {
		t.Fatalf("test key not found in blackboard")
	}

	w2.GetVec2D(e2, VELOCITY_).X = 12.0
	w2.Update(FRAME_MS)
	if w.GetVec2D(e, POSITION_).Equals(*w2.GetVec2D(e2, POSITION_)) {
		t.Fatalf("entity %d did not move", e.ID)
	}
}
