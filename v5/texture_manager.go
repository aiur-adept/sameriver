package sameriver

import "github.com/veandco/go-sdl2/sdl"

type TextureManager struct {
	Textures map[string]*sdl.Texture
}

func NewTextureManager() *TextureManager {
	return &TextureManager{
		Textures: make(map[string]*sdl.Texture),
	}
}

func (tm *TextureManager) LoadTexture(renderer *sdl.Renderer, filename string, kind string) {
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
