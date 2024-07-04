package sameriver

import "github.com/veandco/go-sdl2/sdl"

type TileMap struct {
	Width  int32
	Height int32
	Tiles  [][]string
}

func NewTileMap(width, height int32) *TileMap {
	tmap := &TileMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]string, height),
	}
	for y := range tmap.Tiles {
		tmap.Tiles[y] = make([]string, width)
	}
	return tmap
}

func (tm *TileMap) DrawTiles(tileManager *TileManager, viewport *Viewport, window *sdl.Window) {
	startX := int32(viewport.X) / tileManager.TileDimension.W
	endX := (int32(viewport.X) + int32(viewport.Width) + tileManager.TileDimension.W - 1) / tileManager.TileDimension.W
	startY := int32(viewport.Y) / tileManager.TileDimension.H
	endY := (int32(viewport.Y) + int32(viewport.Height) + tileManager.TileDimension.H - 1) / tileManager.TileDimension.H

	for y := startY; y < endY && y < int32(len(tm.Tiles)); y++ {
		for x := startX; x < endX && x < int32(len(tm.Tiles[y])); x++ {
			tileX := x * tileManager.TileDimension.W
			tileY := y * tileManager.TileDimension.H
			if tm.Tiles[y][x] != "" {
				tileManager.DrawTile(tm.Tiles[y][x], tileX, tileY, viewport, window)
			}
		}
	}
}

func (tm *TileMap) SetTile(x, y int32, kind string) {
	tm.Tiles[y][x] = kind
}
