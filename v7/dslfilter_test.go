package sameriver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSLBasic(t *testing.T) {

	//
	// PREPARE WORLD
	//
	w := testingWorld()
	ps := NewPhysicsSystem()
	cs := NewCollisionSystem(FRAME_DURATION / 2)

	items := NewItemSystem(nil)
	inventories := NewInventorySystem()
	w.RegisterSystems(ps, cs, items, inventories)

	items.CreateArchetype(map[string]any{
		"name":        "yoke",
		"displayName": "a yoke for cattle",
		"flavourText": "one of mankind's greatest inventions... an ancestral gift!",
		"properties": map[string]int{
			"value": 25,
		},
		"tags": []string{"item.agricultural"},
		"entity": map[string]any{
			"sprite": "yoke",
			"box":    [2]float64{0.2, 0.2},
		},
	})

	e := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_:  Vec2D{0, 0},
			BOX_:       Vec2D{1, 1},
			INVENTORY_: inventories.Create(nil),
		},
	})

	yoke := items.CreateItemSimple("yoke")
	items.SpawnItemEntity(Vec2D{0, 5}, yoke)

	positions := []Vec2D{
		Vec2D{0, 10}, // close ox
		Vec2D{0, 30}, // far ox
	}
	tags := []string{
		"a",
		"b",
	}

	oxen := make([]*Entity, len(positions))
	for i := 0; i < len(positions); i++ {
		oxen[i] = w.Spawn(map[string]any{
			"components": map[ComponentID]any{
				POSITION_: positions[i],
				BOX_:      Vec2D{3, 2},
				STATE_: map[string]int{
					"yoked": 0,
				},
			},
			"tags": []string{"ox", tags[i]},
		})
	}
	field := w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{0, 100},
			BOX_:      Vec2D{30, 30},
			STATE_: map[string]int{
				"tilled": 0,
			},
		},
		"tags": []string{"field"},
	})
	w.CreateBlackboard("somebb").Set("field", field.ID)
	// TEST
	//

	Logger.Println("1")
	entities, err := w.EFDSLFilterEntity(e, "HasTag(ox)")
	assert.NoError(t, err)
	assert.ElementsMatch(t, oxen, entities)

	Logger.Println("2")
	entities, err = w.EFDSLFilterEntity(e, "HasComponent(position)")
	assert.NoError(t, err)
	assert.Equal(t, 5, len(entities)) // e, item.yoke, 2 oxen, field

	Logger.Println("3")
	entities, err = w.EFDSLFilterEntity(e, "WithinDistance(self, 15)")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(entities)) // e, item.yoke, close ox

	Logger.Println("4")
	w.GetIntMap(oxen[0], STATE_).Set("yoked", 1)
	entities, err = w.EFDSLFilterEntity(e, "State(yoked, 1)")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))
	assert.Equal(t, oxen[0], entities[0])

	// Test Entity.DSLFilterSort
	Logger.Println("5")
	filterSortEntities, err := w.EFDSLFilterSortEntity(e, "HasTag(ox); Closest(self)")
	assert.NoError(t, err)
	assert.Equal(t, oxen[0], filterSortEntities[0]) // a close
	assert.Equal(t, oxen[1], filterSortEntities[1]) // b far

	//
	// world
	//

	// Test World.DSLFilter
	Logger.Println("6")
	worldEntities, err := w.EFDSLFilter("HasTag(ox)")
	assert.NoError(t, err)
	assert.ElementsMatch(t, oxen, worldEntities)
	// Test World.DSLFilterSort
	Logger.Println("7")
	worldFilterSortEntities, err := w.EFDSLFilterSort("HasTag(ox); Closest(bb.somebb.field)")
	assert.NoError(t, err)
	assert.Equal(t, oxen[1], worldFilterSortEntities[0]) // b close
	assert.Equal(t, oxen[0], worldFilterSortEntities[1]) // a far
}
