// The collision detection in this sytem has 4 main parts:
//
// 1. a method to check collisions invoked by the game every game loop
//
// 2. an UpdatedEntityList of entities having Position and HitBox
//
// 3. a special data structure which holds rate limiters for each possible
// 	collision
//
//
//
// Datastructure (4.) - triangular rateLimiters array
//
// The rate limiters data structed is "collision-indexed", meaning it is indexed
// [i][j], where i and j are ID's and i < j. That is, each pairing of ID's
// is produced by matching each ID with all those greater than it.
//
// A collision-indexed data structure of ResettableRateLimiters
// used to avoid notifying of collisions too often. The need for this arises
// from the fact that we want to run the collision-checking logic as often as
// possible, but we don't want to send collision events at 30 times a second.
// These rate limiters rate-limit the sending of messages on a channel when we
// detect collisions, in order to save resources (internally they use a
// sync.Once which can be reset either by a natural delay or externally, in
// a goroutine-safe way)
//
//          j
//
//      0 1 2 3 4
//     0  r r r r
//     1    r r r
//  i  2      r r
//     3        r
//     4
//
package engine

type CollisionSystem struct {
	// Reference to entity manager to reach components
	entityManager *EntityManager
	// Reference to event bus to publish collisions
	eventBus *EventBus
	// targetted entities
	collidableEntities *UpdatedEntityList
	// an array of rate limiters to avoid the problem where we send out a
	// collision event every single loop. we want to check for collisions as
	// often as possible, but we don't want to send out collision events that
	// often, as it will put a load on anything reading these events
	rateLimiterArray CollisionRateLimiterArray
	// How the collision system communicates collision events
	ev *EventBus
}

func (s *CollisionSystem) Init(
	entityManager *EntityManager,
	eventBus *EventBus) {

	// take down references to em and ev
	s.entityManager = entityManager
	s.eventBus = eventBus
	// get a regularly updated list of the entities which are collidable
	// (position and hitbox)
	query := EntityQueryFromComponentBitArray(
		"collidable",
		MakeComponentBitArray([]ComponentType{
			BOX_COMPONENT}))
	s.collidableEntities = s.entityManager.GetUpdatedEntityList(query)
	// add a callback to the UpdatedEntityList of collidable entities
	// so that whenever an entity is removed, we will reset its rate limiters
	// in the collision rate limiter array (to guard against an entity
	// despawning, a new entity spawning with its ID, and failing a collision
	// test (rare prehaps, but an edge case we nonetheless want to avoid)
	s.collidableEntities.addCallback(
		func(signal EntitySignal) {
			entity := signal.entity
			if entity.ID < 0 {
				entity.ID = -(entity.ID + 1)
				s.rateLimiterArray.Reset(entity)
			}
		})
}

// Test collision between two entities
func (s *CollisionSystem) TestCollision(i uint16, j uint16) bool {
	return s.entityManager.ComponentsData.Box[i].HasIntersection(
		&s.entityManager.ComponentsData.Box[j])
}

// Iterates through the entities in the UpdatedEntityList using a handshake
// pattern, where, given a sorted list of ID's corresponding to collidable
// entities, i is compared with all ID's after i, then i + 1 is compared with
// all entities after i + 1, etc. (basically we iterate through the
// collision-indexed rate-limiter 2d triangular array row by row, left to right)
//
// If a collision is confirmed by checking their positions and bounding boxes,
// we attempt to send a collision event through the channel to be processed
// by goroutine 2 ("Event filtering and sending"), but we rate-limit sending
// events for each possible collision [i][j] using the rate limiter at [i][j]
// in rateLimiters, so if we already sent one within the timeout, we just move on.
func (s *CollisionSystem) Update(dt_ms uint16) {

	entities := s.collidableEntities.Entities

	// NOTE: The ID's in collidableEntities are in sorted order,
	// so the rateLimiterArray access condition that i < j is respected
	// check each possible collison between entities in the list by doing a
	// handshake pattern
	for ix := uint16(0); ix < uint16(len(entities)); ix++ {
		for jx := ix + 1; jx < uint16(len(entities)); jx++ {
			// get the entity ID's
			i := entities[ix]
			j := entities[jx]
			// check the collision
			if s.TestCollision(uint16(i.ID), uint16(j.ID)) {
				// if colliding, send the message (rate-limited)
				s.rateLimiterArray.
					GetRateLimiter(i.ID, j.ID).
					Do(func() {
						s.eventBus.Publish(COLLISION_EVENT,
							CollisionData{EntityA: i, EntityB: j})
					})
				// TODO: move both entities away from their common center?
				// generalized callback function probably best (with a set of
				// predefined ones)
			}
		}
	}
}
