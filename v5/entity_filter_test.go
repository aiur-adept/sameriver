package sameriver

import (
	"testing"
)

func TestEntityFilter(t *testing.T) {
	w := testingWorld()

	pos := Vec2D{0, 0}
	e := testingSpawnPosition(w, pos)
	q := EntityFilter{
		"positionFilter",
		func(e *Entity) bool {
			return *e.GetVec2D(POSITION_) == pos
		},
	}
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}

func TestEntityFilterFromTag(t *testing.T) {
	w := testingWorld()

	tag := "tag1"
	e := testingSpawnTagged(w, tag)
	q := EntityFilterFromTag(tag)
	if !q.Test(e) {
		t.Fatal("Filter did not return true")
	}
}

func TestEntityFilterFromCanBe(t *testing.T) {
	w := testingWorld()
	ox := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{0, 0},
			BOX_:      Vec2D{3, 2},
			STATE_: map[string]int{
				"yoked": 0,
			},
		},
		"tags": []string{"ox"},
	})
	q := EntityFilterFromCanBe(map[string]int{"yoked": 1})
	if !q.Test(ox) {
		t.Fatal("Should've responded to ox that can be yoked")
	}
	ox.GetIntMap(STATE_).SetValidInterval("yoked", 0, 0)
	if q.Test(ox) {
		t.Fatal("Should've failed for unyokable ox")
	}
}
