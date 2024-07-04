package sameriver

import "github.com/veandco/go-sdl2/sdl"

type TileManager struct {
	renderer *sdl.Renderer
	Tiles    map[string]*Tile
	Sprites  map[string]*Sprite
}

func NewTileManager(renderer *sdl.Renderer) *TileManager {
	return &TileManager{
		renderer: renderer,
		Tiles:    make(map[string]*Tile),
		Sprites:  make(map[string]*Sprite),
	}
}

// function to load a tile
func (tm *TileManager) LoadTile(kind string, filename string) {
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

func (tm *TileManager) DrawTile(kind string, x, y, w, h int32) {
	destRect := sdl.Rect{x, y, w, h}
	tm.renderer.Copy(tm.Sprites[kind].Texture, &tm.Tiles[kind].srcRect, &destRect)
}
