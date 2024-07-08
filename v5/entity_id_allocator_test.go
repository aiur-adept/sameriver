package sameriver

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestEntityIDAllocatorAllocateID(t *testing.T) {
	eia := NewEntityIDAllocator(MAX_ENTITIES)
	eia.allocateID()
}

func TestEntityIDAllocatorDeallocateID(t *testing.T) {
	eia := NewEntityIDAllocator(MAX_ENTITIES)
	e := eia.allocateID()
	eia.deallocate(e)
	if eia.Allocated != 0 {
		t.Fatal("didn't update allocated count")
	}
	if !(len(eia.AvailableIDs) == 1 && eia.AvailableIDs[0] == e.ID) {
		t.Fatal("didn't add deallocated ID to list of available IDs")
	}
}

func TestEntityIDAllocatorAllocateMaxIDs(t *testing.T) {
	eia := NewEntityIDAllocator(MAX_ENTITIES)
	for i := 0; i < MAX_ENTITIES; i++ {
		eia.allocateID()
	}
	eia.expand(1)
	eia.allocateID()
}

func TestEntityIDAllocatorReallocateDeallocatedID(t *testing.T) {
	eia := NewEntityIDAllocator(MAX_ENTITIES)
	var e *Entity
	for i := 0; i < MAX_ENTITIES; i++ {
		allocated := eia.allocateID()
		if i == MAX_ENTITIES/2 {
			e = allocated
		}
	}
	eia.deallocate(e)
	e = eia.allocateID()
	if e.ID != MAX_ENTITIES/2 {
		t.Fatal("should have used deallocated ID to serve new allocate request")
	}
}

func TestEntityIDAllocatorSaveLoad(t *testing.T) {
	w := testingWorld()

	w.Em.EntityIDAllocator.allocateID()
	e2 := w.Em.EntityIDAllocator.allocateID()
	w.Em.EntityIDAllocator.deallocate(e2)

	jsonStr, err := json.MarshalIndent(w.Em.EntityIDAllocator, "", "  ")
	if err != nil {
		t.Fatal("error marshalling entity ID allocator")
	}
	Logger.Println(string(jsonStr))

	et2 := EntityIDAllocator{}
	et2.UnmarshalJSON([]byte(jsonStr))

	if w.Em.EntityIDAllocator.Capacity != et2.Capacity {
		t.Fatal("capacity didn't match")
	}
	if w.Em.EntityIDAllocator.Active != et2.Active {
		t.Fatal("active didn't match")
	}
	// deep equals et and et2 AvailableIDs
	if !reflect.DeepEqual(w.Em.EntityIDAllocator.AvailableIDs, et2.AvailableIDs) {
		t.Fatal("available IDs didn't match")
	}
}
