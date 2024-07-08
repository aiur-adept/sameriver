package sameriver

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEntityInvalidComponentAccess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Should have paniced")
		}
	}()
	w := testingWorld()
	e := w.Spawn(nil)
	w.GetVec2D(e, 1337)
}

func TestEntitySaveLoad(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(ps, cs)
	e := testingSpawnPhysics(w)

	jsonStr, err := json.Marshal(e)
	if err != nil {
		t.Fatal("error marshalling entity")
	}
	fmt.Println(string(jsonStr))
	e2 := Entity{}
	json.Unmarshal([]byte(jsonStr), &e2)
	if e.ID != e2.ID {
		t.Fatal("did not save and load entity correctly")
	}
}

func TestEntitySaveLoadSlice(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(ps, cs)
	e := testingSpawnPhysics(w)

	entities := []Entity{*e}

	jsonStr, err := json.Marshal(entities)
	if err != nil {
		t.Fatal("error marshalling entity")
	}
	fmt.Println(string(jsonStr))
	e2 := []Entity{}
	json.Unmarshal([]byte(jsonStr), &e2)
	if e.ID != e2[0].ID {
		t.Fatal("did not save and load entity correctly")
	}
}
