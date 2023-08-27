package audio

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOscillatorFunc(t *testing.T) {
	t.Run("should not return nil for all instruments", func(t *testing.T) {
		for i := Instrument(0); i < InstrumentSfx7; i++ {
			require.NotNil(t, oscillatorFunc(i))
		}
	})
}
