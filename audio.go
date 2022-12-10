// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "sync"

const channelCount = 4
const maxChannelNo = channelCount - 1

// Audio returns AudioSystem object.
func Audio() *AudioSystem {
	return audio
}

var audio = &AudioSystem{
	effects: make([]SoundEffect, 256),
}

// AudioSystem provides methods for generating audio stream.
// AudioSystem contains 256 sound effects, starting at 0, ending at 255.
// AudioSystem has 4 channels. Each channel is able to play one sound effect
// at a time.
//
// AudioSystem is safe to use in a concurrent manner.
type AudioSystem struct {
	mutex         sync.Mutex
	effects       []SoundEffect
	effectsPlayed [4]effectPlayed
}

type effectPlayed struct {
	playing     bool
	effect      byte
	currentNote byte
	noteStart   byte
	noteEnd     byte
}

// Set sets the SoundEffect with given number.
func (s *AudioSystem) Set(number byte, e SoundEffect) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.effects[number] = e
}

// Get returns a copy of SoundEffect with given number. Modifying the returned SoundEffect
// will not alter AudioSystem. After making changes you must run [AudioSystem.Set] to update
// the sound effect.
func (s *AudioSystem) Get(number byte) SoundEffect {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.effects[number]
}

type SoundEffect struct {
	Speed byte // 1 is the fastest (~8.33ms). For 120 the length of one note becomes 1 second. 0 means SoundEffect takes no time.
	Notes [32]Note
}

type Note struct {
	Pitch  byte
	Volume byte // 0 - quiet. 255 - loudest. 255 values is way too much!
}

// Read method is used by back-end to read generated audio stream and play it back to the user. The sample rate is 44100,
// 16 bit depth and stereo (2 audio channels).
//
// Read is (usually) executed concurrently with main game loop. Back-end could decide about buffer size, although
// the higher the size the higher the lag. Usually the buffer is 8KB, which is 46ms of audio.
func (s *AudioSystem) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for i := 0; i < len(p); i++ {
		p[i] = 0
	}

	for _, e := range s.effectsPlayed {
		if !e.playing {
			continue
		}

		for i := 0; i < len(p); i++ {
			p[i] = byte(i)
		}
	}

	return len(p), nil
}

func (s *AudioSystem) Play(soundEffect, channelNo, noteStart, noteEnd byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if channelNo > maxChannelNo {
		return
	}

	s.effectsPlayed[channelNo] = effectPlayed{
		playing:     true,
		effect:      soundEffect,
		currentNote: 0,
		noteStart:   0,
		noteEnd:     0,
	}
}

func (s *AudioSystem) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.effects = make([]SoundEffect, 256)
	s.effectsPlayed = [4]effectPlayed{}
}

func Sfx(soundEffect byte) {
	audio.Play(soundEffect, 0, 0, 31)
}
