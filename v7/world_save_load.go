package sameriver

import (
	"encoding/json"
	"os"
)

func (w *World) Save(filename string) {
	jsonObj, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, jsonObj, 0644)
}

func LoadWorld(filename string) *World {
	jsonObj, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	wTemp := NewWorld(nil)
	json.Unmarshal(jsonObj, wTemp)
	w := NewWorld(map[string]any{
		"seed":                wTemp.Seed,
		"width":               wTemp.Width,
		"height":              wTemp.Height,
		"distanceHasherGridX": wTemp.DistanceHasherGridX,
		"distanceHasherGridY": wTemp.DistanceHasherGridY,
	})
	json.Unmarshal(jsonObj, w)
	return w
}
