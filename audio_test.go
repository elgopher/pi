// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
)

func TestAudioSystem_Read(t *testing.T) {
	t.Run("should read 0 bytes when buffer is empty", func(t *testing.T) {
		n, err := pi.Audio().Read(nil)
		assert.Zero(t, n)
		assert.NoError(t, err)
	})

	t.Run("should clear the buffer with 0 when no channels are used", func(t *testing.T) {
		buffer := []byte{1, 2, 3, 4}
		n, err := pi.Audio().Read(buffer)
		require.NotZero(t, n)
		require.NoError(t, err)
		expected := make([]byte, n)
		assert.Equal(t, expected, buffer)
	})
}
