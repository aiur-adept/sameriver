package sameriver

import "github.com/veandco/go-sdl2/sdl"

type TileManager struct {
	renderer      *sdl.Renderer
	Tiles         map[string]*Tile
	Sprites       map[string]*Sprite
	TileDimension sdl.Rect
}

func NewTileManager(renderer *sdl.Renderer, tileWidth, tileHeight int32) *TileManager {
	return &TileManager{
		renderer:      renderer,
		Tiles:         make(map[string]*Tile),
		Sprites:       make(map[string]*Sprite),
		TileDimension: sdl.Rect{0, 0, tileWidth, tileHeight},
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

func (tm *TileManager) DrawTile(kind string, x, y int32, viewport *Viewport, window *sdl.Window) {
	// get the tile position relative to the viewport
	destRect := viewport.DestRect(window, x, y, tm.TileDimension.W, tm.TileDimension.H)
	tm.renderer.Copy(tm.Sprites[kind].Texture, &tm.Tiles[kind].srcRect, &destRect)
}
