package sameriver

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestIDGeneratorUnique(t *testing.T) {
	IDGen := NewIDGenerator()
	IDs := make(map[int]bool)
	for i := 0; i < 1024*1024; i++ {
		ID := IDGen.Next()
		if _, ok := IDs[ID]; ok {
			t.Fatal("produced same ID already produced")
		}
		IDs[ID] = true
	}
}

func TestIDGeneratorUniqueRemoval(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	IDGen := NewIDGenerator()
	IDs := make(map[int]bool)
	for i := 0; i < 1024*1024; i++ {
		if i > 1024 && i%2 == 0 {
			ID := rand.Intn(len(IDGen.Universe))
			IDGen.Free(ID)
			delete(IDs, ID)
		} else {
			ID := IDGen.Next()
			if _, ok := IDs[ID]; ok {
				t.Fatal("produced same ID already produced")
			}
			IDs[ID] = true
		}
	}
}

func TestIDGeneratorSaveLoad(t *testing.T) {
	IDGen := NewIDGenerator()
	IDGen.Next()
	IDGen.Next()
	IDGen.Next()

	jsonStr := IDGen.String()
	IDGen2 := IDGenerator{}
	err := IDGen2.UnmarshalJSON([]byte(jsonStr))
	if err != nil {
		t.Fatal(err)
	}

	// check deep equals on IDGen and IDGen2
	if !reflect.DeepEqual(IDGen.Universe, IDGen2.Universe) {
		t.Fatal("Universe maps didn't match")
	}
	if !reflect.DeepEqual(IDGen.Freed, IDGen2.Freed) {
		t.Fatal("Freed maps didn't match")
	}
	if IDGen.X != IDGen2.X {
		t.Fatal("X didn't match")
	}
}
