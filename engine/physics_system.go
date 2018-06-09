/**
  *
  *
  *
  *
**/

package engine

type PhysicsSystem struct {
	// to filter, lookup entities
	em *EntityManager
	// targetted entities
	physicsEntities *UpdatedEntityList
}

func (s *PhysicsSystem) Init(em *EntityManager) {
	// take down a reference to entity manager
	s.em = em
	// get a regularly updated list of the entities which have physics
	// (position, velocity and hitbox)
	query := EntityQueryFromComponentBitArray(
		"physical",
		MakeComponentBitArray([]ComponentType{
			BOX_COMPONENT,
			VELOCITY_COMPONENT}))
	s.physicsEntities = s.em.GetUpdatedEntityList(query)
}

// apply velocity to position of entities
// NOTE: this is called from Update and is covered by its mutex on the
// components
func (s *PhysicsSystem) applyPhysics(entity EntityToken, dt_ms uint16) {
	// read the position and velocity, using dt to compute dx, dy
	pos := s.em.Components.Position.Data[entity.ID]
	vel := s.em.Components.Velocity.Data[entity.ID]
	dx := int16(vel[0] * float32(dt_ms/4))
	dy := int16(vel[1] * float32(dt_ms/4))
	box := s.em.Components.HitBox.Data[entity.ID]
	// prevent from leaving the world in X
	if pos[0]+dx <
		int16(box[0]/2) {
		pos[0] = int16(box[0] / 2)
	} else if pos[0]+dx >
		int16(WORLD_WIDTH)-int16(box[0]/2) {
		pos[0] = int16(WORLD_WIDTH) - int16(box[0]/2)
	} else {
		pos[0] += dx
	}
	// prevent from leaving the world in Y
	if pos[1]+dy <
		int16(box[1]/2) {
		pos[1] = int16(box[1] / 2)
	} else if pos[1]+dy >
		int16(WORLD_HEIGHT)-int16(box[1]/2) {
		pos[1] = int16(WORLD_HEIGHT) - int16(box[1]/2)
	} else {
		pos[1] += dy
	}
	// set the new position which has been computed
	s.em.Components.Position.Data[entity.ID] = pos
}

func (s *PhysicsSystem) Update(dt_ms uint16) {

	s.physicsEntities.Mutex.Lock()
	defer s.physicsEntities.Mutex.Unlock()

	for _, e := range s.physicsEntities.Entities {
		// apply the physics only if this entity isn't already locked
		// (atomic operations are cheap, so this isn't a bad thing to
		// do for each entity during each Update())
		if s.em.attemptLockEntityOnce(e) {
			s.applyPhysics(e, dt_ms)
			s.em.releaseEntity(e)
		}
	}
}
