package sameriver

import "encoding/json"

// used by the EntityManager to hold info about the allocated entities
type EntityIDAllocator struct {
	// the ID Generator given by the world the entity manager is in
	IdGen *IDGenerator
	// list of available entity ID's which have previously been deallocated
	AvailableIDs []int
	// map of allocated entities
	AllocatedEntities map[int]*Entity
	// storage for the actual entity objects
	Entities []Entity
	// how many entities are allocated
	Allocated int
	// how many entities are Active
	Active int
	// Capacity of how many ID's we can allocate without expanding
	Capacity int
}

func NewEntityIDAllocator(capacity int, IDGen *IDGenerator) *EntityIDAllocator {
	return &EntityIDAllocator{
		IdGen:             IDGen,
		Entities:          make([]Entity, capacity),
		AllocatedEntities: make(map[int]*Entity),
		Capacity:          capacity,
	}
}

func (a *EntityIDAllocator) expand(n int) {
	a.Capacity += n
	a.Entities = append(a.Entities, make([]Entity, n)...)
}

// get the ID for a new e. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (a *EntityIDAllocator) allocateID() *Entity {

	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	var ID int
	n_avail := len(a.AvailableIDs)
	if n_avail > 0 {
		// there is an ID available for a previously deallocated e.
		// pop it from the list and continue with that as the ID
		ID = a.AvailableIDs[n_avail-1]
		a.AvailableIDs = a.AvailableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		ID = a.Allocated
	}
	entity := Entity{ID: ID, NonNil: true}
	a.Entities[ID] = entity
	a.AllocatedEntities[ID] = &entity
	a.Allocated++
	return &entity
}

func (a *EntityIDAllocator) deallocate(e *Entity) {
	// guards against false deallocation (edge case, but hey)
	if a.Entities[e.ID].NonNil {
		a.AvailableIDs = append(a.AvailableIDs, e.ID)
		a.Entities[e.ID] = Entity{ID: e.ID, NonNil: false}
		delete(a.AllocatedEntities, e.ID)
		a.Allocated--
	}
}

func (a *EntityIDAllocator) String() string {
	jsonStr, err := json.Marshal(a)
	if err != nil {
		Logger.Println(err)
	}
	return string(jsonStr)
}

func (a *EntityIDAllocator) UnmarshalJSON(data []byte) error {
	type Alias EntityIDAllocator
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}
	err := json.Unmarshal(data, &aux)
	return err
}
