package sameriver

import (
	"encoding/json"
	"fmt"
)

type Entity struct {
	NonNil     bool
	ID         int
	Active     bool
	Despawned  bool
	Lists      []string
	Mind       Blackboard
	Components []string
}

func (e *Entity) String() string {
	jsonStr, _ := json.Marshal(e)
	return string(jsonStr)
}

func (e *Entity) HasList(listName string) bool {
	for _, list := range e.Lists {
		if list == listName {
			return true
		}
	}
	return false
}

func (e *Entity) RemoveList(listName string) bool {
	for i, list := range e.Lists {
		if list == listName {
			e.Lists = append(e.Lists[:i], e.Lists[i+1:]...)
			return true
		}
	}
	return false
}

func (e *Entity) LogicUnitName(name string) string {
	return fmt.Sprintf("%d-%s", e.ID, name)
}
