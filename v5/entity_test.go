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
	e.GetVec2D(1337)
}

func TestEntitySaveLoad(t *testing.T) {
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(ps, cs)
	e := testingSpawnPhysics(w)

	jsonStr := e.String()
	fmt.Println(jsonStr)
	e2 := Entity{World: w}
	json.Unmarshal([]byte(jsonStr), &e2)
	if e.ID != e2.ID {
		t.Fatal("did not save and load entity correctly")
	}
}
