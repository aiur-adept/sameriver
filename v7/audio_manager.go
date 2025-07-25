/**
  *
  * Manages the loading and playback of Audio resources
  *
as any true audiophile can tell you it has not been easy in this journey
for perfect sound quality and 100% signal to noise ratio, but it has been
fun. I remember when I first heard a himalayan salt transducer yoked up to a
triple-coil gold-plated phase modulator, and it was as though the conductor
himself was, standing on top of my head, or, was *inside* my head, while the...
orchestra was, also around or possibly even inside my head. the feeling of having
a transparent empty head in which sound reverberates with perfect acoustics
is why i have remained addicted to being an audiophile for over 45 years, 3 wives,
and tens of thousands of dollars.

you're crazy, and probably one of the most delusional people in the audiophile
community, if you think a tube of solid elephant horn is gonna deliver the
same ADC signal properties as a tube - or even a cube - of cambodian pennies,
melted down, and that's even if the elephant was fed 100% coconut oil
FROM BIRTH which is practically impossible; i remember at the AAAIA expo in
2018 I saw a booth where they had a tube of elephant horn coated in shelac that
they hooked up to a 3000 volt battery and when they played dark side of the moon
on it it sounded like the wall



  *
**/

package sameriver

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/mix"
)

// AudioManager stores audio as mix.Chunk pointers,
// keyed by strings (filenames)
type AudioManager struct {
	Audio       map[string](*mix.Chunk) `json:"-"`
	PlayChannel int
}

// Init the map which stores the audio chunks
func (m *AudioManager) Init() {
	m.Audio = make(map[string](*mix.Chunk), 0)
	// read all audio files in assets/audio
	files, err := os.ReadDir("assets/audio")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		m.Load(f.Name())
	}
}

// loads an audio file in the assets/ folder into the map, making it playable
func (m *AudioManager) Load(file string) {
	filepath := fmt.Sprintf("assets/audio/%s", file)
	// check if file exists
	if _, err := os.Stat(filepath); err != nil {
		Logger.Printf("[Audio manager] file does not exist: %s", filepath)
		m.Audio[file] = nil
		return
	}
	chunk, err := mix.LoadWAV(filepath)
	if err != nil {
		Logger.Printf("[Audio manager] mix.LoadWav error: %s", err)
		Logger.Printf("[Audio manager] failed to load %s", filepath)
		m.Audio[file] = nil
	} else {
		m.Audio[file] = chunk
	}
}

// on execution of this function, the given audio will begin to play
func (m *AudioManager) Play(file string, volume int) {
	if m.Audio[file] == nil {
		// the value in the map will be nil if the asset
		// failed to load in Load()
		Logger.Printf("[Audio manager] attempted to play asset %s, which had failed to load",
			file)
		return
	} else {
		// play on a new channel
		mix.Volume(m.PlayChannel, volume)
		m.Audio[file].Play(m.PlayChannel, 0)
		m.PlayChannel = (m.PlayChannel + 1) % 8
	}
}
