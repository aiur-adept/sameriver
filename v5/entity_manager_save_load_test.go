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
}
