package sameriver

import (
	"encoding/json"
)

type Entity struct {
	NonNil     bool
	ID         int
	Active     bool
	Despawned  bool
	Lists      []string
	Mind       map[string]any
	Components []string
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

func (e *Entity) String() string {
	jsonStr, _ := json.Marshal(e)
	return string(jsonStr)
}
