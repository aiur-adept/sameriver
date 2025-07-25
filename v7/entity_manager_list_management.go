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

// get a previously-created UpdatedEntityList by tag, or nil if does not exist
func (m *EntityManager) GetUpdatedEntityListByTag(
	tag string) *UpdatedEntityList {
	return m.getUpdatedEntityList(m.w.EntityFilterFromTag(tag), false)
}

// get a previously-created UpdatedEntityList by name, or nil if does not exist
func (m *EntityManager) GetUpdatedEntityListByName(
	name string) *UpdatedEntityList {

	if list, ok := m.Lists[name]; ok {
		return list
	} else {
		return nil
	}
}

func (m *EntityManager) GetUpdatedEntityListByComponents(names []ComponentID) *UpdatedEntityList {
	strs := make([]string, 0)
	for _, name := range names {
		strs = append(strs, m.ComponentsTable.Strings[name])
	}
	name := strings.Join(strs, ",")
	return m.GetSortedUpdatedEntityList(
		m.w.EntityFilterFromComponentBitArray(
			name,
			m.ComponentsTable.BitArrayFromIDs(names)))
}

func (m *EntityManager) getUpdatedEntityList(
	q EntityFilter, sorted bool) *UpdatedEntityList {

	// helper func that goes through already-existing entities to add them
	// to the list
	processExisting := func(q EntityFilter, list *UpdatedEntityList) {
		for _, e := range m.EntityIDAllocator.AllocatedEntities {
			if e.NonNil && q.Test(e) {
				list.Signal(EntitySignal{ENTITY_ADD, e})
			}
		}
	}

	// return the list if it already exists (this is why Filter names should
	// be unique if they expect to be unique!)
	// TODO: document this requirement
	if list, exists := m.Lists[q.Name]; exists {
		return list
	}
	// register a Filter watcher for the Filter given
	var list *UpdatedEntityList
	if sorted {
		list = NewSortedUpdatedEntityList(q.Name)
	} else {
		list = NewUpdatedEntityList(q.Name)
	}
	list.Filter = &q
	processExisting(q, list)
	m.Lists[q.Name] = list
	return list
}

// send add / remove signal to all lists according to active state of
// entity and whether its in the list
func (m *EntityManager) notifyActiveState(e *Entity, active bool) {
	for _, list := range m.Lists {
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
	for listname, list := range m.Lists {
		if list.Filter.Test(e) {
			// ensure the listname is added (if needed) to the Entity's .Lists []string
			if !e.HasList(listname) {
				e.Lists = append(e.Lists, listname)
			}
			list.Signal(EntitySignal{ENTITY_ADD, e})
		}
	}
	// check whether the entity needs to be removed from any lists it's on
	toRemove := make([]*UpdatedEntityList, 0)
	toRemoveNames := make([]string, 0)
	for _, listName := range e.Lists {
		list := m.Lists[listName]
		if list.Filter != nil && !list.Filter.Test(e) {
			toRemove = append(toRemove, list)
			toRemoveNames = append(toRemoveNames, listName)
		}
	}
	for _, list := range toRemove {
		// ensure the listname is removed
		list.Signal(EntitySignal{ENTITY_REMOVE, e})
	}
	for _, listName := range toRemoveNames {
		e.RemoveList(listName)
	}
}
