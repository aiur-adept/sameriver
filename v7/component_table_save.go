package sameriver

import (
	"encoding/json"
	"os"

	"github.com/golang-collections/go-datastructures/bitarray"
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
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	ct.ComponentBitArrays = make([]bitarray.BitArray, ct.Capacity)
	for eid, e := range ct.ComponentStrings {
		for k, _ := range e {
			ct.orStringIntoBitArray(eid, k)
		}
	}
	return nil
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
