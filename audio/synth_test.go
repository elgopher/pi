// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio_test

import (
	_ "embed"
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/audio"
)

const (
	maxSfxNo                     = 63
	maxPatternNo                 = 63
	maxChannels                  = 4
	durationOfNoteWhenSpeedIsOne = 183
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

//go:embed "internal/valid-save"
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

func TestSynthesizer_ReadSamples(t *testing.T) {
	t.Run("should not do anything when buffer has 0 length", func(t *testing.T) {
		s := audio.Synthesizer{}
		buffer := make([]float64, 0)
		s.ReadSamples(buffer)
	})

	t.Run("when no sound is playing ReadSamples should return silence", func(t *testing.T) {
		s := audio.Synthesizer{}
		buffer := []float64{1, 2, 3, 4}
		s.ReadSamples(buffer)
		assertSilence(t, buffer)
	})
}

func TestSynthesizer_Sfx(t *testing.T) {
	t.Run("should not panic when channel is not in range 0-3", func(t *testing.T) {
		channels := []audio.Channel{
			audio.ChannelStop - 1, audio.ChannelStop, audio.ChannelAny, audio.Channel3 + 1,
		}

		for _, ch := range channels {
			assert.NotPanics(t, func() {
				s := audio.Synthesizer{}
				// when
				s.Sfx(0, ch, 0, 0)
				s.ReadSamples(make([]float64, 256))
			})
		}
	})

	t.Run("should play sound on a given channel", func(t *testing.T) {
		for channelNo := 0; channelNo < maxChannels; channelNo++ {
			testName := fmt.Sprintf("channel %d", channelNo)

			t.Run(testName, func(t *testing.T) {
				s := audio.Synthesizer{}
				s.SetSfx(0, validEffect)

				buffer := make([]float64, 32)
				// when
				s.Sfx(0, audio.Channel(channelNo), 0, 1)
				// then
				s.ReadSamples(buffer)
				assertAllValuesBetween(t, -1.0, 1.0, buffer)
				assertAllValuesDifferent(t, buffer)
			})
		}
	})

	t.Run("should sum samples from all channels", func(t *testing.T) {
		singleChannelBuffer := generateSamples(validEffect, 1)

		s := audio.Synthesizer{}
		s.SetSfx(0, validEffect)
		// when
		for ch := 0; ch < maxChannels; ch++ {
			s.Sfx(0, audio.Channel(ch), 0, 1)
		}
		// then
		allChannelBuffer := make([]float64, 1)
		s.ReadSamples(allChannelBuffer)

		expectedSample := singleChannelBuffer[0] * maxChannels
		assert.InDelta(t, expectedSample, allChannelBuffer[0], 0.0000001)
		assertAllValuesBetween(t, -4.0, 4.0, allChannelBuffer)
	})

	t.Run("should stop playing on a given channel", func(t *testing.T) {
		for channelNo := audio.Channel(0); channelNo < maxChannels; channelNo++ {
			testName := fmt.Sprintf("channel %d", channelNo)

			t.Run(testName, func(t *testing.T) {
				s := audio.Synthesizer{}
				s.SetSfx(0, validEffect)
				s.Sfx(0, channelNo, 0, 31)
				s.ReadSamples(make([]float64, 1))
				// when
				s.Sfx(-1, channelNo, 0, 0)
				// then
				buffer := make([]float64, 16)
				s.ReadSamples(buffer)
				assertSilence(t, buffer)
			})
		}
	})

	t.Run("should play sound with specified pitch", func(t *testing.T) {
		const bufferSize = 256

		c0 := validEffect
		c0.Notes[0].Pitch = audio.PitchC0
		c0buffer := generateSamples(c0, bufferSize)

		a3 := validEffect
		a3.Notes[0].Pitch = audio.PitchA3
		a3buffer := generateSamples(a3, bufferSize)
		// then
		assert.NotEqual(t, c0buffer, a3buffer, "buffers for pitch C0 and A3 should be different but are the same")
	})

	t.Run("should play sound with specified instrument", func(t *testing.T) {
		const bufferSize = 256

		triangle := validEffect
		triangle.Notes[0].Instrument = audio.InstrumentTriangle
		triangleBuffer := generateSamples(triangle, bufferSize)

		organ := validEffect
		organ.Notes[0].Instrument = audio.InstrumentOrgan
		organBuffer := generateSamples(organ, bufferSize)
		// then
		assert.NotEqual(t, triangleBuffer, organBuffer, "buffers for triangle and organ waves should be different but are the same")
	})

	t.Run("should not play first note when volume is 0", func(t *testing.T) {
		e := validEffect
		e.Speed = 1
		e.Notes[0].Volume = 0

		buffer := generateSamples(e, durationOfNoteWhenSpeedIsOne)
		assertSilence(t, buffer)
	})

	t.Run("should not play first note when volume is 0 and speed is 2", func(t *testing.T) {
		e := validEffect
		e.Speed = 2
		e.Notes[0].Volume = 0
		e.Notes[1].Volume = audio.VolumeLoudest

		buffer := generateSamples(e, 2*durationOfNoteWhenSpeedIsOne)
		assertSilence(t, buffer)
	})

	t.Run("should not play second note when volume is 0", func(t *testing.T) {
		for speed := 1; speed <= 2; speed++ {
			testName := fmt.Sprintf("speed=%d", speed)

			t.Run(testName, func(t *testing.T) {
				e := validEffect
				e.Speed = byte(speed)
				e.Notes[0].Volume = audio.VolumeLoudest
				e.Notes[1].Volume = audio.VolumeSilence
				e.Notes[2].Volume = audio.VolumeLoudest

				synth := audio.Synthesizer{}
				synth.SetSfx(0, e)
				synth.Sfx(0, audio.Channel0, 0, 31)
				buffer := make([]float64, speed*durationOfNoteWhenSpeedIsOne)
				// skip 1st note
				synth.ReadSamples(buffer)
				// when
				synth.ReadSamples(buffer)
				assertSilence(t, buffer)

				t.Run("and next note with max volume should be played", func(t *testing.T) {
					synth.ReadSamples(buffer)
					assertNotSilence(t, buffer)
				})
			})
		}
	})

	t.Run("should play sfx to the end", func(t *testing.T) {
		e := validEffect
		e.Speed = 1
		for i := 0; i < len(e.Notes); i++ {
			e.Notes[i].Volume = audio.VolumeLoudest
		}

		synth := audio.Synthesizer{}
		synth.SetSfx(0, e)
		synth.Sfx(0, audio.Channel0, 0, 31)

		buffer := make([]float64, len(e.Notes)*durationOfNoteWhenSpeedIsOne)
		// read the entire sfx
		synth.ReadSamples(buffer)
		// and then read silence
		synth.ReadSamples(buffer)
		assertSilence(t, buffer)
	})

	t.Run("should play sfx from the beginning", func(t *testing.T) {
		e := validEffect
		e.Speed = 1
		synth := audio.Synthesizer{}
		synth.SetSfx(0, e)
		synth.Sfx(0, audio.Channel0, 0, 31) // first call to Sfx
		buffer1 := make([]float64, len(e.Notes)*durationOfNoteWhenSpeedIsOne)
		synth.ReadSamples(buffer1) // read entire sound
		// when
		synth.Sfx(0, audio.Channel0, 0, 31)
		// then
		buffer2 := make([]float64, len(e.Notes)*durationOfNoteWhenSpeedIsOne)
		synth.ReadSamples(buffer2)
		assert.True(t, dominantFrequency(buffer1) == dominantFrequency(buffer2), "frequency should be the same")
	})

	t.Run("should change oscillator frequency when second note has different pitch", func(t *testing.T) {
		e := audio.SoundEffect{
			Notes: [32]audio.Note{
				{Pitch: audio.PitchC0, Volume: audio.VolumeLoudest},
				{Pitch: audio.PitchDs5, Volume: audio.VolumeLoudest},
			},
			Speed: 255,
		}
		synth := audio.Synthesizer{}
		synth.SetSfx(0, e)
		synth.Sfx(0, 0, 0, 31)
		buffer1 := make([]float64, 255*durationOfNoteWhenSpeedIsOne)
		synth.ReadSamples(buffer1)
		buffer2 := make([]float64, 255*durationOfNoteWhenSpeedIsOne)
		// when
		synth.ReadSamples(buffer2)
		// then
		assert.True(t, dominantFrequency(buffer1) < dominantFrequency(buffer2), "frequency of pitch C1 must be smaller than D#5")
	})
}

func clone(s []byte) []byte {
	cloned := make([]byte, len(s))
	copy(cloned, s)
	return cloned
}

func assertSilence(t *testing.T, buffer []float64) {
	zeroBuffer := make([]float64, len(buffer))
	assert.Equal(t, zeroBuffer, buffer, "buffer should have zeroes only (silence)")
}

func assertNotSilence(t *testing.T, buffer []float64) {
	zeroBuffer := make([]float64, len(buffer))
	assert.NotEqual(t, zeroBuffer, buffer, "buffer should not have zeroes only (no silence)")
}

func assertAllValuesDifferent(t *testing.T, buffer []float64) {
	current := buffer[0]
	for i := 1; i < len(buffer); i++ {
		require.NotEqualf(t, current, buffer[i], "adjacent buffer values at %d and %d are the same but should be different", i-1, i)
		current = buffer[i]
	}
}

func generateSamples(e audio.SoundEffect, bufferSize int) []float64 {
	synth := audio.Synthesizer{}
	synth.SetSfx(0, e)
	synth.Sfx(0, audio.Channel0, 0, 31)
	buffer := make([]float64, bufferSize)
	synth.ReadSamples(buffer)
	return buffer
}

func assertAllValuesBetween(t *testing.T, minInclusive, maxInclusive float64, buffer []float64) {
	for i, b := range buffer {
		require.Truef(t, b >= minInclusive && b <= maxInclusive, "buffer[%d] is not between [%f,%f]", i, minInclusive, maxInclusive)
	}
}

func dominantFrequency(input []float64) int {
	maxAmplitude := 0.0
	dominantFrequencyIndex := 0

	for i, value := range fft(input) {
		amplitude := cmplx.Abs(value)
		if amplitude > maxAmplitude {
			maxAmplitude = amplitude
			dominantFrequencyIndex = i
		}
	}

	return int(
		float64(dominantFrequencyIndex) * audio.SampleRate / float64(len(input)),
	)
}

// fft runs Fast Fourier Transform on given input.
func fft(input []float64) []complex128 {
	freqs := make([]complex128, len(input))
	hfft(input, freqs, len(input), 1)
	return freqs
}

// code by Dylan Meeus, from GoAudio library: https://github.com/DylanMeeus/GoAudio
func hfft(input []float64, freqs []complex128, n, step int) {
	if n == 1 {
		freqs[0] = complex(input[0], 0)
		return
	}

	h := n / 2

	hfft(input, freqs, h, 2*step)
	hfft(input[step:], freqs[h:], h, 2*step)

	for k := 0; k < h; k++ {
		a := -2 * math.Pi * float64(k) * float64(n)
		e := cmplx.Rect(1, a) * freqs[k+h]
		freqs[k], freqs[k+h] = freqs[k]+e, freqs[k]-e
	}
}
