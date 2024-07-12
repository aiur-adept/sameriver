package sameriver

import (
	"fmt"
	"os"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type SpriteSystem struct {
	w              *World
	SpriteEntities *UpdatedEntityList

	tm         *TextureManager
	NilTexture *sdl.Texture
}

func NewSpriteSystem(renderer *sdl.Renderer, tm *TextureManager) *SpriteSystem {
	s := &SpriteSystem{}
	s.tm = tm
	s.LoadFiles(renderer)
	s.generateNilTexture(renderer)
	return s
}

func (s *SpriteSystem) GetSprite(name string, dimX, dimY int) Sprite {
	_, ok := s.tm.Textures[name]
	if !ok {
		name = "__nil_texture__"
	}
	// query dimensions of texture
	_, _, w, h, err := s.tm.Textures[name].Query()
	if err != nil {
		panic(err)
	}
	return Sprite{
		Texture: name,                   // texture
		FrameX:  0,                      // frame
		FrameY:  0,                      // frame
		FrameW:  uint8(w / int32(dimX)), // frame width
		FrameH:  uint8(h / int32(dimY)), // frame height
		DimX:    uint8(dimX),            // width
		DimY:    uint8(dimY),            // height
		Visible: true,                   // visible
		Flip:    sdl.FLIP_NONE,          // flip
	}
}

func (s *SpriteSystem) generateNilTexture(renderer *sdl.Renderer) {
	surface, err := sdl.CreateRGBSurface(
		0,          // flags
		8,          // width
		8,          //height
		int32(32),  // depth
		0xff000000, // rgba masks
		0x00ff0000,
		0x0000ff00,
		0x000000ff)
	if err != nil {
		panic(err)
	}
	rect := sdl.Rect{0, 0, 8, 8}
	color := uint32(0x9fddbcff) // feijoa
	surface.FillRect(&rect, color)
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	s.tm.Textures["__nil_texture__"] = texture
}

func (s *SpriteSystem) LoadFiles(renderer *sdl.Renderer) {
	files, err := os.ReadDir("assets/images/sprites")
	if err != nil {
		Logger.Println(err)
		logWarning("could not open assets/images/sprites; skipping SpriteSystem.LoadFiles()")
		return
	}
	for _, f := range files {
		var err error
		log_err := func(err error) {
			Logger.Printf("[Sprite manager] failed to load %s", f.Name())
			panic(err)
		}
		// get image, convert to texture, and store
		// image to texture
		surface, err := img.Load(fmt.Sprintf("assets/images/sprites/%s", f.Name()))
		if err != nil {
			log_err(err)
			continue
		}
		mapkey := strings.Split(f.Name(), ".png")[0]
		s.tm.Textures[mapkey], err = renderer.CreateTextureFromSurface(surface)
		if err != nil {
			log_err(err)
			continue
		}
		surface.Free()
	}
}

func (s *SpriteSystem) Render(renderer *sdl.Renderer, e *Entity, sprite *Sprite) {
	texture := s.tm.Textures[sprite.Texture]

	pos := s.w.GetVec2D(e, POSITION_)
	box := s.w.GetVec2D(e, BOX_)

	pos.ShiftCenterToBottomLeft(*box)
	defer pos.ShiftBottomLeftToCenter(*box)

	srcRect := sdl.Rect{
		X: int32(int32(sprite.FrameW) * int32(sprite.FrameX)),
		Y: int32(int32(sprite.FrameH) * int32(sprite.FrameY)),
		W: int32(sprite.FrameW),
		H: int32(sprite.FrameH),
	}
	destRect := sdl.Rect{
		X: int32(s.w.GetVec2D(e, POSITION_).X),
		Y: int32(s.w.GetVec2D(e, POSITION_).Y),
		W: int32(s.w.GetVec2D(e, BOX_).X),
		H: int32(s.w.GetVec2D(e, BOX_).Y),
	}
	renderer.Copy(texture, &srcRect, &destRect)
}

// System funcs

func (s *SpriteSystem) GetComponentDeps() []any {
	return []any{
		BASESPRITE_, SPRITE, "BASESPRITE",
	}
}

func (s *SpriteSystem) LinkWorld(w *World) {
	s.w = w

	s.SpriteEntities = w.GetUpdatedEntityListByComponents([]ComponentID{BASESPRITE_})
}

func (s *SpriteSystem) Update(dt_ms float64) {
	// nil?
}

func (s *SpriteSystem) Expand(n int) {
	// nil?
}
