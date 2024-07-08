package sameriver

import (
	"regexp"
	"testing"
	"time"
)

func TestWorldSpawn(t *testing.T) {

	w := testingWorld()

	pos := Vec2D{11, 11}
	e := testingSpawnPosition(w, pos)
	if *w.GetVec2D(e, POSITION_) != pos {
		t.Fatal("failed to apply component data")
	}
	total, _ := w.NumEntities()
	if total == 0 {
		t.Fatal("failed to Spawn simple Spawn request entity")
	}
	if !(e.Despawned == false &&
		e.Active == true) {
		t.Fatal("entity struct not populated properly")
	}
	if len(w.GetCurrentEntitiesSet()) == 0 {
		t.Fatal("entity not added to current entities set on spawn")
	}
}

func TestWorldSpawnFail(t *testing.T) {
	w := testingWorld()

	for i := 0; i < MAX_ENTITIES; i++ {
		testingSpawnSimple(w)
	}
	testingSpawnSimple(w)
	if len(w.Components().Vec2DMap[POSITION_]) == MAX_ENTITIES {
		t.Fatal("Did not resize component data tables")
	}
}

func TestWorldQueueSpawn(t *testing.T) {
	w := testingWorld()
	testingQueueSpawnSimple(w)
	// sleep long enough for the event to appear on the channel
	time.Sleep(FRAME_DURATION)
	w.Update(FRAME_MS / 2)
	total, _ := w.NumEntities()
	if total != 1 {
		t.Fatal("should have spawned an entity after processing spawn " +
			"request channel")
	}
}

func TestWorldDespawn(t *testing.T) {
	w := testingWorld()
	// fill up the spawnSubscription channel
	for i := 0; i < EVENT_SUBSCRIBER_CHANNEL_CAPACITY; i++ {
		testingQueueSpawnSimple(w)
	}
	// spawn two more entities (one simple, one unique)
	testingQueueSpawnSimple(w)
	testingQueueSpawnUnique(w)
	// sleep long enough for the events to appear on the channel
	time.Sleep(FRAME_DURATION)
	// update *twice*, allowing the extra events to process despite having seen
	// a full spawn subscription channel the first time
	w.Update(FRAME_MS / 2)
	w.Update(FRAME_MS / 2)
	total, _ := w.NumEntities()
	if total != EVENT_SUBSCRIBER_CHANNEL_CAPACITY+2 {
		t.Fatal("should have spawned entities after processing spawn " +
			"request channel")
	}
}

func TestWorldDespawnAll(t *testing.T) {
	w := testingWorld()
	for i := 0; i < 64; i++ {
		testingSpawnSimple(w)
	}
	for i := 0; i < 64; i++ {
		testingQueueSpawnSimple(w)
	}
	w.DespawnAll()
	total, _ := w.NumEntities()
	if total != 0 {
		t.Fatal("did not despawn all entities")
	}
	w.Update(FRAME_MS / 2)
	total, _ = w.NumEntities()
	if total != 0 {
		t.Fatal("DespawnAll() did not discard pending spawn requests")
	}
}

func TestWorldEntityHasComponent(t *testing.T) {
	w := testingWorld()
	pos := Vec2D{11, 11}
	e := testingSpawnPosition(w, pos)
	if !w.EntityHasComponent(e, POSITION_) {
		t.Fatal("failed to set or get entity component bit array")
	}
}

func TestWorldEntitiesWithTag(t *testing.T) {
	w := testingWorld()
	tag := "tag1"
	testingSpawnTagged(w, tag)
	tagged := w.UpdatedEntitiesWithTag(tag)
	if tagged.Length() == 0 {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
}

func TestWorldSpawnUnique(t *testing.T) {
	w := testingWorld()

	uniqueTag := "the chosen one"
	e, err := w.UniqueTaggedEntity(uniqueTag)
	if !(e == nil && err != nil) {
		t.Fatal("should return err if unique entity not found")
	}
	e = testingSpawnUnique(w)

	eRetrieved, err := w.UniqueTaggedEntity(uniqueTag)
	if !(eRetrieved == e && err == nil) {
		t.Fatal("did not return unique entity")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic (aragorn voice: are you frightened? not nearly frightened enough)")
		}
	}()

	testingSpawnUnique(w)
}

func TestWorldTagUntagEntity(t *testing.T) {
	w := testingWorld()
	e := testingSpawnSimple(w)
	tag := "tag1"
	w.TagEntity(e, tag)
	tagged := w.UpdatedEntitiesWithTag(tag)
	empty := tagged.Length() == 0
	if empty {
		t.Fatal("failed to find Spawned entity in EntitiesWithTag")
	}
	if !w.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag() saw entity as untagged")
	}
	w.UntagEntity(e, tag)
	empty = tagged.Length() == 0
	if !empty {
		t.Fatal("entity was still in EntitiesWithTag after untag")
	}
	if w.EntityHasTag(e, tag) {
		t.Fatal("EntityHasTag() saw entity as still having removed tag")
	}
}

func TestWorldTagEntities(t *testing.T) {
	w := testingWorld()
	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		e := testingSpawnSimple(w)
		entities = append(entities, e)
	}
	w.TagEntities(entities, tag)
	for _, e := range entities {
		if !w.EntityHasTag(e, tag) {
			t.Fatal("entity's taglist was not modified by TagEntities")
		}
	}
}

func TestWorldUntagEntities(t *testing.T) {
	w := testingWorld()
	entities := make([]*Entity, 0)
	tag := "tag1"
	for i := 0; i < 32; i++ {
		e := testingSpawnTagged(w, tag)
		entities = append(entities, e)
	}
	w.UntagEntities(entities, tag)
	for _, e := range entities {
		if w.EntityHasTag(e, tag) {
			t.Fatal("entity's taglist was not modified by UntagEntities")
		}
	}
}

func TestWorldDeactivateActivateEntity(t *testing.T) {
	w := testingWorld()
	e := testingSpawnSimple(w)
	tag := "tag1"
	w.TagEntity(e, tag)
	tagged := w.UpdatedEntitiesWithTag(tag)
	w.Deactivate(e)
	if tagged.Length() != 0 {
		t.Fatal("entity was not removed from list after Deactivate()")
	}
	_, active := w.NumEntities()
	if active != 0 {
		t.Fatal("didn't update active count")
	}
	w.Activate(e)
	if tagged.Length() != 1 {
		t.Fatal("entity was not reinserted to list after Activate()")
	}
	_, active = w.NumEntities()
	if active != 1 {
		t.Fatal("didn't update active count")
	}
}

func TestWorldGetUpdatedEntityList(t *testing.T) {
	w := testingWorld()
	name := "ILoveDya"
	nameToo := "ILoveDya!!!"
	list := w.GetUpdatedEntityList(
		NewEntityFilter(
			name,
			func(e *Entity) bool {
				return true
			}),
	)
	testingSpawnSimple(w)
	if list.Length() != 1 {
		t.Fatal("failed to update UpdatedEntityList")
	}
	list2 := w.GetUpdatedEntityList(
		NewEntityFilter(
			nameToo,
			func(e *Entity) bool {
				return true
			}),
	)
	if list2.Length() != 1 {
		t.Fatal("failed to created UpdatedEntityList relative to existing entities")
	}
}

func TestWorldGetSortedUpdatedEntityList(t *testing.T) {
	w := testingWorld()
	list := w.GetSortedUpdatedEntityList(
		NewEntityFilter(
			"filter",
			func(e *Entity) bool {
				return true
			}),
	)
	e8 := &Entity{ID: 8, Active: true, Despawned: false}
	e0 := &Entity{ID: 0, Active: true, Despawned: false}
	list.Signal(EntitySignal{ENTITY_ADD, e8})
	list.Signal(EntitySignal{ENTITY_ADD, e0})
	first, _ := list.FirstEntity()
	if first.ID != 0 {
		t.Fatal("didn't insert in order")
	}
}

func TestWorldGetUpdatedEntityListByName(t *testing.T) {
	w := testingWorld()
	name := "ILoveDya"
	if w.GetUpdatedEntityListByName(name) != nil {
		t.Fatal("should return nil if not found")
	}
	list := w.GetUpdatedEntityList(
		NewEntityFilter(
			name,
			func(e *Entity) bool {
				return false
			}),
	)
	if w.GetUpdatedEntityListByName(name) != list {
		t.Fatal("GetUpdatedEntityListByName did not find list")
	}
}

func TestWorldGetCurrentEntitiesSet(t *testing.T) {
	w := testingWorld()
	if !(len(w.GetCurrentEntitiesSet()) == 0) {
		t.Fatal("initially, len(GetCurrentEntitiesSet()) should be 0")
	}
	e := testingSpawnSimple(w)
	if !(len(w.GetCurrentEntitiesSet()) == 1) {
		t.Fatal("after spawn, len(GetCurrentEntitiesSet()) should be 1")
	}
	w.Deactivate(e)
	if !(len(w.GetCurrentEntitiesSet()) == 1) {
		t.Fatal("after deactivate, len(GetCurrentEntitiesSet()) should be 1")
	}
	w.Despawn(e)
	if !(len(w.GetCurrentEntitiesSet()) == 0) {
		t.Fatal("after despawn, len(GetCurrentEntitiesSet()) should be 0")
	}
}

func TestWorldString(t *testing.T) {
	w := testingWorld()
	if w.String() == "" {
		t.Fatal("string implementation cannot be empty string")
	}
}

func TestWorldDumpEntities(t *testing.T) {
	w := testingWorld()
	e := testingSpawnSimple(w)
	tag := "tag1"
	w.TagEntity(e, tag)
	s := w.DumpEntities()
	if ok, _ := regexp.MatchString("tag", s); !ok {
		t.Fatal("tag data for each entity wasn't produced in EntityManager.Dump()")
	}
}
