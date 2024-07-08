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
