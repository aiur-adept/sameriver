package sameriver

import (
	"bytes"
	"encoding/json"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type Entity struct {
	NonNil            bool
	ID                int
	World             *World `json:"-"`
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray `json:"-"`
	Lists             []string
	Mind              map[string]any
}

func (e *Entity) GetMind(name string) any {
	if v, ok := e.Mind[name]; ok {
		return v
	}
	return nil
}

func (e *Entity) SetMind(name string, val any) {
	e.Mind[name] = val
}

func (e *Entity) HasTag(tag string) bool {
	return e.GetTagList(GENERICTAGS_).Has(tag)
}

func (e *Entity) HasTags(tags ...string) bool {
	has := true
	for _, t := range tags {
		has = has && e.GetTagList(GENERICTAGS_).Has(t)
	}
	return has
}

func (e *Entity) HasComponent(name ComponentID) bool {
	b, _ := e.ComponentBitArray.GetBit(uint64(e.World.Em.ComponentsTable.Ixs[name]))
	return b
}

func (e *Entity) HasComponents(names ...ComponentID) bool {
	has := true
	for _, name := range names {
		b, _ := e.ComponentBitArray.GetBit(uint64(e.World.Em.ComponentsTable.Ixs[name]))
		has = has && b
	}
	return has
}

func (e *Entity) DistanceFrom(x *Entity) float64 {
	ePos, eBox := e.GetVec2D(POSITION_), e.GetVec2D(BOX_)
	xPos, xBox := x.GetVec2D(POSITION_), x.GetVec2D(BOX_)
	return RectDistance(*ePos, *eBox, *xPos, *xBox)
}

func (e *Entity) DistanceFromRect(pos Vec2D, box Vec2D) float64 {
	ePos, eBox := e.GetVec2D(POSITION_), e.GetVec2D(BOX_)
	return RectDistance(*ePos, *eBox, pos, box)
}

func (e *Entity) MarshalJSON() ([]byte, error) {
	type Alias Entity
	return json.Marshal(&struct {
		*Alias
		ComponentBitArray []string `json:"ComponentBitArray"`
	}{
		Alias:             (*Alias)(e),
		ComponentBitArray: e.World.Em.ComponentsTable.BitArrayToStringArray(e.ComponentBitArray),
	})
}

func (e *Entity) String() string {
	jsonStr, _ := json.Marshal(e)
	return string(jsonStr)
}

func (e *Entity) UnmarshalJSON(data []byte) error {
	type Alias Entity
	aux := &struct {
		*Alias
		ComponentBitArray []string `json:"ComponentBitArray"`
	}{
		Alias: (*Alias)(e),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	e.ComponentBitArray = e.World.Em.ComponentsTable.BitArrayFromStrings(aux.ComponentBitArray)
	return nil
}

func EntitySliceToString(entities []*Entity) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, e := range entities {
		buf.WriteString(e.String())
		if i != len(entities)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
