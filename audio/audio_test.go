// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/audio"
)

const (
	maxSfxNo     = 63
	maxPatternNo = 63
)

var (
	validEffect = audio.SoundEffect{
		Notes: [32]audio.Note{
			{
				Pitch:      audio.PitchC2,
				Instrument: audio.InstrumentOrgan,
				Volume:     audio.VolumeLoudest,
				Effect:     audio.EffectSlide,
			},
		},
		Speed:     3,
		LoopStart: 4,
		LoopStop:  5,
		Noiz:      false,
		Buzz:      true,
		Detune:    0,
		Reverb:    1,
		Dampen:    2,
	}

	validPattern = audio.Pattern{
		Sfx: [4]audio.PatternSfx{
			{SfxNo: 1, Enabled: true},
			{SfxNo: 2, Enabled: false},
		},
		BeginLoop:    false,
		EndLoop:      true,
		StopAtTheEnd: false,
	}
)

func TestSynthesizer_SetSfx(t *testing.T) {
	t.Run("should update sfx", func(t *testing.T) {
		for i := 0; i < maxSfxNo; i++ {
			s := audio.Synthesizer{}
			// when
			s.SetSfx(i, validEffect)
			// then
			actual := s.GetSfx(i)
			assert.Equal(t, validEffect, actual)
		}
	})

	t.Run("should clamp parameters", func(t *testing.T) {
		s := audio.Synthesizer{}
		e := validEffect
		e.LoopStart = 64
		e.LoopStop = 64
		e.Detune = 3
		e.Reverb = 3
		e.Dampen = 3
		e.Notes[2].Volume = audio.VolumeLoudest + 1
		e.Notes[1].Pitch = 64
		e.Notes[3].Instrument = 16
		e.Notes[4].Effect = 8
		// when
		s.SetSfx(1, e)
		// then
		actual := s.GetSfx(1)
		assert.Equal(t, byte(63), actual.LoopStart)
		assert.Equal(t, byte(63), actual.LoopStop)
		assert.Equal(t, byte(2), actual.Detune)
		assert.Equal(t, byte(2), actual.Reverb)
		assert.Equal(t, byte(2), actual.Dampen)
		assert.Equal(t, audio.VolumeLoudest, actual.Notes[2].Volume)
		assert.Equal(t, audio.Pitch(63), actual.Notes[1].Pitch)
		assert.Equal(t, audio.Instrument(15), actual.Notes[3].Instrument)
		assert.Equal(t, audio.Effect(7), actual.Notes[4].Effect)
	})

	t.Run("should not set sfx number which is out of range", func(t *testing.T) {
		t.Run("maxSfxNo+1", func(t *testing.T) {
			s := audio.Synthesizer{}
			s.SetSfx(maxSfxNo+1, validEffect)
			sfx := s.GetSfx(maxSfxNo + 1)
			assert.Equal(t, audio.SoundEffect{}, sfx) // not found
		})

		t.Run("negative", func(t *testing.T) {
			s := audio.Synthesizer{}
			s.SetSfx(-1, validEffect)
			sfx := s.GetSfx(-1)
			assert.Equal(t, audio.SoundEffect{}, sfx) // not found
		})
	})
}

func TestSynthesizer_GetSfx(t *testing.T) {
	t.Run("should return default sfx", func(t *testing.T) {
		s := audio.Synthesizer{}
		actual := s.GetSfx(0)
		assert.Equal(t, audio.SoundEffect{}, actual)
	})

	t.Run("should return zero-value sfx for number out range", func(t *testing.T) {
		s := audio.Synthesizer{}
		for i := 0; i <= maxSfxNo; i++ {
			s.SetSfx(i, validEffect)
		}

		tests := []int{-1, -255, 256}
		for _, sfxNo := range tests {
			// when
			actual := s.GetSfx(sfxNo)
			// then
			assert.Equalf(t, audio.SoundEffect{}, actual, "sfxNo=%d", sfxNo)
		}
	})
}

func TestSynthesizer_SetMusic(t *testing.T) {
	t.Run("should update pattern", func(t *testing.T) {
		for i := 0; i < maxPatternNo; i++ {
			s := audio.Synthesizer{}
			// when
			s.SetMusic(i, validPattern)
			// then
			actual := s.GetMusic(i)
			assert.Equal(t, validPattern, actual)
		}
	})

	t.Run("should clamp parameters", func(t *testing.T) {
		s := audio.Synthesizer{}
		e := validPattern
		e.Sfx[1].SfxNo = 64
		// when
		s.SetMusic(1, e)
		// then
		actual := s.GetMusic(1)
		assert.Equal(t, byte(maxSfxNo), actual.Sfx[1].SfxNo)
	})

	t.Run("should not set pattern number which is out of range", func(t *testing.T) {
		t.Run("maxPatternNo+1", func(t *testing.T) {
			s := audio.Synthesizer{}
			s.SetMusic(maxPatternNo+1, validPattern)
			pattern := s.GetMusic(maxPatternNo + 1)
			assert.Equal(t, audio.Pattern{}, pattern) // not found
		})

		t.Run("negative", func(t *testing.T) {
			s := audio.Synthesizer{}
			s.SetMusic(-1, validPattern)
			pattern := s.GetMusic(-1)
			assert.Equal(t, audio.Pattern{}, pattern) // not found
		})
	})
}

func TestSynthesizer_GetMusic(t *testing.T) {
	t.Run("should return default pattern", func(t *testing.T) {
		s := audio.Synthesizer{}
		actual := s.GetMusic(0)
		assert.Equal(t, audio.Pattern{}, actual)
	})

	t.Run("should return zero-value pattern for number out range", func(t *testing.T) {
		s := audio.Synthesizer{}
		for i := 0; i <= maxPatternNo; i++ {
			s.SetMusic(i, validPattern)
		}

		tests := []int{-1, -255, 256}
		for _, patterNo := range tests {
			// when
			actual := s.GetMusic(patterNo)
			// then
			assert.Equalf(t, audio.Pattern{}, actual, "patternNo=%d", patterNo)
		}
	})
}

func TestSynthesizer_Read(t *testing.T) {
	t.Run("should read 0 bytes when buffer is empty", func(t *testing.T) {
		s := audio.Synthesizer{}
		n, err := s.ReadSamples(nil)
		assert.Zero(t, n)
		assert.NoError(t, err)
	})

	t.Run("should clear the buffer with 0 when no channels are used", func(t *testing.T) {
		buffer := []float32{1, 2, 3, 4}
		s := audio.Synthesizer{}
		n, err := s.ReadSamples(buffer)
		require.NotZero(t, n)
		require.NoError(t, err)
		expected := make([]float32, n)
		assert.Equal(t, expected, buffer)
	})
}

//go:embed "valid-save"
var validSave []byte

func TestSynthesizer_Save(t *testing.T) {
	t.Run("should save in binary format", func(t *testing.T) {
		s := audio.Synthesizer{}
		s.SetSfx(3, validEffect)
		s.SetMusic(4, validPattern)
		bytes, err := s.Save()
		require.NoError(t, err)
		assert.Equal(t, validSave, bytes)
	})
}

func TestSynthesizer_Load(t *testing.T) {
	t.Run("should load state", func(t *testing.T) {
		s := audio.Synthesizer{}
		err := s.Load(validSave)
		require.NoError(t, err)
		assert.Equal(t, validEffect, s.GetSfx(3))
		assert.Equal(t, validPattern, s.GetMusic(4))
	})

	t.Run("should return error when state is empty", func(t *testing.T) {
		s := audio.Synthesizer{}
		err := s.Load([]byte{})
		assert.Error(t, err)
	})

	t.Run("should return error when version is not supported", func(t *testing.T) {
		s := audio.Synthesizer{}
		err := s.Load([]byte{2})
		assert.Error(t, err)
	})

	t.Run("should return error when state has invalid length", func(t *testing.T) {
		s := audio.Synthesizer{}
		err := s.Load([]byte{1, 0, 0, 0})
		assert.Error(t, err)
	})

	t.Run("should clamp sfx when sfx parameter is out of range", func(t *testing.T) {
		save := clone(validSave)
		const sfx0note0pitch = 1
		save[sfx0note0pitch] = 64
		s := audio.Synthesizer{}
		// when
		err := s.Load(save)
		// then
		require.NoError(t, err)
		assert.Equal(t, audio.Pitch(63), s.GetSfx(0).Notes[0].Pitch)
	})

	t.Run("should clamp pattern when pattern parameter is out of range", func(t *testing.T) {
		save := clone(validSave)
		const pattern0sfx0sfxNo = 8705 // 1 byte for version, 8704 bytes for sound effects
		save[pattern0sfx0sfxNo] = 64
		s := audio.Synthesizer{}
		// when
		err := s.Load(save)
		// then
		require.NoError(t, err)
		assert.Equal(t, byte(63), s.GetMusic(0).Sfx[0].SfxNo)
	})
}

func clone(s []byte) []byte {
	cloned := make([]byte, len(s))
	copy(cloned, s)
	return cloned
}
