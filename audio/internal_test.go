package audio

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOscillatorFunc(t *testing.T) {
	t.Run("should not return nil for all instruments", func(t *testing.T) {
		for i := Instrument(0); i < InstrumentSfx7; i++ {
			require.NotNil(t, oscillatorFunc(i))
		}
	})
}

func TestPitchToFreq(t *testing.T) {
	tests := map[Pitch]float64{
		PitchA2:  440,
		PitchC0:  65.40639132515,
		PitchDs5: 2489.0158697766,
		PitchC1:  130.8127826503,
		PitchC2:  261.6255653006,
	}
	for pitch, expectedFreq := range tests {
		actualFreq := pitchToFreq(pitch)
		assert.InDelta(t, expectedFreq, actualFreq, 0.000000001)
	}
}

const durationOfNoteWhenSpeedIsOne = 183

func TestSingleNoteSamples(t *testing.T) {
	assert.Equal(t, durationOfNoteWhenSpeedIsOne, singleNoteSamples(0))
	assert.Equal(t, durationOfNoteWhenSpeedIsOne, singleNoteSamples(1))
	assert.Equal(t, 2*durationOfNoteWhenSpeedIsOne, singleNoteSamples(2))
	assert.Equal(t, 255*durationOfNoteWhenSpeedIsOne, singleNoteSamples(255))
}
