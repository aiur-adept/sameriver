package sameriver

import (
	"encoding/json"
	"os"
)

func (ct *ComponentTable) Save(filename string) {
	json, err := json.Marshal(ct)
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, json, 0644)
}

func ComponentTableFromJSON(filename string) *ComponentTable {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var ct ComponentTable
	err = json.Unmarshal(data, &ct)
	if err != nil {
		panic(err)
	}
	return &ct
}
