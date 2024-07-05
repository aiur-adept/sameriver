package sameriver

import (
	"os"
	"testing"
)

func TestComponentTableSave(t *testing.T) {

	// normal setup
	w := testingWorld()
	p := NewPhysicsSystem()
	w.RegisterSystems(p)
	e := testingSpawnPhysics(w)
	*e.GetVec2D(VELOCITY_) = Vec2D{1, 1}

	w.em.components.Save("test.json")

	ct := ComponentTableFromJSON("test.json")
	Logger.Println(ct)

	// check if the component table is the same as the original
	if ct.Vec2DMap[VELOCITY_][e.ID] != *e.GetVec2D(VELOCITY_) {
		t.Errorf("Vec2DMap[%v][0] = %v, want %v", VELOCITY_, ct.Vec2DMap[VELOCITY_][0], *e.GetVec2D(VELOCITY_))
	}

	os.Remove("test.json")
}
