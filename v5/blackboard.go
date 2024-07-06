package sameriver

import (
	"encoding/json"
	"fmt"
)

type Blackboard struct {
	Name   string
	State  map[string]any
	Events *EventBus
}

func NewBlackboard(name string) *Blackboard {
	return &Blackboard{
		Name:   name,
		State:  make(map[string]any),
		Events: NewEventBus("blackboard-" + name),
	}
}

func (b *Blackboard) Has(k string) bool {
	_, ok := b.State[k]
	return ok
}

func (b *Blackboard) Get(k string) any {
	return b.State[k]
}

func (b *Blackboard) Set(k string, v any) {
	b.State[k] = v
}

func (b *Blackboard) Remove(k string) {
	delete(b.State, k)
}

// NOTE: every type in State must be marshalable to json
func (b *Blackboard) String() string {
	// marshal to json
	json, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("error marshalling blackboard: %v", err)
	}
	return string(json)
}

// NOTE: every type in State must be unmarshalable from json
func (b *Blackboard) UnmarshalJSON(data []byte) error {
	result := json.Unmarshal(data, &b.State)
	if result != nil {
		return result
	}
	b.Events = NewEventBus("blackboard-" + b.Name)
	return nil
}
