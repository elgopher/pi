// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package audio provides functions to play sound effects and music.
// Call Sfx() to play sound effect and Music() to play music.
//
// Sound effects and music can be changed using functions: SetSfx and SetMusic.
package audio

import (
	"bytes"
	"fmt"

	"github.com/elgopher/pi"
)

// Sfx starts playing sound effect with given sfxNo on specified channel.
//
// sfxNo is the number of sound effect to play (0-63). sfxNo=-1 stops
// playing sound on the given channel. sfxNo=-2 releases the sound
// of the given channel from looping.
//
// channel is the channel number from 0-3. Channel=-1 chooses available
// channel automatically. Channel=-2 stops playing the given sound effect
// on any channel. See pi.Channel.
//
// offset is the note position to start playing (0-31).
//
// length is the number of notes to play (0-31).
func Sfx(sfxNo int, channel Channel, offset, length int) {
	audioSystem.Sfx(sfxNo, channel, offset, length)
}

// Music starts playing music with given patterNo.
//
// patterNo is the number of music pattern to play (0-63). patternNo=-1 stops
// playing music.
//
// fadeMs fades in (or out) the music volume over a duration, given as number
// of milliseconds.
//
// channelMask is a bitfield indicating which of the four sound channels
// should be reserved for music.
func Music(patterNo int, fadeMs int, channelMask byte) {
	audioSystem.Music(patterNo, fadeMs, channelMask)
}

// GetStat returns information about currently played sound effects and music.
func GetStat() Stat {
	return audioSystem.Stat()
}

// SetSfx updates the sound effect. sfxNo is 0-63. Updating sfx number which
// is higher than 63 does not do anything.
//
// SoundEffect parameters are clamped when out of range.
// For example, sfx note volume equal 8 will be silently clamped to 7.
func SetSfx(sfxNo int, e SoundEffect) {
	audioSystem.SetSfx(sfxNo, e)
}

// GetSfx returns sfx with given number. sfxNo is 0-63. Trying to get
// sfx number higher than 63 will result in returning empty SoundEffect (zero-value).
func GetSfx(sfxNo int) SoundEffect {
	return audioSystem.GetSfx(sfxNo)
}

// SetMusic updates the music pattern. patternNo is 0-63. Updating pattern number which
// is higher than 63 does not do anything.
//
// Pattern parameters are clamped when out of range.
// For example, pattern sfx number equal 64 will be silently clamped to 63.
func SetMusic(patternNo int, p Pattern) {
	audioSystem.SetMusic(patternNo, p)
}

// GetMusic returns music pattern with given number. patterNo is 0-63. Trying to get
// pattern number higher than 63 will result in returning empty Pattern (zero-value).
func GetMusic(patterNo int) Pattern {
	return audioSystem.GetMusic(patterNo)
}

// SaveAudio stores audio system state to byte slice. State is stored in binary form.
// The format is described in Synthesizer.Save source code.
func SaveAudio() ([]byte, error) {
	return audioSystem.Save()
}

// LoadAudio restores audio system state from byte slice. State is restored from binary form.
// The format is described in Synthesizer.Save source code.
func LoadAudio(b []byte) error {
	return audioSystem.Load(b)
}

var audioSystem AudioSystem = &Synthesizer{}

type Channel int8

const (
	Channel0    Channel = 0
	Channel1    Channel = 1
	Channel2    Channel = 2
	Channel3    Channel = 3
	ChannelAny  Channel = -1
	ChannelStop Channel = -2
)

type Stat struct {
	Sfx           [4]int // -1 means no sfx on channel
	Note          [4]int // -1 means no sfx on channel
	Pattern       int    // currently played music pattern
	PatternsCount int    // the number of music patterns played since the most recent call to Music().
	TicksCount    int    // the number of ticks (notes or rests) played on the current pattern.
}

type SoundEffect struct {
	Notes     [32]Note
	Speed     byte // 1 is the fastest (~8.33ms). For 120 the length of one note becomes 1 second. 0 means SoundEffect takes no time.
	LoopStart byte // 0-63. Notes 32-63 are silent.
	LoopStop  byte // 0-63. Notes 32-63 are silent.
	Detune    byte // 0 (disabled), 1 or 2
	Reverb    byte // 0 (disabled), 1 or 2
	Dampen    byte // 0 (disabled), 1 or 2
	Noiz      bool
	Buzz      bool
}

type Note struct {
	Pitch      Pitch      // 0-63
	Instrument Instrument // 0-15
	Volume     Volume     // 0-7
	Effect     Effect     // 0-7
}

type Volume byte

const (
	VolumeQuiet   Volume = 0
	VolumeLoudest Volume = 7
)

type Pitch byte

const (
	PitchC0  Pitch = 0
	PitchCs0 Pitch = 1 // C#0
	PitchD0  Pitch = 2
	PitchDs0 Pitch = 3
	PitchE0  Pitch = 4
	PitchF0  Pitch = 5
	PitchFs0 Pitch = 6
	PitchG0  Pitch = 7
	PitchGs0 Pitch = 8
	PitchA0  Pitch = 9
	PitchAs0 Pitch = 10
	PitchB0  Pitch = 11

	PitchC1  Pitch = 12
	PitchCs1 Pitch = 13
	PitchD1  Pitch = 14
	PitchDs1 Pitch = 15
	PitchE1  Pitch = 16
	PitchF1  Pitch = 17
	PitchFs1 Pitch = 18
	PitchG1  Pitch = 19
	PitchGs1 Pitch = 20
	PitchA1  Pitch = 21
	PitchAs1 Pitch = 22
	PitchB1  Pitch = 23

	PitchC2  Pitch = 24
	PitchCs2 Pitch = 25
	PitchD2  Pitch = 26
	PitchDs2 Pitch = 27
	PitchE2  Pitch = 28
	PitchF2  Pitch = 29
	PitchFs2 Pitch = 30
	PitchG2  Pitch = 31
	PitchGs2 Pitch = 32
	PitchA2  Pitch = 33
	PitchAs2 Pitch = 34
	PitchB2  Pitch = 35

	PitchC3  Pitch = 36
	PitchCs3 Pitch = 37
	PitchD3  Pitch = 38
	PitchDs3 Pitch = 39
	PitchE3  Pitch = 40
	PitchF3  Pitch = 41
	PitchFs3 Pitch = 42
	PitchG3  Pitch = 43
	PitchGs3 Pitch = 44
	PitchA3  Pitch = 45
	PitchAs3 Pitch = 46
	PitchB3  Pitch = 47

	PitchC4  Pitch = 48
	PitchCs4 Pitch = 49
	PitchD4  Pitch = 50
	PitchDs4 Pitch = 51
	PitchE4  Pitch = 52
	PitchF4  Pitch = 53
	PitchFs4 Pitch = 54
	PitchG4  Pitch = 55
	PitchGs4 Pitch = 56
	PitchA4  Pitch = 57
	PitchAs4 Pitch = 58
	PitchB4  Pitch = 59

	PitchC5  Pitch = 60
	PitchCs5 Pitch = 61
	PitchD5  Pitch = 62
	PitchDs5 Pitch = 63
)

type Instrument byte

const (
	InstrumentTriangle  Instrument = 0
	InstrumentTiltedSaw Instrument = 1
	InstrumentSaw       Instrument = 2
	InstrumentSquare    Instrument = 3
	InstrumentPulse     Instrument = 4
	InstrumentOrgan     Instrument = 5
	InstrumentNoise     Instrument = 6
	InstrumentPhaser    Instrument = 7
	InstrumentSfx0      Instrument = 8
	InstrumentSfx1      Instrument = 9
	InstrumentSfx2      Instrument = 10
	InstrumentSfx3      Instrument = 11
	InstrumentSfx4      Instrument = 12
	InstrumentSfx5      Instrument = 13
	InstrumentSfx6      Instrument = 14
	InstrumentSfx7      Instrument = 15
)

type Effect byte

const (
	EffectNoEffect Effect = 0
	EffectSlide    Effect = 1
	EffectVibrato  Effect = 2
	EffectDrop     Effect = 3
	EffectFadeIn   Effect = 4
	EffectFadeOut  Effect = 5
	EffectArpFast  Effect = 6
	EffectArpSlow  Effect = 7
)

// Pattern is a music pattern
type Pattern struct {
	Sfx          [4]PatternSfx
	BeginLoop    bool
	EndLoop      bool
	StopAtTheEnd bool
}

type PatternSfx struct {
	SfxNo   byte // 0-63
	Enabled bool
}

// Synthesizer is used by a back-end. It is an AudioSystem implementation.
// Plus it provides method for synthesizing audio samples.
type Synthesizer struct {
	sfx     map[byte]SoundEffect
	pattern map[byte]Pattern
}

// ReadSamples method is used by a back-end to read generated audio samples and play them back to the user. The sample rate is 22050.
// One channel is used (mono).
//
// ReadSamples is (usually) executed concurrently with main game loop. Back-end should add proper synchronization to avoid race conditions.
// Back-end could decide about buffer size, although the higher the size the higher the lag. Usually the buffer is 512 samples,
// which is 23ms of audio.
func (s *Synthesizer) ReadSamples(b []float32) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}

	// generate silence for now
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}

	return len(b), nil
}

func (s *Synthesizer) Sfx(sfxNo int, channel Channel, offset, length int) {
	fmt.Println("Sfx is not implemented yet. Sorry...")
}

func (s *Synthesizer) Music(patterNo int, fadeMs int, channelMask byte) {
	fmt.Println("Music is not implemented yet. Sorry...")
}

func (s *Synthesizer) Stat() Stat {
	fmt.Println("Stat is not implemented yet. Sorry...")
	return Stat{}
}

const maxSfxNo = 63

func (s *Synthesizer) SetSfx(sfxNo int, effect SoundEffect) {
	if sfxNo < 0 || sfxNo > maxSfxNo {
		return
	}

	if s.sfx == nil {
		s.sfx = map[byte]SoundEffect{}
	}

	effect.LoopStart = pi.MinInt(63, effect.LoopStart)
	effect.LoopStop = pi.MinInt(63, effect.LoopStop)
	effect.Detune = pi.MinInt(2, effect.Detune)
	effect.Reverb = pi.MinInt(2, effect.Reverb)
	effect.Dampen = pi.MinInt(2, effect.Dampen)

	for i := 0; i < len(effect.Notes); i++ {
		volume := effect.Notes[i].Volume
		effect.Notes[i].Volume = pi.MinInt(7, volume)

		pitch := effect.Notes[i].Pitch
		effect.Notes[i].Pitch = pi.MinInt(63, pitch)

		instrument := effect.Notes[i].Instrument
		effect.Notes[i].Instrument = pi.MinInt(15, instrument)

		eff := effect.Notes[i].Effect
		effect.Notes[i].Effect = pi.MinInt(7, eff)
	}

	s.sfx[byte(sfxNo)] = effect
}

func (s *Synthesizer) GetSfx(sfxNo int) (e SoundEffect) {
	if sfxNo < 0 || sfxNo > maxSfxNo {
		return e
	}

	return s.sfx[byte(sfxNo)]
}

const maxPatternNo = 63

func (s *Synthesizer) SetMusic(patternNo int, pattern Pattern) {
	if patternNo < 0 || patternNo > maxPatternNo {
		return
	}

	if s.pattern == nil {
		s.pattern = map[byte]Pattern{}
	}

	for i := 0; i < len(pattern.Sfx); i++ {
		sfxNo := pattern.Sfx[i].SfxNo
		pattern.Sfx[i].SfxNo = pi.MinInt(maxSfxNo, sfxNo)
	}

	s.pattern[byte(patternNo)] = pattern
}

func (s *Synthesizer) GetMusic(patterNo int) (p Pattern) {
	if patterNo < 0 || patterNo > maxPatternNo {
		return p
	}

	return s.pattern[byte(patterNo)]
}

const schemaVersion = 1

func (s *Synthesizer) Load(state []byte) error {
	if len(state) == 0 {
		return fmt.Errorf("state is empty")
	}

	version := state[0]
	if schemaVersion != version {
		return fmt.Errorf("state version %d is not supported. Only %d is supported.", version, schemaVersion)
	}

	const expectedStateLen = 9409
	if len(state) != expectedStateLen {
		return fmt.Errorf("invalid length of state. Must be %d.", expectedStateLen)
	}

	offset := 1

	for sfxNo := 0; sfxNo <= maxSfxNo; sfxNo++ {
		var sfx SoundEffect

		for j, note := range sfx.Notes {
			note.Pitch = Pitch(state[offset])
			offset++
			note.Instrument = Instrument(state[offset])
			offset++
			note.Volume = Volume(state[offset])
			offset++
			note.Effect = Effect(state[offset])
			offset++

			sfx.Notes[j] = note
		}

		sfx.Speed = state[offset]
		offset++
		sfx.LoopStart = state[offset]
		offset++
		sfx.LoopStop = state[offset]
		offset++
		sfx.Detune = state[offset]
		offset++
		sfx.Reverb = state[offset]
		offset++
		sfx.Dampen = state[offset]
		offset++
		sfx.Noiz = byteToBool(state[offset])
		offset++
		sfx.Buzz = byteToBool(state[offset])
		offset++

		s.SetSfx(sfxNo, sfx)
	}

	for patterNo := 0; patterNo <= maxPatternNo; patterNo++ {
		var pattern Pattern

		for j, sfx := range pattern.Sfx {
			sfx.SfxNo = state[offset]
			offset++
			sfx.Enabled = byteToBool(state[offset])
			offset++

			pattern.Sfx[j] = sfx
		}

		pattern.BeginLoop = byteToBool(state[offset])
		offset++
		pattern.EndLoop = byteToBool(state[offset])
		offset++
		pattern.StopAtTheEnd = byteToBool(state[offset])
		offset++

		s.SetMusic(patterNo, pattern)
	}

	return nil
}

func (s *Synthesizer) Save() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	buffer.WriteByte(schemaVersion)

	for i := 0; i <= maxSfxNo; i++ {
		sfx := s.GetSfx(i)
		for _, note := range sfx.Notes {
			buffer.WriteByte(byte(note.Pitch))
			buffer.WriteByte(byte(note.Instrument))
			buffer.WriteByte(byte(note.Volume))
			buffer.WriteByte(byte(note.Effect))
		}
		buffer.WriteByte(sfx.Speed)
		buffer.WriteByte(sfx.LoopStart)
		buffer.WriteByte(sfx.LoopStop)
		buffer.WriteByte(sfx.Detune)
		buffer.WriteByte(sfx.Reverb)
		buffer.WriteByte(sfx.Dampen)
		buffer.WriteByte(boolToByte(sfx.Noiz))
		buffer.WriteByte(boolToByte(sfx.Buzz))
	}

	for i := 0; i <= maxPatternNo; i++ {
		pattern := s.GetMusic(i)
		for _, sfx := range pattern.Sfx {
			buffer.WriteByte(sfx.SfxNo)
			buffer.WriteByte(boolToByte(sfx.Enabled))
		}
		buffer.WriteByte(boolToByte(pattern.BeginLoop))
		buffer.WriteByte(boolToByte(pattern.EndLoop))
		buffer.WriteByte(boolToByte(pattern.StopAtTheEnd))
	}

	return buffer.Bytes(), nil
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}

	return 0
}

func byteToBool(b byte) bool {
	return b == 1
}

// SetAudioSystem is executed by the back-end to replace audio system with his own implementation.
func SetAudioSystem(s AudioSystem) {
	audioSystem = s
}

type AudioSystem interface {
	Sfx(sfxNo int, channel Channel, offset, length int)
	Music(patterNo int, fadeMs int, channelMask byte)
	Stat() Stat
	SetSfx(sfxNo int, e SoundEffect)
	GetSfx(sfxNo int) SoundEffect
	SetMusic(patternNo int, _ Pattern)
	GetMusic(patterNo int) Pattern
	Save() ([]byte, error)
	Load([]byte) error
}
