package engine

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/golang-collections/go-datastructures/bitarray"
)

type Entity struct {
	ID                int
	World             *World
	WorldID           int
	Active            bool
	Despawned         bool
	ComponentBitArray bitarray.BitArray
	ListsMutex        sync.RWMutex
	Lists             []*UpdatedEntityList
	Logics            map[string]*LogicUnit
}

func (e *Entity) LogicUnitName(name string) string {
	return fmt.Sprintf("entity-logic-%d-%s", e.ID, name)
}

func (e *Entity) makeLogicUnit(name string, F func()) *LogicUnit {
	return &LogicUnit{
		name:    e.LogicUnitName(name),
		f:       F,
		active:  true,
		worldID: e.World.IdGen.Next(),
	}
}

func (e *Entity) AddLogic(name string, F func()) *LogicUnit {
	l := e.makeLogicUnit(name, F)
	e.Logics[name] = l
	e.World.addEntityLogic(e, l)
	return l
}

func (e *Entity) RemoveLogic(name string) {
	if _, ok := e.Logics[name]; !ok {
		panic(fmt.Sprintf("Trying to remove logic %s - but entity doesn't have it!", name))
	}
	e.World.removeEntityLogic(e, e.Logics[name])
	delete(e.Logics, name)
}

func (e *Entity) RemoveAllLogics() {
	for _, l := range e.Logics {
		e.World.removeEntityLogic(e, l)
	}
}

func (e *Entity) ActivateLogics() {
	for _, logic := range e.Logics {
		logic.active = true
	}
}

func (e *Entity) DeactivateLogics() {
	for _, logic := range e.Logics {
		logic.active = false
	}
}

func EntitySliceToString(entities []*Entity) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, e := range entities {
		buf.WriteString(fmt.Sprintf("%d", e.ID))
		if i != len(entities)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
