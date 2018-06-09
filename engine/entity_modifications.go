/*
 * AtomicEntityModify or AtomicEntitiesModify, that is, which produce functions
 * of type func(EntityToken) or func([]EntityToken)
 * Special case for functions which operate directly using the EntityToken,
 * they simply *are* the function which would be output, to prevent patterns
 * like DespawnAtomic()(entity), which just looks cluttered
 *
 * To make this clear, they are always suffixed with "Atomic", so that a function
 * suffixed "Atomic" should appear out of place if not inside an instance of
 * Atomic.+Modify()
 *
 */

package engine

// despawn function (this is how entities *must* be despawned, other than
// DespawnAll() which should not happen during normal gameplay)
// TODO: figure out whether DespawnAll works properly and doesn't race with
// anything
func (m *EntityManager) DespawnAtomic(entity EntityToken) {

	// NOTE: we don't need to gen validate here because we can only
	// despawn, changing gen, once all current modifications (by the logic of
	// AtomicEntityModify are valid for the gen of the entity) are RUnlock()'d,
	// and we can only call DespawnAtomic inside an AtomicEntityModify for
	// which gen matched
	if m.entityTable.despawnFlags[entity.ID].CAS(0, 1) {
		go func() {
			// lock the activeModification lock as a *writer* (compare with
			// the lock as RLock() in lockEntityComponent), so that calls to
			// AtomicEntityModify which want to lock the active modification
			// lock will not proceed until the lock is released in
			// EntityManager.processDespawnChannel() (at which point they will
			// proceed, but immediately fail the genValidate() call, since the
			// entity was despawned, causing the AtomicEntityModify call for
			// the despawned entity to return flase)
			m.activeModificationLocks[entity.ID].Lock()
			m.despawnChannel <- entity
		}()
	}
}

// set an entity Active and notify all active entity lists
func (m *EntityManager) Activate(entity EntityToken) {
	m.setActiveState(entity, true)
}

// set an entity inactive and notify all active entity lists
func (m *EntityManager) Deactivate(entity EntityToken) {
	m.setActiveState(entity, false)
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntityAtomic(tag string) func(EntityToken) {
	return func(entity EntityToken) {

		// add the tag to the taglist component
		m.Components.TagList.Data[entity.ID].Add(tag)
		// if the entity is active, it has already been checked by all lists,
		// thus generate a new signal to add it to the list of the tag
		if m.Components.Active.Data[entity.ID] {
			m.tags.createEntitiesWithTagListIfNeeded(tag)
			m.activeEntityLists.checkActiveEntity(entity)
		}
	}
}

// remove a tag from an entity
func (m *EntityManager) UntagEntityAtomic(tag string) func(EntityToken) {
	return func(entity EntityToken) {
		m.Components.TagList.Data[entity.ID].Remove(tag)
		m.activeEntityLists.checkActiveEntity(entity)
	}
}

// Tag each of the entities in the provided array of ID's with the given tag
func (m *EntityManager) TagEntitiesAtomic(tag string) func([]EntityToken) {
	return func(entities []EntityToken) {

		for _, entity := range entities {
			m.TagEntityAtomic(tag)(entity)
		}
	}
}

// Apply a given component set to the entity
func (em *EntityManager) ApplyComponentSetAtomic(
	cs ComponentSet) func(EntityToken) {

	return func(entity EntityToken) {
		if cs.Color != nil {
			em.Components.Color.Data[entity.ID] = *cs.Color
		}
		if cs.Health != nil {
			em.Components.Health.Data[entity.ID] = *cs.Health
		}
		if cs.HitBox != nil {
			em.Components.HitBox.Data[entity.ID] = *cs.HitBox
		}
		if cs.Logic != nil {
			em.Components.Logic.Data[entity.ID] = *cs.Logic
		}
		if cs.Mind != nil {
			em.Components.Mind.Data[entity.ID] = *cs.Mind
		}
		if cs.Position != nil {
			em.Components.Position.Data[entity.ID] = *cs.Position
		}
		if cs.Sprite != nil {
			em.Components.Sprite.Data[entity.ID] = *cs.Sprite
		}
		if cs.TagList != nil {
			em.Components.TagList.Data[entity.ID] = *cs.TagList
		}
		if cs.Velocity != nil {
			em.Components.Velocity.Data[entity.ID] = *cs.Velocity
		}
	}
}

// a function which generates an entity modification to augment the health
// of an entity
func (m *EntityManager) IncrementHealthAtomic(change uint8) func(EntityToken) {
	return func(e EntityToken) {
		healthNow := m.Components.Health.Data[e.ID] + change
		if healthNow > 255 {
			healthNow = 255
		} else if healthNow < 0 {
			healthNow = 0
		}
		m.Components.Health.Data[e.ID] = healthNow
	}
}
