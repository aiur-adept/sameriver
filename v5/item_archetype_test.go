package sameriver

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestItemArchetypeSaveLoad(t *testing.T) {

	i := NewItemSystem(nil)

	i.CreateArchetype(map[string]any{
		"name":        "sword_iron",
		"displayName": "iron sword",
		"flavourText": "a good irons word, decently sharp",
		"properties": map[string]int{
			"damage":     3,
			"value":      20,
			"durability": 5,
		},
		"tags": []string{"weapon"},
	})

	jsonStr := i.Archetypes["sword_iron"].String()

	arch := ItemArchetype{}
	err := json.Unmarshal([]byte(jsonStr), &arch)
	if err != nil {
		t.Errorf("Error unmarshalling: %v", err)
	}
	fmt.Println(arch)
}
