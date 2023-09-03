// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"bytes"
	"fmt"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio/internal"
)

const SampleRate = internal.SampleRate

// Synthesizer is used by a back-end. It is an System implementation.
// Plus it provides method for synthesizing audio samples.
type Synthesizer struct {
	sfx      map[byte]SoundEffect
	pattern  map[byte]Pattern
	channels [internal.Channels]channel
}

// ReadSamples method is used by a back-end to read generated audio samples and play them back to the user. The sample rate is 22050.
// One channel is used (mono).
//
// ReadSamples is (usually) executed concurrently with main game loop. Back-end should add proper synchronization to avoid race conditions.
// Back-end could decide about buffer size, although the higher the size the higher the lag. Usually the buffer is 512 samples,
// which is 23ms of audio.
//
// Values written to the buffer are usually in range between -1.0 and 1.0, but sometimes they can exceed the range
// (for example due to audio channels summing). Min is -4.0, max is 4.0 inclusive.
func (s *Synthesizer) ReadSamples(buffer []float64) {
	if len(buffer) == 0 {
		return
	}

	for i := 0; i < len(buffer); i++ {
		buffer[i] = s.readSample()
	}
}

func (s *Synthesizer) readSample() float64 {
	var sampleChannels internal.SampleChannels

	for i, ch := range s.channels {
		sfx := s.GetSfx(ch.sfxNo)
		sampleChannels[i] = ch.readSample(sfx)
		s.channels[i] = ch
	}

	return sampleChannels.Sum()
}

func (c *channel) readSample(sfx SoundEffect) float64 {
	if !c.playing {
		return 0
	}

	var sample float64

	volume := float64(sfx.noteAt(c.noteNo).Volume) / 7
	sample = c.oscillator.NextSample() * volume

	c.frame += 1
	noteHasEnded := c.frame == c.noteEndFrame
	if noteHasEnded {
		c.moveToNextNote(sfx)
	}

	return sample
}

func (c *channel) moveToNextNote(sfx SoundEffect) {
	c.notesToGo--
	c.noteNo++

	if c.noteNo == int(sfx.LoopStop) && !c.loopingDisabled {
		c.noteNo = int(sfx.LoopStart)
	}

	maxLenReached := sfx.LoopStop == 0 && int(sfx.LoopStart) == c.noteNo
	if c.notesToGo == 0 || maxLenReached {
		c.playing = false
		return
	}

	c.noteEndFrame += singleNoteSamples(sfx.Speed)

	note := sfx.noteAt(c.noteNo)

	c.oscillator.Func = oscillatorFunc(note.Instrument)
	c.oscillator.FreqHz = pitchToFreq(note.Pitch)
}

func (s *Synthesizer) Sfx(sfxNo int, ch Channel, offset, length int) {
	if sfxNo == -2 {
		s.disableLooping(ch)
		return
	}

	if ch >= ChannelStop && ch <= Channel3 {
		s.stopSfx(sfxNo)
	}

	if ch == ChannelAny {
		ch = s.findAvailableChannel()
	}

	if ch < Channel0 || ch > Channel3 {
		return
	}

	if sfxNo == -1 {
		s.channels[ch].playing = false
		return
	}

	offset = pi.MidInt(offset, 0, 31)

	s.channels[ch].playing = true

	sfx := s.GetSfx(sfxNo)

	s.channels[ch].sfxNo = sfxNo
	s.channels[ch].frame = 0
	s.channels[ch].noteNo = offset
	s.channels[ch].notesToGo = length

	s.channels[ch].noteEndFrame = singleNoteSamples(sfx.Speed)

	note0 := sfx.Notes[offset]
	s.channels[ch].oscillator.Func = oscillatorFunc(note0.Instrument)
	s.channels[ch].oscillator.FreqHz = pitchToFreq(note0.Pitch)
}

func (s *Synthesizer) disableLooping(ch Channel) {
	if ch == ChannelAny || ch == ChannelStop {
		for i := range s.channels {
			s.channels[i].loopingDisabled = true
		}
		return
	}

	if ch >= Channel0 && ch <= Channel3 {
		s.channels[ch].loopingDisabled = true
	}
}

func (s *Synthesizer) stopSfx(no int) {
	for i, c := range s.channels {
		if c.playing && c.sfxNo == no {
			c.playing = false
			s.channels[i] = c
			return
		}
	}
}

func (s *Synthesizer) findAvailableChannel() Channel {
	for i, c := range s.channels {
		if !c.playing {
			return Channel(i)
		}
	}

	return Channel3
}

func (s *Synthesizer) Music(patterNo int, fadeMs int, channelMask byte) {
	fmt.Println("Music is not implemented yet. Sorry...")
}

func (s *Synthesizer) Stat() Stat {
	fmt.Println("Stat is not implemented yet. Sorry...")

	stat := Stat{}
	for i, c := range s.channels {
		if c.playing {
			stat.Sfx[i] = c.sfxNo
		} else {
			stat.Sfx[i] = -1
		}
	}
	return stat
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

type channel struct {
	sfxNo           int
	noteNo          int
	notesToGo       int
	frame           int
	noteEndFrame    int
	oscillator      internal.Oscillator
	playing         bool
	loopingDisabled bool
}
