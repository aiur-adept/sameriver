package sameriver

import (
	"testing"
)

func TestEntityManagerConstruct(t *testing.T) {
	em := NewEntityManager(testingWorld())
	if em == nil {
		t.Fatal("Could not construct NewEntityManager()")
	}
}

func TestEntityManagerSpawn(t *testing.T) {
	EntityManagerInterfaceTestSpawn(testingWorld().Em, t)
}

func TestWorldSpawn(t *testing.T) {
	EntityManagerInterfaceTestSpawn(testingWorld(), t)
}

func TestEntityManagerSpawnFail(t *testing.T) {
	EntityManagerInterfaceTestSpawnFail(testingWorld().Em, t)
}

func TestWorldSpawnFail(t *testing.T) {
	EntityManagerInterfaceTestSpawnFail(testingWorld(), t)
}

func TestEntityManagerQueueSpawn(t *testing.T) {
	EntityManagerInterfaceTestQueueSpawn(testingWorld().Em, t)
}

func TestWorldQueueSpawn(t *testing.T) {
	EntityManagerInterfaceTestQueueSpawn(testingWorld(), t)
}

func TestEntityManagerDespawn(t *testing.T) {
	EntityManagerInterfaceTestDespawn(testingWorld().Em, t)
}

func TestWorldDespawn(t *testing.T) {
	EntityManagerInterfaceTestDespawn(testingWorld(), t)
}

func TestEntityManagerDespawnAll(t *testing.T) {
	EntityManagerInterfaceTestDespawnAll(testingWorld().Em, t)
}

func TestWorldDespawnAll(t *testing.T) {
	EntityManagerInterfaceTestDespawnAll(testingWorld(), t)
}

func TestEntityManagerEntityHasComponent(t *testing.T) {
	EntityManagerInterfaceTestEntityHasComponent(testingWorld().Em, t)
}

func TestWorldEntityHasComponent(t *testing.T) {
	EntityManagerInterfaceTestEntityHasComponent(testingWorld(), t)
}

func TestEntityManagerEntitiesWithTag(t *testing.T) {
	EntityManagerInterfaceTestEntitiesWithTag(testingWorld().Em, t)
}

func TestWorldEntitiesWithTag(t *testing.T) {
	EntityManagerInterfaceTestEntitiesWithTag(testingWorld(), t)
}

func TestEntityManagerSpawnUnique(t *testing.T) {
	EntityManagerInterfaceTestSpawnUnique(testingWorld().Em, t)
}

func TestWorldSpawnUnique(t *testing.T) {
	EntityManagerInterfaceTestSpawnUnique(testingWorld(), t)
}

func TestEntityManagerTagUntagEntity(t *testing.T) {
	EntityManagerInterfaceTestTagUntagEntity(testingWorld().Em, t)
}

func TestWorldTagUntagEntity(t *testing.T) {
	EntityManagerInterfaceTestTagUntagEntity(testingWorld(), t)
}

func TestEntityManagerTagEntities(t *testing.T) {
	EntityManagerInterfaceTestTagEntities(testingWorld().Em, t)
}

func TestWorldTagEntities(t *testing.T) {
	EntityManagerInterfaceTestTagEntities(testingWorld(), t)
}

func TestEntityManagerUntagEntities(t *testing.T) {
	EntityManagerInterfaceTestUntagEntities(testingWorld().Em, t)
}

func TestWorldUntagEntities(t *testing.T) {
	EntityManagerInterfaceTestUntagEntities(testingWorld(), t)
}

func TestEntityManagerDeactivateActivateEntity(t *testing.T) {
	EntityManagerInterfaceTestDeactivateActivateEntity(testingWorld().Em, t)
}

func TestWorldDeactivateActivateEntity(t *testing.T) {
	EntityManagerInterfaceTestDeactivateActivateEntity(testingWorld(), t)
}

func TestEntityManagerGetUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityList(testingWorld().Em, t)
}

func TestWorldGetUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityList(testingWorld(), t)
}

func TestEntityManagerGetSortedUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetSortedUpdatedEntityList(testingWorld().Em, t)
}

func TestWorldGetSortedUpdatedEntityList(t *testing.T) {
	EntityManagerInterfaceTestGetSortedUpdatedEntityList(testingWorld(), t)
}

func TestEntityManagerGetUpdatedEntityListByName(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityListByName(testingWorld().Em, t)
}

func TestWorldGetUpdatedEntityListByName(t *testing.T) {
	EntityManagerInterfaceTestGetUpdatedEntityListByName(testingWorld(), t)
}

func TestEntityManagerGetCurrentEntitiesSet(t *testing.T) {
	EntityManagerInterfaceTestGetCurrentEntitiesSet(testingWorld().Em, t)
}

func TestWorldGetCurrentEntitiesSet(t *testing.T) {
	EntityManagerInterfaceTestGetCurrentEntitiesSet(testingWorld(), t)
}

func TestEntityManagerString(t *testing.T) {
	EntityManagerInterfaceTestString(testingWorld().Em, t)
}

func TestWorldString(t *testing.T) {
	EntityManagerInterfaceTestString(testingWorld(), t)
}

func TestEntityManagerDumpEntities(t *testing.T) {
	EntityManagerInterfaceTestDumpEntities(testingWorld().Em, t)
}

func TestWorldDumpEntities(t *testing.T) {
	EntityManagerInterfaceTestDumpEntities(testingWorld(), t)
}
