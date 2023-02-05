package sameriver

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSpatialHashInsertion(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	testData := map[[2]Vec2D][][2]int{
		[2]Vec2D{Vec2D{5, 5}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{1, 1}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{4, 4}, Vec2D{1, 1}}:   [][2]int{[2]int{0, 0}},
		[2]Vec2D{Vec2D{1, 11}, Vec2D{1, 1}}:  [][2]int{[2]int{0, 1}},
		[2]Vec2D{Vec2D{11, 11}, Vec2D{1, 1}}: [][2]int{[2]int{1, 1}},
		[2]Vec2D{Vec2D{41, 41}, Vec2D{1, 1}}: [][2]int{[2]int{4, 4}},
		[2]Vec2D{Vec2D{99, 99}, Vec2D{1, 1}}: [][2]int{[2]int{9, 9}},
		[2]Vec2D{Vec2D{11, 99}, Vec2D{1, 1}}: [][2]int{[2]int{1, 9}},
	}
	entityCells := make(map[*Entity][][2]int)
	for posbox, cells := range testData {
		e, _ := testingSpawnSpatial(w, posbox[0], posbox[1])
		entityCells[e] = cells
	}
	w.Update(FRAME_DURATION_INT / 2)
	for e, cells := range entityCells {
		for _, cell := range cells {
			inCell := false
			for _, entity := range sh.hasher.Entities(cell[0], cell[1]) {
				if entity == e {
					inCell = true
				}
			}
			if !inCell {
				t.Fatal(fmt.Sprintf("%v,%v was not mapped to cell %v",
					e.GetVec2D("Position"),
					e.GetVec2D("Box"),
					cell))
			}
		}
	}
}

func TestSpatialHashMany(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	for i := 0; i < 300; i++ {
		testingSpawnSpatial(w,
			Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			Vec2D{5, 5})
	}
	w.Update(FRAME_DURATION_INT / 2)
	n_entities := w.em.entityTable.active
	seen := make(map[*Entity]bool)
	found := 0
	table := sh.hasher.TableCopy()
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			cell := table[x][y]
			for _, e := range cell {
				if _, ok := seen[e]; !ok {
					found++
					seen[e] = true
				}
			}
		}
	}
	if found != n_entities {
		t.Fatal("Some entities were not in any cell!")
	}
}

func TestSpatialHashLargeEntity(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	pos := Vec2D{20, 20}
	box := Vec2D{5, 5}
	cells := [][2]int{
		[2]int{1, 1},
		[2]int{1, 2},
		[2]int{2, 1},
		[2]int{2, 2},
	}
	e, _ := testingSpawnSpatial(w, pos, box)
	w.Update(FRAME_DURATION_INT / 2)
	for _, cell := range cells {
		inCell := false
		for _, entity := range sh.hasher.Entities(cell[0], cell[1]) {
			if entity == e {
				inCell = true
			}
		}
		if !inCell {
			t.Fatal(fmt.Sprintf("%v,%v was not mapped to cell %v",
				e.GetVec2D("Position"),
				e.GetVec2D("Box"),
				cell))
		}
	}
}

func TestSpatialHashCellsWithinDistance(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)

	box := Vec2D{0.01, 0.01}

	// we're checking the radius at 0, 0, the corner of the world
	cells := sh.hasher.CellsWithinDistance(Vec2D{0, 0}, box, 25.0)
	if len(cells) != 8 {
		t.Fatal(fmt.Sprintf("circle centered at 0, 0 of radius 25 should touch 8 cells; got %d: %v", len(cells), cells))
	}
	cells = sh.hasher.CellsWithinDistance(Vec2D{0, 0}, box, 29.0)
	if len(cells) != 9 {
		t.Fatal(fmt.Sprintf("circle centered at 0, 0 of radius 29 should touch 9 cells; got %d: %v", len(cells), cells))
	}
	// now check from a position not quite at the corner
	cells = sh.hasher.CellsWithinDistance(Vec2D{20, 20}, box, 29.0)
	if len(cells) != 25 {
		t.Fatal(fmt.Sprintf("circle centered at 20, 20 of radius 29 should touch 25 cells; got %d: %v", len(cells), cells))
	}
	cells = sh.hasher.CellsWithinDistance(Vec2D{20, 20}, box, 7.0)
	if len(cells) != 4 {
		t.Fatal(fmt.Sprintf("circle centered at 20, 20 of radius 7 should touch 4 cells; got %d: %v", len(cells), cells))
	}
	cells = sh.hasher.CellsWithinDistance(Vec2D{25, 25}, box, 1.0)
	if len(cells) != 1 {
		t.Fatal(fmt.Sprintf("circle centered at 25, 25 of radius 1 should touch 1 cell; got %d: %v", len(cells), cells))
	}
}

/*
func TestSpatialHashEntitiesPotentiallyWithinDistance(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)

	e, _ := testingSpawnPosition(w, Vec2D{20, 20})

	entities = sh.EntitiesPotentiallyWithinDistance(e, 29.0)
	if len(entities) != 8 {
		t.Fatal(fmt.Sprintf("circle centered at 20, 20 of radius 29 should have caught 8 entities; got %d", len(entities)))
	}

}
*/

func TestSpatialHashTableCopy(t *testing.T) {
	w := NewWorld(100, 100)
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	testingSpawnSpatial(w, Vec2D{1, 1}, Vec2D{1, 1})
	w.Update(FRAME_DURATION_INT / 2)
	w.Update(FRAME_DURATION_INT / 2)
	table := sh.hasher.Table
	tableCopy := sh.hasher.TableCopy()
	if table[0][0][0] != tableCopy[0][0][0] {
		t.Fatal("CurrentTableCopy() doesn't return a copy")
	}
	table[0][0] = table[0][0][:0]
	if len(tableCopy[0][0]) == 0 {
		t.Fatal("CurrentTableCopy() doesn't return a copy")
	}
}

func TestSpatialHashTableToString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	w := testingWorld()
	sh := NewSpatialHashSystem(10, 10)
	w.RegisterSystems(sh)
	s0 := sh.hasher.String()
	for i := 0; i < 500; i++ {
		testingSpawnSpatial(w,
			Vec2D{rand.Float64() * 1024, rand.Float64() * 1024},
			Vec2D{5, 5})
	}
	w.Update(FRAME_DURATION_INT)
	s1 := sh.hasher.String()
	if len(s1) < len(s0) {
		t.Fatal("spatial hash did not show entities in its String() representation")
	}
}
