package engine

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
)

// Provides services related to entities
type EntityManager struct {
	// Component data for entities
	ComponentsData *ComponentsDataTable
	// EntityTable stores: a list of allocated EntityTokens and a
	// list of available IDs from previous deallocations
	entityTable EntityTable
	// updated entity lists of entities with given tags
	entitiesWithTag map[string]*UpdatedEntityList
	// classes stores references to entity classes, which can be
	// retrieved by string ("crow", "turtle", "bear") in GetEntityClass()
	classes map[string]EntityClass
	// ActiveEntityListCollection is used by GetUpdatedEntityList to
	// store EntityQueryWatchers and references to UpdatedEntityLists used
	// to implement GetUpdatedEntityList
	activeEntityLists *ActiveEntityListCollection
	// used to communicate with other systems
	eventBus *EventBus
	// Channel for spawn entity requests (processed as a batch each Update())
	spawnSubscription EventChannel
	// Channel for despawn entity requests (processed as a batch each Update())
	despawnSubscription EventChannel
	// spawnMutex prevents despawn / spawn events from occurring while we
	// convert the entire EntityManager to string (expensive!)
	spawnMutex sync.Mutex
}

// Construct a new entity manager
func NewEntityManager(eventBus *EventBus) *EntityManager {
	em := &EntityManager{}
	em.ComponentsData = NewComponentsDataTable(em)
	em.activeEntityLists = NewActiveEntityListCollection(em)
	em.entitiesWithTag = make(map[string]*UpdatedEntityList)
	em.classes = make(map[string]EntityClass)
	em.eventBus = eventBus
	em.spawnSubscription = eventBus.Subscribe(
		"EntityManager::SpawnRequest",
		NewSimpleEventQuery(SPAWNREQUEST_EVENT))
	em.despawnSubscription = eventBus.Subscribe(
		"EntityManager::DespawnRequest",
		NewSimpleEventQuery(DESPAWNREQUEST_EVENT))
	return em
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update() {
	m.processDespawnChannel()
	m.processSpawnChannel()
}

// set an entity Active and notify all active entity lists
func (m *EntityManager) Activate(entity *EntityToken) {
	m.setActiveState(entity, true)
}

// set an entity inactive and notify all active entity lists
func (m *EntityManager) Deactivate(entity *EntityToken) {
	m.setActiveState(entity, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(entity *EntityToken, state bool) {
	// only act if the state is different to that which exists
	if entity.active != state {
		// start / stop logic accordingly
		m.ComponentsData.Logic[entity.ID].Active = state
		// set active state
		entity.active = state
		// notify any listening lists
		m.activeEntityLists.notifyActiveState(entity, state)
	}
}

// Get a list of entities which will be updated whenever an entity becomes
// active / inactive
func (m *EntityManager) GetUpdatedEntityList(
	q EntityQuery) *UpdatedEntityList {

	return m.activeEntityLists.GetUpdatedEntityList(q)
}

// get a previously-created UpdatedEntityList by name, or nil if does not exist
func (m *EntityManager) GetUpdatedEntityListByName(
	name string) *UpdatedEntityList {

	if list, ok := m.activeEntityLists.lists[name]; ok {
		return list
	} else {
		return nil
	}
}

// Gets the first entity with the given tag. Warns to console if the entity is
// not unique. Returns an error if the entity doesn't exist
func (m *EntityManager) UniqueTaggedEntity(tag string) (*EntityToken, error) {
	list := m.EntitiesWithTag(tag)
	if list.Length() == 0 {
		errorMsg := fmt.Sprintf("tried to fetch unique entity %s, but did "+
			"not exist", tag)
		return nil, errors.New(errorMsg)
	}
	if list.Length() > 1 {
		tagsDebug("⚠ more than one entity tagged with %s, but "+
			"GetUniqueTaggedEntity was called. This is a logic error. "+
			"Returning the first entity.", tag)
	}
	return list.FirstEntity()
}

func (m *EntityManager) EntitiesWithTag(
	tag string) *UpdatedEntityList {

	m.createEntitiesWithTagListIfNeeded(tag)
	return m.entitiesWithTag[tag]
}

func (m *EntityManager) createEntitiesWithTagListIfNeeded(tag string) {
	_, exists := m.entitiesWithTag[tag]
	if !exists {
		m.entitiesWithTag[tag] =
			m.GetUpdatedEntityList(EntityQueryFromTag(tag))
	}
}

// Boolean check of whether a given entity has a given component
func (m *EntityManager) EntityHasComponent(
	entity *EntityToken, COMPONENT int) bool {

	b, _ := m.entityTable.componentBitArrays[entity.ID].GetBit(uint64(COMPONENT))
	return b
}

// apply the given tag to the given entity
func (m *EntityManager) TagEntity(entity *EntityToken, tag string) {
	// add the tag to the taglist component
	m.ComponentsData.TagList[entity.ID].Add(tag)
	// if the entity is active, it has already been checked by all lists,
	// thus generate a new signal to add it to the list of the tag
	if entity.active {
		m.createEntitiesWithTagListIfNeeded(tag)
		m.activeEntityLists.checkActiveEntity(entity)
	}
}

// Tag each of the entities in the provided list
func (m *EntityManager) TagEntities(entities []*EntityToken, tag string) {
	for _, entity := range entities {
		m.TagEntity(entity, tag)
	}
}

// Remove a tag from an entity
func (m *EntityManager) UntagEntity(entity *EntityToken, tag string) {
	m.ComponentsData.TagList[entity.ID].Remove(tag)
	m.activeEntityLists.checkActiveEntity(entity)
}

// Remove a tag from each of the entities in the provided list
func (m *EntityManager) UntagEntities(entities []*EntityToken, tag string) {
	for _, entity := range entities {
		m.UntagEntity(entity, tag)
	}
}

// Register an entity class (subsequently retrievable)
func (m *EntityManager) AddEntityClass(c EntityClass) {
	m.classes[c.Name()] = c
}

// Get an entity class by name
func (m *EntityManager) GetEntityClass(name string) EntityClass {
	return m.classes[name]
}

// Somewhat expensive conversion of entire entity list to string, locking
// spawn/despawn from occurring while we read the entities (best to use for
// debugging, very ocassional diagnostic output)
func (m *EntityManager) String() string {
	m.spawnMutex.Lock()
	defer m.spawnMutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, entity := range m.entityTable.currentEntities {
		tags := m.ComponentsData.TagList[entity.ID]
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			entity.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
