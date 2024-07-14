package sameriver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Provides services related to entities
type EntityManager struct {
	// the world this EntityManager is inside
	w *World
	// Component data for entities
	ComponentsTable ComponentTable
	// EntityIDALlocator stores: a list of allocated Entitys and a
	// list of available IDs from previous deallocations
	EntityIDAllocator EntityIDAllocator
	// updated entity Lists created by the user according to provided filters
	Lists map[string]*UpdatedEntityList `json:"-"`
	// updated entity lists of entities with given tags
	entitiesWithTag map[string]*UpdatedEntityList
	// callbacks to call when an entity is despawned
	despawnCallbacks []func(e *Entity)
	// entities which have been tagged uniquely
	uniqueEntities map[string]*Entity
	// entities that are active
	ActiveEntities map[int]bool `json:"-"`
	// Channel for spawn entity requests (processed as a batch each Update())
	spawnSubscription *EventChannel
	// Channel for despawn entity requests (processed as a batch each Update())
	despawnSubscription *EventChannel
}

// Construct a new entity manager
func NewEntityManager(w *World) *EntityManager {
	em := &EntityManager{
		w:                   w,
		ComponentsTable:     NewComponentTable(MAX_ENTITIES),
		EntityIDAllocator:   NewEntityIDAllocator(MAX_ENTITIES),
		Lists:               make(map[string]*UpdatedEntityList),
		entitiesWithTag:     make(map[string]*UpdatedEntityList),
		uniqueEntities:      make(map[string]*Entity),
		ActiveEntities:      make(map[int]bool),
		spawnSubscription:   w.Events.Subscribe(SimpleEventFilter("spawn-request")),
		despawnSubscription: w.Events.Subscribe(SimpleEventFilter("despawn-request")),
	}
	return em
}

func (m *EntityManager) GetEntity(id int) *Entity {
	return m.EntityIDAllocator.AllocatedEntities[id]
}

func (m *EntityManager) Components() *ComponentTable {
	return &m.ComponentsTable
}

// called once per scene Update() for scenes holding an entity manager
func (m *EntityManager) Update(allowance_ms float64) float64 {
	t0 := time.Now()
	// TODO: base spawning off allowance. Spawn enough and do no more.
	m.processSpawnChannel()
	dt_ms := float64(time.Since(t0).Nanoseconds()) / 1e6
	return allowance_ms - dt_ms
}

// set an entity Active and notify all active entity lists
func (m *EntityManager) Activate(e *Entity) {
	m.setActiveState(e, true)
}

// set an entity inactive and notify all active entity lists
func (m *EntityManager) Deactivate(e *Entity) {
	m.setActiveState(e, false)
}

// sets the active state on an entity and notifies all watchers
func (m *EntityManager) setActiveState(e *Entity, state bool) {
	// only act if the state is different to that which exists
	if e.Active != state {
		if state {
			m.EntityIDAllocator.Active++
			m.ActiveEntities[e.ID] = true
		} else {
			m.EntityIDAllocator.Active--
			delete(m.ActiveEntities, e.ID)
		}
		// set active state
		e.Active = state
		// notify any listening lists
		m.notifyActiveState(e, state)
	}
}

// Gets the first entity with the given tag. Warns to console if the entity is
// not unique. Returns an error if the entity doesn't exist
func (m *EntityManager) UniqueTaggedEntity(tag string) (*Entity, error) {
	if e, ok := m.uniqueEntities[tag]; ok {
		return e, nil
	} else {
		errorMsg := fmt.Sprintf("tried to fetch unique entity %s, but did "+
			"not exist", tag)
		return nil, errors.New(errorMsg)
	}
}

func (m *EntityManager) UpdatedEntitiesWithTag(tag string) *UpdatedEntityList {
	m.createEntitiesWithTagListIfNeeded(tag)
	return m.entitiesWithTag[tag]
}

func (m *EntityManager) createEntitiesWithTagListIfNeeded(tag string) {
	if _, exists := m.entitiesWithTag[tag]; !exists {
		m.entitiesWithTag[tag] =
			m.GetUpdatedEntityList(m.w.EntityFilterFromTag(tag))
	}
}

// apply the given tags to the given entity
func (m *EntityManager) TagEntity(e *Entity, tags ...string) {
	for _, tag := range tags {
		m.w.GetTagList(e, GENERICTAGS_).Add(tag)
		if e.Active {
			m.createEntitiesWithTagListIfNeeded(tag)
		}
	}
	if e.Active {
		m.checkActiveEntity(e)
	}
}

// Tag each of the entities in the provided list
func (m *EntityManager) TagEntities(entities []*Entity, tag string) {
	for _, e := range entities {
		m.TagEntity(e, tag)
	}
}

// Remove a tag from an entity
func (m *EntityManager) UntagEntity(e *Entity, tag string) {
	m.w.GetTagList(e, GENERICTAGS_).Remove(tag)
	m.checkActiveEntity(e)
}

// Remove a tag from each of the entities in the provided list
func (m *EntityManager) UntagEntities(entities []*Entity, tag string) {
	for _, e := range entities {
		m.UntagEntity(e, tag)
	}
}

// get the maximum number of entities without a resizing and reallocating of
// components and system data (if Expand() is not a nil function for that system)
func (m *EntityManager) MaxEntities() int {
	return m.EntityIDAllocator.Capacity
}

// Get the number of allocated entities (not number of active, mind you)
func (m *EntityManager) NumEntities() (total int, active int) {
	return m.EntityIDAllocator.Allocated, m.EntityIDAllocator.Active
}

// returns a map of all active entities
func (m *EntityManager) GetActiveEntitiesSet() map[int]bool {
	return m.ActiveEntities
}

// return a map where the keys are the current entities (aka an idiomatic
// go "set")
func (m *EntityManager) GetCurrentEntitiesSet() map[int]*Entity {
	result := make(map[int]*Entity, m.EntityIDAllocator.Allocated)
	for ID, e := range m.EntityIDAllocator.AllocatedEntities {
		if e.NonNil {
			result[ID] = e
		}
	}
	return result
}

func (m *EntityManager) GetEntityByID(ID int) *Entity {
	return m.EntityIDAllocator.AllocatedEntities[ID]
}

func (m *EntityManager) ApplyComponentSet(e *Entity, spec map[ComponentID]any) {
	m.ComponentsTable.ApplyComponentSet(e, spec)
}

func (m *EntityManager) String() string {
	json, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (m *EntityManager) Save(filename string) error {
	json, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, json, 0644)
}

// dump entities with tags
func (m *EntityManager) DumpEntities() string {
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for _, e := range m.EntityIDAllocator.AllocatedEntities {
		tags := m.w.GetTagList(e, GENERICTAGS_)
		entityRepresentation := fmt.Sprintf("{id: %d, tags: %v}",
			e.ID, tags)
		buffer.WriteString(entityRepresentation)
		buffer.WriteString(",\n")
	}
	buffer.WriteString("]")
	return buffer.String()
}
