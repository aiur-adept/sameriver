package sameriver

import (
	"testing"
)

func TestSteeringSystem(t *testing.T) {
	w := testingWorld()
	ss := NewSteeringSystem()
	w.RegisterSystems(ss)
	e := testingSpawnSteering(w)
	vel := *e.GetVec2D(VELOCITY_)

	w.Update(1)
	w.Update(FRAME_MS / 2)
	if *e.GetVec2D(VELOCITY_) == vel {
		t.Fatal("failed to steer velocity")
	}
}
