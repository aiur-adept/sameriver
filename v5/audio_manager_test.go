package sameriver

import (
	"testing"
	"time"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

func TestAudioManagerInitAndPlay(t *testing.T) {
	if err := sdl.Init(sdl.INIT_AUDIO); err != nil {
		Logger.Println(err)
		return
	}
	defer sdl.Quit()

	mix.Init(mix.INIT_MP3)

	if err := mix.OpenAudio(48000, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		Logger.Println(err)
		return
	}
	defer sdl.CloseAudio()

	sdl.PauseAudio(false)

	// Create an instance of AudioManager
	manager := AudioManager{}

	// Initialize AudioManager
	manager.Init()

	// Attempt to play the "bell.wav" sound
	manager.Play("bell.wav")
	time.Sleep(1000 * time.Millisecond)
	manager.Play("bell.wav")
	time.Sleep(1000 * time.Millisecond)
	manager.Play("bell.wav")
	time.Sleep(1000 * time.Millisecond)
}
