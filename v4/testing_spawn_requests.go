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
			POSITION_: pos,
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
			POSITION_: pos,
			BOX_:      box,
		}})
}

func testingSpawnCollision(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{10, 10},
			BOX_:      Vec2D{4, 4},
		}})
}

func testingSpawnCollisionRandom(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_: Vec2D{100 * rand.Float64(), 100 * rand.Float64()},
			BOX_:      Vec2D{5, 5},
		}})
}

func testingSpawnSteering(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
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

func testingSpawnPhysics(em EntityManagerInterface) *Entity {
	return em.Spawn(map[string]any{
		"components": map[ComponentID]any{
			POSITION_:     Vec2D{10, 10},
			VELOCITY_:     Vec2D{0, 0},
			ACCELERATION_: Vec2D{0, 0},
			BOX_:          Vec2D{1, 1},
			MASS_:         3.0,
			RIGIDBODY_:    true,
		}})
}
