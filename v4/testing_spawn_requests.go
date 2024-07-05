package sameriver

import (
	"math/rand"
)

func testingQueueSpawnSimple(em EntityManagerInterface) {
	em.QueueSpawn(nil)
}

func testingQueueSpawnUnique(em EntityManagerInterface) {
	em.QueueSpawn(map[string]any{
		"uniqueTag": "the chosen one",
	})
}

func testingSpawnUnique(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"uniqueTag": "the chosen one",
	})
}

func testingSpawnSimple(em EntityManagerInterface) *Entity {
	return em.Spawn(nil)
}

func testingSpawnPosition(
	em EntityManagerInterface, pos Vec2D) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION: pos,
		}})
}

func testingSpawnTagged(
	em EntityManagerInterface, tag string) *Entity {
	return em.Spawn(map[string]any{
		"tags": []string{tag},
	})
}

func testingSpawnSpatial(
	em EntityManagerInterface, pos Vec2D, box Vec2D) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION: pos,
			_BOX:      box,
		}})
}

func testingSpawnCollision(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION: Vec2D{10, 10},
			_BOX:      Vec2D{4, 4},
		}})
}

func testingSpawnCollisionRandom(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION: Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			_BOX:      Vec2D{5, 5},
		}})
}

func testingSpawnSteering(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION:       Vec2D{0, 0},
			_VELOCITY:       Vec2D{0, 0},
			_ACCELERATION:   Vec2D{0, 0},
			_MAXVELOCITY:    3.0,
			_MOVEMENTTARGET: Vec2D{1, 1},
			_STEER:          Vec2D{0, 0},
			_MASS:           3.0,
		}})
}

func testingSpawnPhysics(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			_POSITION:     Vec2D{10, 10},
			_VELOCITY:     Vec2D{0, 0},
			_ACCELERATION: Vec2D{0, 0},
			_BOX:          Vec2D{1, 1},
			_MASS:         3.0,
		}})
}
