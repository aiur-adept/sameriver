package sameriver

// User facing function which is used to drain the state of the
// entity manager, and will also kill any pending spawn requests
func (m *EntityManager) DespawnAll() {
	// drain the spawn request channel of pending spawns
	for len(m.spawnSubscription.C) > 0 {
		<-m.spawnSubscription.C
	}
	// iterate all entities which have been allocated and despawn them
	for _, e := range m.GetCurrentEntitiesSet() {
		m.Despawn(e)
	}
}

// Despawn an entity
func (m *EntityManager) Despawn(e *Entity) {
	// guard against multiple logics per tick despawning an entity
	if !e.Despawned {
		e.Despawned = true
		m.EntityIDAllocator.deallocate(e)
		m.setActiveState(e, false)
		for _, cb := range m.despawnCallbacks {
			cb(e)
		}
	}
}

func (m *EntityManager) AddDespawnCallback(cb func(e *Entity)) {
	m.despawnCallbacks = append(m.despawnCallbacks, cb)
}
