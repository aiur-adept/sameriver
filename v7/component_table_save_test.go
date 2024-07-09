package sameriver

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestComponentTableSave(t *testing.T) {

	// normal setup
	w := testingWorld()
	p := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)
	w.RegisterSystems(p, cs)
	e := testingSpawnPhysics(w)
	*w.GetVec2D(e, VELOCITY_) = Vec2D{1, 1}

	w.Em.ComponentsTable.Save("test.json")

	ct := ComponentTableFromJSON("test.json")
	Logger.Println(ct)

	// check if the component table is the same as the original
	if ct.Vec2DMap[VELOCITY_][e.ID] != *w.GetVec2D(e, VELOCITY_) {
		t.Errorf("Vec2DMap[%v][0] = %v, want %v", VELOCITY_, ct.Vec2DMap[VELOCITY_][0], *w.GetVec2D(e, VELOCITY_))
	}

	// os.Remove("test.json")
}

func TestComponentTableSaveState(t *testing.T) {

	// normal setup
	w := testingWorld()

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			STATE_: map[string]int{
				"health": 100,
			},
		},
	})

	fmt.Println(w.Em.ComponentsTable.IntMapMap[STATE_][e.ID].M)

	w.Em.ComponentsTable.Save("test.json")

	ct := ComponentTableFromJSON("test.json")

	//check if the component table is the same as the original
	if ct.IntMapMap[STATE_][e.ID].M["health"] != 100 {
		t.Errorf("IntMap[%v][0] = %v, want %v", STATE_, ct.IntMapMap[STATE_][e.ID].M["health"], 100)
	}

	os.Remove("test.json")
}

func TestComponentTableSaveTime(t *testing.T) {

	// normal setup
	w := testingWorld()

	TIME_ := GENERICTAGS_ + 1
	w.RegisterComponents(
		[]any{TIME_, TIME, "TIME_BORN"},
	)

	born := time.Now()
	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			TIME_: born,
		},
	})

	w.Em.ComponentsTable.Save("test.json")
	fmt.Println(w.Em.ComponentsTable.TimeMap[TIME_][e.ID])

	ct := ComponentTableFromJSON("test.json")

	fmt.Println(ct.TimeMap[TIME_][e.ID])

	//check if the component table is the same as the original
	if !ct.TimeMap[TIME_][e.ID].Equal(born) {
		t.Errorf("IntMap[%v][0] = %v, want %v", STATE_, ct.TimeMap[TIME_][e.ID], born)
	}

	os.Remove("test.json")
}

func TestComponentTableSaveTimeAccumulator(t *testing.T) {

	// normal setup
	w := testingWorld()

	TIME_ACCUM_ := GENERICTAGS_ + 1
	w.RegisterComponents(
		[]any{TIME_ACCUM_, TIMEACCUMULATOR, "TIME_ACCUM"},
	)

	accum := NewTimeAccumulator(1000)
	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			TIME_ACCUM_: accum,
		},
	})

	w.Em.ComponentsTable.Save("test.json")
	fmt.Println(w.Em.ComponentsTable.TimeAccumulatorMap[TIME_ACCUM_][e.ID])

	ct := ComponentTableFromJSON("test.json")

	fmt.Println(ct.TimeAccumulatorMap[TIME_ACCUM_][e.ID])

	//check if the component table is the same as the original
	if ct.TimeAccumulatorMap[TIME_ACCUM_][e.ID] != accum {
		t.Errorf("IntMap[%v][0] = %v, want %v", STATE_, ct.TimeAccumulatorMap[TIME_ACCUM_][e.ID], accum)
	}

	os.Remove("test.json")
}

func TestComponentTableSaveBitArrays(t *testing.T) {

	// normal setup
	w := testingWorld()

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			STATE_: map[string]int{
				"health": 100,
			},
		},
	})

	fmt.Println(w.Em.ComponentsTable.IntMapMap[STATE_][e.ID].M)

	w.Em.ComponentsTable.Save("test.json")

	ct := ComponentTableFromJSON("test.json")

	fmt.Println(ct.ComponentBitArrays[e.ID])

	os.Remove("test.json")
}
