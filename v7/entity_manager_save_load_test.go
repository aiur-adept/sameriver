package sameriver

import (
	"encoding/json"
	"os"
	"testing"
)

func TestEntityManagerSaveLoad(t *testing.T) {
	w := testingWorld()
	p := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(p, cs)

	e := testingSpawnPhysics(w)
	_ = testingSpawnPhysics(w)

	w.Em.Save("test.json")

	// read jsonStr from test.json
	jsonStr, err := os.ReadFile("test.json")
	if err != nil {
		t.Fatal(err)
	}

	em2 := NewEntityManager(w)
	// unmarshal from jsonStr
	err = json.Unmarshal([]byte(jsonStr), &em2)
	if err != nil {
		t.Fatal(err)
	}

	if em2.GetEntityByID(e.ID) == nil {
		t.Fatal("entity not in em2")
	}

	// check if position is the same
	if w.GetVec2D(em2.GetEntityByID(e.ID), POSITION_) != w.GetVec2D(e, POSITION_) {
		t.Fatal("entity position not the same")
	}
	// check if mass is the same
	if w.GetFloat64(em2.GetEntityByID(e.ID), MASS_) != w.GetFloat64(e, MASS_) {
		t.Fatal("entity mass not the same")
	}
}
