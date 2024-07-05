package sameriver

import (
	"encoding/json"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type TileMap struct {
	tm     *TileManager
	Width  int32
	Height int32
	Tiles  [][]string
}

func NewTileMap(tm *TileManager, width, height int32) *TileMap {
	tmap := &TileMap{
		tm:     tm,
		Width:  width,
		Height: height,
		Tiles:  make([][]string, height),
	}
	for y := range tmap.Tiles {
		tmap.Tiles[y] = make([]string, width)
	}
	return tmap
}

func (tm *TileMap) DrawTiles(viewport *Viewport, window *sdl.Window) {

	startX := int32(viewport.X) / tm.tm.Dimension
	endX := (int32(viewport.X) + int32(viewport.Width) + tm.tm.Dimension - 1) / tm.tm.Dimension
	startY := int32(viewport.Y) / tm.tm.Dimension
	endY := (int32(viewport.Y) + int32(viewport.Height) + tm.tm.Dimension - 1) / tm.tm.Dimension

	for y := startY; y < endY && y < int32(len(tm.Tiles)); y++ {
		for x := startX; x < endX && x < int32(len(tm.Tiles[y])); x++ {
			tileX := x * tm.tm.Dimension
			tileY := y * tm.tm.Dimension
			if tm.Tiles[y][x] != "" {
				tm.tm.DrawTile(tm.Tiles[y][x], tileX, tileY, viewport, window)
			}
		}
	}
}

func (tm *TileMap) SetTile(x, y int32, kind string) {
	tm.Tiles[y][x] = kind
}

func (tm *TileMap) Save(filename string) {
	// save to json
	data := map[string]interface{}{
		"tile_manager": tm,
		"width":        tm.Width,
		"height":       tm.Height,
		"tiles":        tm.Tiles,
	}
	obj, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, obj, 0644)
}

// func LoadTileMap(renderer *sdl.Renderer, filename string) *TileMap {
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var obj map[string]interface{}
// 	err = json.Unmarshal(data, &obj)
// 	if err != nil {
// 		panic(err)
// 	}
// 	width := int32(obj["width"].(float64))
// 	height := int32(obj["height"].(float64))
// 	tmap := NewTileMap(tm, width, height)
// 	tmap.Tiles = obj["tiles"].([][]string)
// 	return tmap
// }
