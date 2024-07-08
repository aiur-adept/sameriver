package sameriver

import (
	"encoding/json"
	"math"
)

type IDGenerator struct {
	Universe map[int]bool
	Freed    map[int]bool
	X        int
}

func NewIDGenerator() IDGenerator {
	return IDGenerator{
		Universe: make(map[int]bool),
		Freed:    make(map[int]bool),
	}
}

func (g *IDGenerator) Next() (ID int) {
	// try to get ID from already-available freed IDs
	if len(g.Freed) > 0 {
		// get first of freed (break immediately)
		for freeID := range g.Freed {
			ID = freeID
			delete(g.Freed, freeID)
			break
		}
	} else {
		// if there are no free id's, we're chock-full up to the latest
		// value of x++
		g.X++
		ID = g.X
		if ID > math.MaxUint32/64 {
			panic("tried to generate more than (2^32 - 1) / 64 simultaneous " +
				"ID's without free. This is surely a logic error. If you're" +
				"from the future and you can run 4,294,967,295 entities..." +
				"well, that's wild")
		}
	}
	g.Universe[ID] = true
	return ID
}

func (g *IDGenerator) Free(ID int) {
	delete(g.Universe, ID)
	g.Freed[ID] = true
}

func (g *IDGenerator) String() string {
	jsonStr, err := json.Marshal(g)
	if err != nil {
		Logger.Println(err)
	}
	return string(jsonStr)
}

func (g *IDGenerator) UnmarshalJSON(data []byte) error {
	type Alias IDGenerator
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(g),
	}
	return json.Unmarshal(data, &aux)
}
