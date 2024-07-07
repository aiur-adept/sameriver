package sameriver

import (
	"reflect"
	"testing"
)

func TestEntityIDAllocatorAllocateID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	et.allocateID()
}

func TestEntityIDAllocatorDeallocateID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	e := et.allocateID()
	et.deallocate(e)
	if et.Allocated != 0 {
		t.Fatal("didn't update allocated count")
	}
	if !(len(et.AvailableIDs) == 1 && et.AvailableIDs[0] == e.ID) {
		t.Fatal("didn't add deallocated ID to list of available IDs")
	}
}

func TestEntityIDAllocatorAllocateMaxIDs(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	for i := 0; i < MAX_ENTITIES; i++ {
		et.allocateID()
	}
	et.allocateID()
	et.expand(1)
	et.allocateID()
}

func TestEntityIDAllocatorReallocateDeallocatedID(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	var e *Entity
	for i := 0; i < MAX_ENTITIES; i++ {
		allocated := et.allocateID()
		if i == MAX_ENTITIES/2 {
			e = allocated
		}
	}
	et.deallocate(e)
	e = et.allocateID()
	if e.ID != MAX_ENTITIES/2 {
		t.Fatal("should have used deallocated ID to serve new allocate request")
	}
}

func TestEntityIDAllocatorSaveLoad(t *testing.T) {
	et := NewEntityIDAllocator(MAX_ENTITIES, NewIDGenerator())
	et.allocateID()
	e2 := et.allocateID()
	et.deallocate(e2)

	jsonStr := et.String()
	Logger.Println(jsonStr)

	et2 := EntityIDAllocator{}
	et2.UnmarshalJSON([]byte(jsonStr))

	if et.Capacity != et2.Capacity {
		t.Fatal("capacity didn't match")
	}
	if et.Active != et2.Active {
		t.Fatal("active didn't match")
	}
	// deep equals et and et2 AvailableIDs
	if !reflect.DeepEqual(et.AvailableIDs, et2.AvailableIDs) {
		t.Fatal("available IDs didn't match")
	}
}
