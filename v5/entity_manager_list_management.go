package sameriver

import "strings"

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(q EntityFilter) *UpdatedEntityList {
	return m.getUpdatedEntityList(q, false)
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive, sorted by ID
func (m *EntityManager) GetSortedUpdatedEntityList(
	q EntityFilter) *UpdatedEntityList {
	return m.getUpdatedEntityList(q, true)
}

// get a previously-created UpdatedEntityList by name, or nil if does not exist
func (m *EntityManager) GetUpdatedEntityListByName(
	name string) *UpdatedEntityList {

	if list, ok := m.lists[name]; ok {
		return list
	} else {
		return nil
	}
}

func (m *EntityManager) GetUpdatedEntityListByComponents(names []ComponentID) *UpdatedEntityList {
	strs := make([]string, 0)
	for _, name := range names {
		strs = append(strs, m.components.Strings[name])
	}
	name := strings.Join(strs, ",")
	return m.GetSortedUpdatedEntityList(
		EntityFilterFromComponentBitArray(
			name,
			m.components.BitArrayFromIDs(names)))
}

func (m *EntityManager) getUpdatedEntityList(
	q EntityFilter, sorted bool) *UpdatedEntityList {

	// helper func that goes through already-existing entities to add them
	// to the list
	processExisting := func(q EntityFilter, list *UpdatedEntityList) {
		for e := range m.entityIDAllocator.currentEntities {
			if q.Test(e) {
				list.Signal(EntitySignal{ENTITY_ADD, e})
			}
		}
	}

	// return the list if it already exists (this is why Filter names should
	// be unique if they expect to be unique!)
	// TODO: document this requirement
	if list, exists := m.lists[q.Name]; exists {
		return list
	}
	// register a Filter watcher for the Filter given
	var list *UpdatedEntityList
	if sorted {
		list = NewSortedUpdatedEntityList()
	} else {
		list = NewUpdatedEntityList()
	}
	list.Filter = &q
	processExisting(q, list)
	m.lists[q.Name] = list
	return list
}

// send add / remove signal to all lists according to active state of
// entity and whether its in the list
func (m *EntityManager) notifyActiveState(e *Entity, active bool) {
	for _, list := range m.lists {
		if list.Filter.Test(e) {
			if active {
				list.Signal(EntitySignal{ENTITY_ADD, e})
			} else {
				list.Signal(EntitySignal{ENTITY_REMOVE, e})
			}
		}
	}
}

// check if the entity needs to be added to or removed from any lists
func (m *EntityManager) checkActiveEntity(e *Entity) {
	for _, list := range m.lists {
		if list.Filter.Test(e) {
			list.Signal(EntitySignal{ENTITY_ADD, e})
		}
	}
	// check whether the entity needs to be removed from any lists it's on
	toRemove := make([]*UpdatedEntityList, 0)
	for _, list := range e.Lists {
		if list.Filter != nil && !list.Filter.Test(e) {
			toRemove = append(toRemove, list)
		}
	}
	for _, list := range toRemove {
		list.Signal(EntitySignal{ENTITY_REMOVE, e})
	}
}
