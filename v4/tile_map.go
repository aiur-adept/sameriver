package sameriver

import "github.com/veandco/go-sdl2/sdl"

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

	startX := int32(viewport.X) / tm.tm.TileDimension.W
	endX := (int32(viewport.X) + int32(viewport.Width) + tm.tm.TileDimension.W - 1) / tm.tm.TileDimension.W
	startY := int32(viewport.Y) / tm.tm.TileDimension.H
	endY := (int32(viewport.Y) + int32(viewport.Height) + tm.tm.TileDimension.H - 1) / tm.tm.TileDimension.H

	for y := startY; y < endY && y < int32(len(tm.Tiles)); y++ {
		for x := startX; x < endX && x < int32(len(tm.Tiles[y])); x++ {
			tileX := x * tm.tm.TileDimension.W
			tileY := y * tm.tm.TileDimension.H
			if tm.Tiles[y][x] != "" {
				tm.tm.DrawTile(tm.Tiles[y][x], tileX, tileY, viewport, window)
			}
		}
	}
}

func (tm *TileMap) SetTile(x, y int32, kind string) {
	tm.Tiles[y][x] = kind
}
