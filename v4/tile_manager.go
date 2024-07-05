package sameriver

import (
	"encoding/json"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type TileManager struct {
	renderer      *sdl.Renderer
	Files         map[string]string
	Tiles         map[string]*Tile
	Sprites       map[string]*Sprite
	TileDimension sdl.Rect
}

func NewTileManager(renderer *sdl.Renderer) *TileManager {
	return &TileManager{
		renderer: renderer,
		Files:    make(map[string]string),
		Tiles:    make(map[string]*Tile),
		Sprites:  make(map[string]*Sprite),
	}
}

func (tm *TileManager) SetDimension(dim int32) *TileManager {
	tm.TileDimension = sdl.Rect{0, 0, dim, dim}
	return tm
}

// function to save the tile manager to a file
func (tm *TileManager) Save(filename string) {
	// write tiles to the "tiles" key in the object
	data := map[string]interface{}{
		"files":          tm.Files,
		"tile_dimension": tm.TileDimension,
	}
	obj, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, obj, 0644)
}

// function to load tiles from a file
func LoadTileManager(renderer *sdl.Renderer, filename string) *TileManager {
	// load from the JSON as defined in Save()
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var obj map[string]interface{}
	json.Unmarshal(data, &obj)

	tm := NewTileManager(renderer)
	files := obj["files"].(map[string]interface{})
	for kind, value := range files {
		filename := value.(string)
		tm.Files[kind] = filename
		tm.LoadTile(kind, filename)
	}
	dimension := obj["tile_dimension"].(map[string]interface{})
	tm.TileDimension = sdl.Rect{int32(dimension["X"].(float64)), int32(dimension["Y"].(float64)), int32(dimension["W"].(float64)), int32(dimension["H"].(float64))}
	return tm
}

// function to load a tile
func (tm *TileManager) LoadTile(kind string, filename string) {
	tm.Files[kind] = filename
	sprite := NewSprite(tm.renderer, filename)
	tm.Tiles[kind] = &Tile{
		Kind: kind,
	}
	tm.Sprites[kind] = sprite
	_, _, width, height, err := sprite.Texture.Query()
	if err != nil {
		panic(err)
	}
	tm.Tiles[kind].srcRect = sdl.Rect{0, 0, int32(width), int32(height)}
}

func (tm *TileManager) DrawTile(kind string, x, y int32, viewport *Viewport, window *sdl.Window) {
	// get the tile position relative to the viewport
	destRect := viewport.DestRect(window, x, y, tm.TileDimension.W, tm.TileDimension.H)
	tm.renderer.Copy(tm.Sprites[kind].Texture, &tm.Tiles[kind].srcRect, &destRect)
}
