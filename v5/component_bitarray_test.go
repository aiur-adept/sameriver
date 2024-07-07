package sameriver

import (
	"testing"
)

func TestComponentBitArrayToString(t *testing.T) {
	w := testingWorld()
	b := w.Em.ComponentsTable.BitArrayFromIDs([]ComponentID{POSITION_, BOX_, GENERICTAGS_})
	s := w.Em.ComponentsTable.BitArrayToString(b)
	// TODO: check s
	Logger.Printf(s)
}
