package audio_test

import (
	"io"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/internal/audio"
)

const notZeroFrequency = 440

func TestOscillator_Read(t *testing.T) {
	t.Run("new oscillator should return EOF", func(t *testing.T) {
		out := make([]float64, 1)

		var oscillator audio.Oscillator
		n, err := oscillator.Read(out)
		assert.Equal(t, n, 0)
		assert.ErrorIs(t, err, io.EOF)
	})

	t.Run("should return n=0, no error when out buffer is empty", func(t *testing.T) {
		out := make([]float64, 0)

		var oscillator audio.Oscillator
		oscillator.SetDuration(time.Second)
		oscillator.SetFrequency(notZeroFrequency)
		oscillator.SetWaveForm(fixedWaveForm(1.0))

		n, err := oscillator.Read(out)
		assert.Equal(t, n, 0)
		assert.NoError(t, err)
	})

	t.Run("should generate 1 sample", func(t *testing.T) {
		out := make([]float64, 2)

		var oscillator audio.Oscillator
		oscillator.SetDuration(22676 * time.Nanosecond) // 1 sample duration
		oscillator.SetFrequency(notZeroFrequency)
		oscillator.SetWaveForm(fixedWaveForm(1.0))
		// when
		n, err := oscillator.Read(out)
		// then
		require.NoError(t, err)
		assert.Equal(t, n, 1)
		assert.Equal(t, []float64{1.0, 0}, out)
		// when
		n, err = oscillator.Read(out)
		// then
		assert.ErrorIs(t, err, io.EOF)
		assert.Equal(t, 0, n)
	})

	t.Run("should pass phase to wave form function based on frequency", func(t *testing.T) {
		var oscillator audio.Oscillator
		const samples = 12
		oscillator.SetDuration(samples * 22676 * time.Nanosecond)
		oscillator.SetWaveForm(audio.WaveForm{
			F: func(time float64) float64 {
				return time
			},
		})
		oscillator.SetFrequency(44100 / samples * 2)
		out := make([]float64, samples)
		// when
		n, err := oscillator.Read(out)
		// then
		require.NoError(t, err)
		require.Equal(t, samples, n)
		for i := 0; i < samples; i++ {
			require.Less(t, out[i], 2*math.Pi)
			require.GreaterOrEqual(t, out[i], 0.0)
		}
	})

	// 1000ms  - 44100
	//  2.5ms  - 100 (440HZ)
}

func fixedWaveForm(v float64) audio.WaveForm {
	return audio.WaveForm{
		Name: "fixedWaveForm",
		F: func(time float64) float64 {
			return v
		},
	}
}
