package sameriver

import (
	"bytes"
	"encoding/json"
)

type Entity struct {
	NonNil    bool
	ID        int
	Active    bool
	Despawned bool
	Lists     []string
	Mind      map[string]any
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
