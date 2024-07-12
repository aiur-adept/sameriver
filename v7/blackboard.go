package sameriver

import (
	"encoding/json"
	"fmt"
)

type Blackboard struct {
	Name   string
	State  map[string]any
	Ints   map[string]bool
	Events *EventBus `json:"-"`
}

func NewBlackboard(name string) Blackboard {
	return Blackboard{
		Name:   name,
		State:  make(map[string]any),
		Ints:   make(map[string]bool),
		Events: NewEventBus("blackboard-" + name),
	}
}

func (b Blackboard) Has(k string) bool {
	_, ok := b.State[k]
	return ok
}

func (b Blackboard) Get(k string) any {
	return b.State[k]
}

func (b Blackboard) Set(k string, v any) {
	// cast to float if v is of type int
	if _, ok := v.(int); ok {
		b.State[k] = float64(v.(int))
		b.Ints[k] = true
	} else if _, ok := v.([]int); ok {
		b.State[k] = v
		b.Ints[k] = true
	} else {
		b.State[k] = v
	}
}

func (b Blackboard) Remove(k string) {
	delete(b.State, k)
	if _, ok := b.Ints[k]; ok {
		delete(b.Ints, k)
	}
}

func (bb *Blackboard) UnmarshalJSON(data []byte) error {
	var aux struct {
		Name  string
		State map[string]interface{}
		Ints  map[string]bool
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bb.Name = aux.Name
	bb.State = make(map[string]any)

	for key, value := range aux.State {
		switch v := value.(type) {
		case []interface{}:
			isStr, isFloat, isBool := false, false, false
			strSlice := make([]string, len(v))
			floatSlice := make([]float64, len(v))
			boolSlice := make([]bool, len(v))
			for i, item := range v {
				if str, ok := item.(string); ok {
					strSlice[i] = str
					isStr = true
				} else if float, ok := item.(float64); ok {
					floatSlice[i] = float
					isFloat = true
				} else if bool, ok := item.(bool); ok {
					boolSlice[i] = bool
					isBool = true
				} else {
					return fmt.Errorf("value under key %s is not a valid-typed slice", key)
				}
			}
			if isStr {
				bb.State[key] = strSlice
			} else if isFloat {
				bb.State[key] = floatSlice
			} else if isBool {
				bb.State[key] = boolSlice
			}
		case interface{}:
			// unmarshal string, int, float or bool
			if str, ok := v.(string); ok {
				bb.State[key] = str
			} else if f, ok := v.(float64); ok {
				bb.State[key] = f
			} else if b, ok := v.(bool); ok {
				bb.State[key] = b
			} else {
				return fmt.Errorf("value under key %s is not a valid-typed value", key)
			}
		}
	}

	for i, _ := range aux.Ints {
		if value, ok := bb.State[i].([]float64); ok {
			ints := make([]int, len(value))
			for j, item := range value {
				ints[j] = int(item)
			}
			bb.State[i] = ints
		}
		if value, ok := bb.State[i].(float64); ok {
			bb.State[i] = int(value)
		}
	}

	return nil
}

func (b *Blackboard) MarshalJSON() ([]byte, error) {
	type Alias Blackboard
	aux := &struct {
		State map[string]interface{}
		*Alias
	}{
		State: b.State,
		Alias: (*Alias)(b),
	}
	return json.MarshalIndent(aux, "", "  ")
}
