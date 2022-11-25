package pi

import (
	"math"
	"sync"
)

func Audio() *AudioSystem {
	return audio
}

var audio = &AudioSystem{
	effects: make([]SoundEffect, 64),
}

type AudioSystem struct {
	mutex   sync.Mutex
	effects []SoundEffect

	channels [4]ChannelState // TODO CHANNELS MUST BE STARTED FROM 1 TO HONOR CHANNEL MASK

	music MusicState
}

type ChannelState struct {
	No      byte
	Playing bool
	Sfx     SfxSlice
	Note    byte
	// TODO some additional information about progress of note? ticks? time?
}

type SfxSlice struct {
	No        byte
	NoteStart byte
	NoteEnd   byte
}

type MusicState struct {
	PatterNo    byte
	FadeMs      int
	ChannelMask byte
	// CurrentSfx? CurrentNote?
}

func (s *AudioSystem) SoundEffect(no int) SoundEffect {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.effects[no]
}

func (s *AudioSystem) playSfx(effectNo byte, channel byte, noteStart, noteEnd byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.channels[channel] = ChannelState{
		No: 0,
		Sfx: SfxSlice{
			No:        effectNo,
			NoteStart: noteStart,
			NoteEnd:   noteEnd,
		},
		Note: noteStart,
	}
}

func (s *AudioSystem) stopChannel(channel byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.channels[channel].Playing = false
}

func (s *AudioSystem) SetSoundEffect(effect SoundEffect) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.effects[effect.No] = effect
}

type SoundEffect struct {
	No    byte
	Speed byte
	Notes [32]Note // it's an array because we want to copy it each time (and don't allocate the heap)
}

type Note struct {
	Pitch  uint16
	Volume uint // 0 - 7
	Wave   Wave
}

// Wave is a function generating wave. Phase parameter is [0,2PI). Return value is [-1,1].
type Wave func(phase float64) float64

const tau = 2 * math.Pi

var (
	WaveTriangle Wave = func(phase float64) float64 {
		val := 2.0*(phase*(1.0/tau)) - 1.0
		if val < 0.0 {
			val = -val
		}
		val = 2.0 * (val - 0.5)
		return val
	}

	//
	//WaveTiltedSaw Wave = 2
	//WaveSaw       Wave = 2
	WaveSquare Wave = func(phase float64) float64 {
		val := -1.0
		if phase <= math.Pi {
			val = 1.0
		}
		return val
	}
)

// Sfx plays sound effect with given number on any available channel.
func Sfx(no byte) {
	channel, ok := audio.AvailableChannel()
	if !ok {
		return
	}

	audio.playSfx(no, channel, 0, 31)
}

// SfxChan plays sound effect with given number on specified channel.
func SfxChan(no byte, channel byte) {
	audio.playSfx(no, channel, 0, 31)
}

func SfxChanRange(no byte, channel byte, noteStart, noteEnd byte) {
	audio.playSfx(no, channel, noteStart, noteEnd)
}

func SfxStop(channel byte) {
	audio.stopChannel(channel)
}

func Music(pattern byte) {
	audio.playMusic(pattern, 0, 0)
}

func MusicStop() {}

func MusicFade(pattern byte, fadeMs int) {}

func MusicFadeChanMask(pattern byte, fadeMs int, channelMask byte) {}

func (s *AudioSystem) Read(p []byte) (n int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO should clear p first?

	for _, _ = range s.channels {
		// generate wave of p len for channel
		// add each value to p
	}

	// generate music (only for channels not currently playing)

	return 0, nil
}

func (s *AudioSystem) Channel(i int) ChannelState {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.channels[i]
}

func (s *AudioSystem) SetChannel(state ChannelState) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.channels[state.No] = state
}

func (s *AudioSystem) playMusic(pattern byte, fadeMs int, channelMask byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.music = MusicState{
		PatterNo:    pattern,
		FadeMs:      fadeMs,
		ChannelMask: channelMask,
	}
}

func (s *AudioSystem) Music() MusicState {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.music
}

func (s *AudioSystem) AvailableChannel() (byte, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, channel := range s.channels {
		if !channel.Playing && (s.music.ChannelMask&channel.No) != channel.No {
			return channel.No, true
		}
	}

	return 0, false
}
