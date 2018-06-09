/**

The collision detection in this sytem works using a method invoked by the
game every game loop, a goroutine running to read collision events from a
buffered channel (basically a queue which is easy for a goroutine to read),
an UpdatedEntityList of entities having Position and HitBox, and a special
data structure which holds rate limiters for each possible collision.

Collision detection is called after the entity manager has (de)spawned and
(de)activated entities, after entity-component modifications send by
various logic goroutines have been processed, and after the physics update,
so the list of entities in the collidableEntities UpdatedEntityList is valid,
their positions and hitboxes will not change, and none of them will despawn.



Datastructure: rateLimiters

The rate limiters data structed is "collision-indexed", meaning it is indexed
[i][j], where i and j are ID's and i < j. That is, each pairing of ID's
is produced by matching each ID with all those greater than it.

A collision-indexed data structure of ResettableRateLimiters
used to avoid notifying of collisions too often. The need for this arises
from the fact that we want to run the collision-checking logic as often as
possible, but we don't want to send collision events at 30 times a second.
These rate limiters rate-limit the sending of messages on a channel when we
detect collisions, in order to save resources (internally they use a
sync.Once which can be reset either by a natural delay or externally, in
a goroutine-safe way)

         j

     0 1 2 3 4
    0  r r r r
    1    r r r
 i  2      r r
    3        r
    4



Method: "Detection"

Iterates through the entities in the UpdatedEntityList using a handshake
pattern, where, given a sorted list of ID's corresponding to collidable
entities, i is compared with all ID's after i, then i + 1 is compared with
all entities after i + 1, etc. (basically we iterate through the
collision-indexed rate-limiter 2d triangular array row by row, left to right)

If a collision is confirmed by checking their positions and bounding boxes,
we attempt to send a collision event through the channel to be processed
by goroutine 2 ("Event filtering and sending"), but we rate-limit sending
events for each possible collision [i][j] using the rate limiter at [i][j]
in rateLimiters, so if we already sent one within the timeout, we just move on.



Goroutine: "Event filtering and sending"

This goroutine reads collision events from the buffered channel which the
"Detection" method wrote to, comparing each event to a list of tests supplied
by scubscribers to collision events. If a collision event matches a test
for a subscriber, it get sent to the channel of their CollisionQueryWatcher.

**/

package engine

import (
	"fmt"
)

type CollisionSystem struct {
	// Reference to entity manager to reach components
	em *EntityManager
	// targetted entities
	collidableEntities *UpdatedEntityList
	// an array of rate limiters to avoid the problem where we send out a
	// collision event every single loop. we want to check for collisions as
	// often as possible, but we don't want to send out collision events that
	// often, as it will put a load on anything reading these events
	// (used by the event-checker loop)
	rateLimiterArray CollisionRateLimiterArray
	// How the collision system communicates collision events
	ev *EventBus
}

func (s *CollisionSystem) Init(
	em *EntityManager,
	ev *EventBus) {

	// take down references to em and ev
	s.em = em
	s.ev = ev
	// get a regularly updated list of the entities which are collidable
	// (position and hitbox)
	query := EntityQueryFromComponentBitArray(
		"collidable",
		MakeComponentBitArray([]ComponentType{
			BOX_COMPONENT}))
	s.collidableEntities = s.em.GetUpdatedEntityList(query)
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
				s.rateLimiterArray.ResetAll(entity)
			}
		})
}

// Test collision between two entities
func (s *CollisionSystem) TestCollision(i uint16, j uint16) bool {
	return s.em.Components.Box[i].HasIntersection(&s.em.Components.Box[j])
}

func (s *CollisionSystem) Update(dt_ms uint16) {

	// prevent any updates to the collidableEntities list while we're using it
	s.collidableEntities.Mutex.Lock()
	defer s.collidableEntities.Mutex.Unlock()

	// acquire exclusive lock on the box component (position and bounding box)
	// TODO: have this happen at a higher level of abstraction - see comment
	// in spatial_hash.go in ComputeSpatialHash()
	s.em.Components.accessLocks[BOX_COMPONENT].Lock()
	defer s.em.Components.accessLocks[BOX_COMPONENT].Unlock()

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
				s.rateLimiterArray.GetRateLimiter(
					uint16(i.ID),
					uint16(j.ID)).Do(func() {
					s.ev.Publish(Event{
						Type:        COLLISION_EVENT,
						Description: fmt.Sprintf("collision(%d,%d)", i, j),
						Data: CollisionEvent{
							EntityA: i,
							EntityB: j}})
				})
				// TODO: move both entities away from their common center
			}
		}
	}
}
