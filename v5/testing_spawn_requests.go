package sameriver

import (
	"math/rand"
)

func testingQueueSpawnSimple(w *World) {
	w.QueueSpawn(nil)
}

func testingQueueSpawnUnique(w *World) {
	w.QueueSpawn(map[string]any{
		"uniqueTag": "the chosen one",
	})
}

func testingSpawnUnique(w *World) *Entity {
	return w.Spawn(map[string]any{
		"uniqueTag": "the chosen one",
	})
}

func testingSpawnSimple(w *World) *Entity {
	return w.Spawn(nil)
}

func testingSpawnPosition(
	w *World, pos Vec2D) *Entity {
	return w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: pos,
		}})
}

func testingSpawnTagged(
	w *World, tag string) *Entity {
	return w.Spawn(map[string]any{
		"tags": []string{tag},
	})
}

func testingSpawnSpatial(
	w *World, pos Vec2D, box Vec2D) *Entity {
	return w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: pos,
			BOX_:      box,
		}})
}

func testingSpawnCollision(w *World) *Entity {
	return w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{10, 10},
			BOX_:      Vec2D{4, 4},
		}})
}

func testingSpawnCollisionRandom(w *World) *Entity {
	return w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			BOX_:      Vec2D{5, 5},
		}})
}

func testingSpawnSteering(w *World) *Entity {
	return w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_:       Vec2D{0, 0},
			VELOCITY_:       Vec2D{0, 0},
			ACCELERATION_:   Vec2D{0, 0},
			MAXVELOCITY_:    3.0,
			MOVEMENTTARGET_: Vec2D{1, 1},
			STEER_:          Vec2D{0, 0},
			MASS_:           3.0,
		}})
}

func testingSpawnPhysics(w *World) *Entity {
	return w.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_:     Vec2D{10, 10},
			VELOCITY_:     Vec2D{0, 0},
			ACCELERATION_: Vec2D{0, 0},
			BOX_:          Vec2D{1, 1},
			MASS_:         3.0,
			RIGIDBODY_:    true,
		}})
}
