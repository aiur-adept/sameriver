package sameriver

import (
	"testing"
)

func TestInvalidComponentType(t *testing.T) {
	w := testingWorld()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Should panic if given unregistered component id")
		}
	}()
	var UNREGISTEREDCOMPONENT ComponentID = 1337
	w.Em.ComponentsTable.makeComponentSet(map[ComponentID]any{
		UNREGISTEREDCOMPONENT: Vec2D{0, 0},
	})
}

func TestComponentSetToBitArray(t *testing.T) {
	w := testingWorld()
	b := w.Em.ComponentsTable.BitArrayFromComponentSet(map[ComponentID]any{
		POSITION_: Vec2D{0, 0},
	})
	// TODO: convert to proper string and actually test
	Logger.Println(b)
}

func TestComponentSetApply(t *testing.T) {
	w := testingWorld()
	e := testingSpawnSimple(w)
	l := NewTagList()
	cs := map[ComponentID]any{
		GENERICTAGS_: l,
	}
	w.Em.ComponentsTable.ApplyComponentSet(e, cs)
	eb := w.Em.ComponentsTable.ComponentBitArrays[e.ID]
	csb := w.Em.ComponentsTable.BitArrayFromComponentSet(cs)
	Logger.Println(eb)
	Logger.Println(csb)
	if !eb.Equals(csb) {
		t.Fatal("failed to apply componentset according to bitarray")
	}
}
