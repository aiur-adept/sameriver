package sameriver

import (
	"encoding/json"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type TileManager struct {
	renderer  *sdl.Renderer
	Files     map[string]string  `json:"files"`
	Tiles     map[string]*Tile   `json:"-"`
	Sprites   map[string]*Sprite `json:"-"`
	Dimension int32              `json:"dimension"`
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
	tm.Dimension = dim
	return tm
}

// function to save the tile manager to a file
func (tm *TileManager) Save(filename string) {
	// marshal to jon
	obj, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		panic(err)
	}
	os.WriteFile(filename, obj, 0644)
}

func TileManagerFromJSON(renderer *sdl.Renderer, filename string) *TileManager {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var tm TileManager
	json.Unmarshal(data, &tm)
	tm.renderer = renderer
	tm.Tiles = make(map[string]*Tile)
	tm.Sprites = make(map[string]*Sprite)
	tm.LoadTiles()
	return &tm
}

func (tm *TileManager) LoadTiles() {
	for kind, filename := range tm.Files {
		tm.LoadTile(kind, filename)
	}
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
	destRect := viewport.DestRect(window, x, y, tm.Dimension, tm.Dimension)
	tm.renderer.Copy(tm.Sprites[kind].Texture, &tm.Tiles[kind].srcRect, &destRect)
}
