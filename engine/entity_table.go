package engine

import (
	"errors"
	"fmt"

	"github.com/dt-rush/sameriver/engine/utils"
)

// used by the EntityManager to hold info about the allocated entities
type EntityTable struct {
	// the ID Generator given by the world the entity manager is in
	IDGen *utils.IDGenerator
	// how many entities there are
	n int
	// how many entities are active
	active int
	// list of Entities which have been allocated
	currentEntities []*EntityToken
	// isAllocated is maintained for quick deallocate guard lookups
	isAllocated [MAX_ENTITIES]bool
	// list of available entity ID's which have previously been deallocated
	availableIDs []int
}

func NewEntityTable(IDGen *utils.IDGenerator) *EntityTable {
	return &EntityTable{IDGen: IDGen}
}

// get the ID for a new entity. Only called by SpawnEntity, which locks
// the entityTable, so it's safe that this method operates on that data.
// Returns int32 so that we can return -1 in case we have run out of space
// to spawn entities
func (t *EntityTable) allocateID() (*EntityToken, error) {
	// if maximum entity count reached, fail with message
	if t.n == MAX_ENTITIES {
		msg := fmt.Sprintf("Reached max entity count: %d. "+
			"Will not allocate ID.", MAX_ENTITIES)
		Logger.Println(msg)
		return nil, errors.New(msg)
	}
	// if there is a deallocated entity somewhere in the table before the
	// highest ID, return that ID to the caller
	n_avail := len(t.availableIDs)
	var ID int
	if n_avail > 0 {
		// there is an ID available for a previously deallocated entity.
		// pop it from the list and continue with that as the ID
		ID = t.availableIDs[n_avail-1]
		t.availableIDs = t.availableIDs[:n_avail-1]
	} else {
		// every slot in the table before the highest ID is filled
		ID = t.n
	}
	t.isAllocated[ID] = true
	// Increment the entity count
	t.n++
	// return the token
	entity := EntityToken{
		ID:        ID,
		WorldID:   t.IDGen.Next(),
		Active:    false,
		Despawned: false,
	}
	return &entity, nil
}

func (t *EntityTable) deallocate(e *EntityToken) {
	// guards against false deallocation (edge case, but hey)
	if t.isAllocated[e.ID] {
		t.n--
		t.availableIDs = append(t.availableIDs, e.ID)
		t.IDGen.Free(e.WorldID)
	}
}
