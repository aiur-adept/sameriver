package sameriver

import (
	"encoding/json"
)

type ItemArchetype struct {
	Name        string
	DisplayName string
	FlavourText string
	Properties  map[string]float64
	Tags        TagList
	Entity      map[string]any
}

func (i *ItemArchetype) String() string {
	b, _ := json.MarshalIndent(i, "", "\t")
	return string(b)
}

func (item *ItemArchetype) UnmarshalJSON(data []byte) error {
	// Implement custom unmarshalling logic here
	// This is a simple example that just directly unmarshals into the struct
	type Alias ItemArchetype // Create an alias to avoid recursion
	alias := &Alias{}
	if err := json.Unmarshal(data, alias); err != nil {
		return err
	}
	*item = ItemArchetype(*alias)
	return nil
}
