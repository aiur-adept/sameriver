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
	*e.GetVec2D(VELOCITY) = Vec2D{1, 1}

	w.em.components.Save("test.json")

	ct := ComponentTableFromJSON("test.json")
	Logger.Println(ct)

	// check if the component table is the same as the original
	if ct.Vec2DMap[VELOCITY][e.ID] != *e.GetVec2D(VELOCITY) {
		t.Errorf("Vec2DMap[%v][0] = %v, want %v", VELOCITY, ct.Vec2DMap[VELOCITY][0], *e.GetVec2D(VELOCITY))
	}

	os.Remove("test.json")
}
