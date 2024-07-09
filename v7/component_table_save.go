package sameriver

import (
	"encoding/json"
	"os"
)

func (ct *ComponentTable) String() string {
	str, err := json.MarshalIndent(ct, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (ct *ComponentTable) UnmarshalJSON(data []byte) error {
	type Alias ComponentTable
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(ct),
	}
	return json.Unmarshal(data, &aux)
}

func (ct *ComponentTable) Save(filename string) {
	str := ct.String()
	os.WriteFile(filename, []byte(str), 0644)
}

func ComponentTableFromJSON(filename string) *ComponentTable {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var ct ComponentTable
	err = ct.UnmarshalJSON(data)
	if err != nil {
		panic(err)
	}
	return &ct
}
