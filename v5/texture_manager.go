package sameriver

import (
	"encoding/json"
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
)

type TextureManager struct {
	Files    map[string]string
	Textures map[string]*sdl.Texture `json:"-"`
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		Files:    make(map[string]string),
		Textures: make(map[string]*sdl.Texture),
	}
}

func (tm *TextureManager) Init(renderer *sdl.Renderer) {
	// LoadTexture for each file in assets/textures/
	files, err := filepath.Glob("assets/textures/*.bmp")
	if err != nil {
		panic(err)
	}
	for _, filename := range files {
		kind := filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))])
		tm.loadTexture(renderer, filename, kind)
	}
}

func (tm *TextureManager) LoadFiles(renderer *sdl.Renderer) {
	for kind, filename := range tm.Files {
		tm.loadTexture(renderer, filename, kind)
	}
}

func (tm *TextureManager) loadTexture(renderer *sdl.Renderer, filename string, kind string) {
	surface, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(err)
	}
	// create texture from surface
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	tm.Textures[kind] = texture
	tm.Files[kind] = filename
}

func (tm *TextureManager) CreateSprite(kind string) *Sprite {
	return &Sprite{
		Texture: kind,
		Visible: true,
	}
}

func (tm *TextureManager) Render(renderer *sdl.Renderer, kind string, x, y, w, h int32) {
	_, _, width, height, err := tm.Textures[kind].Query()
	if err != nil {
		panic(err)
	}
	srcRect := sdl.Rect{0, 0, width, height}
	destRect := sdl.Rect{x, y, w, h}
	renderer.Copy(tm.Textures[kind], &srcRect, &destRect)
}

func (tm *TextureManager) String() string {
	// marshal to json
	json, err := json.MarshalIndent(tm, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(json)
}

func (tm *TextureManager) UnmarshalJSON(data []byte) error {
	type Alias TextureManager
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(tm),
	}
	return json.Unmarshal(data, &aux)
}
