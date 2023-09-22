// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"fmt"

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

func (s *Synthesizer) Play(sfxNo, ch, offset, length int) {
	if ch < -1 || ch > 3 {
		return
	}

	if sfxNo < 0 || sfxNo > maxSfxNo {
		return
	}

	s.Stop(sfxNo)

	if ch == -1 {
		ch = s.findAvailableChannel()
	}

	offset = midInt(offset, 0, 31)

	s.channels[ch].playing = true

	sfx := s.GetSfx(sfxNo)

	s.channels[ch].sfxNo = sfxNo
	s.channels[ch].frame = 0
	s.channels[ch].noteNo = offset
	if length <= 0 && sfx.LoopStop <= sfx.LoopStart {
		length = 32
	}
	if length > 32 {
		length = 32
	}
	s.channels[ch].notesToGo = length
	s.channels[ch].loopingDisabled = false

	s.channels[ch].noteEndFrame = singleNoteSamples(sfx.Speed)

	note0 := sfx.Notes[offset]
	s.channels[ch].oscillator.Func = oscillatorFunc(note0.Instrument)
	s.channels[ch].oscillator.FreqHz = pitchToFreq(note0.Pitch)
}

func (s *Synthesizer) Stop(sfxNo int) {
	if sfxNo == -1 {
		for i, c := range s.channels {
			c.playing = false
			s.channels[i] = c
		}
		return
	}

	for i, c := range s.channels {
		if c.playing && c.sfxNo == sfxNo {
			c.playing = false
			s.channels[i] = c
			return
		}
	}
}

func (s *Synthesizer) StopChan(channel int) {
	if channel == -1 {
		s.Stop(-1)
		return
	}
	if channel < 0 || channel > 3 {
		return
	}

	s.channels[channel].playing = false
}

func (s *Synthesizer) StopLoop(channel int) {
	if channel == -1 {
		for i := range s.channels {
			s.channels[i].loopingDisabled = true
		}
		return
	}

	if channel >= 0 && channel <= 3 {
		s.channels[channel].loopingDisabled = true
	}
}

func (s *Synthesizer) findAvailableChannel() int {
	for i, c := range s.channels {
		if !c.playing {
			return i
		}
	}

	return 3
}

func (s *Synthesizer) Music(patterNo int, fadeMs int, channelMask byte) {
	fmt.Println("Music is not implemented yet. Sorry...")
}

func (s *Synthesizer) Stat() Stat {
	stat := Stat{}
	for i, c := range s.channels {
		if c.playing {
			sfxStat := SfxStat{
				SfxNo:     c.sfxNo,
				Note:      c.noteNo,
				Remaining: c.notesToGo,
			}
			sfx := s.GetSfx(c.sfxNo)
			if sfx.hasLoop() {
				sfxStat.Remaining = -1
			} else if sfx.hasLength() {
				sfxStat.Remaining = minInt(sfx.length(), c.notesToGo)
			}
			stat.Sfx[i] = sfxStat
		} else {
			stat.Sfx[i] = SfxStat{
				SfxNo:     -1,
				Note:      -1,
				Remaining: 0,
			}
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

	effect.LoopStart = minInt(63, effect.LoopStart)
	effect.LoopStop = minInt(63, effect.LoopStop)
	effect.Detune = minInt(2, effect.Detune)
	effect.Reverb = minInt(2, effect.Reverb)
	effect.Dampen = minInt(2, effect.Dampen)

	for i := 0; i < len(effect.Notes); i++ {
		volume := effect.Notes[i].Volume
		effect.Notes[i].Volume = minInt(7, volume)

		pitch := effect.Notes[i].Pitch
		effect.Notes[i].Pitch = minInt(63, pitch)

		instrument := effect.Notes[i].Instrument
		effect.Notes[i].Instrument = minInt(15, instrument)

		eff := effect.Notes[i].Effect
		effect.Notes[i].Effect = minInt(7, eff)
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
		pattern.Sfx[i].SfxNo = minInt(maxSfxNo, sfxNo)
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
