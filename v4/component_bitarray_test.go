package sameriver

import (
	"testing"
)

func TestComponentBitArrayToString(t *testing.T) {
	w := testingWorld()
	b := w.em.components.BitArrayFromIDs([]ComponentID{_POSITION, _BOX, _GENERICTAGS})
	s := w.em.components.BitArrayToString(b)
	// TODO: check s
	Logger.Printf(s)
}
