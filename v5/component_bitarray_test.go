package sameriver

import (
	"testing"
)

func TestComponentBitArrayToString(t *testing.T) {
	w := testingWorld()
	b := w.Em.components.BitArrayFromIDs([]ComponentID{POSITION_, BOX_, GENERICTAGS_})
	s := w.Em.components.BitArrayToString(b)
	// TODO: check s
	Logger.Printf(s)
}
