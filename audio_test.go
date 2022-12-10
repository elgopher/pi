// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"sync"
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

func TestAudio(t *testing.T) {
	t.Run("should not data race when run with -race flag", func(t *testing.T) {
		audio := pi.Audio()
		const goroutineCount = 100
		var wg sync.WaitGroup
		wg.Add(goroutineCount)
		for i := 0; i < goroutineCount; i++ {
			go func() {
				audio.Set(0, pi.SoundEffect{})
				_ = audio.Get(0)
				audio.Play(0, 0, 0, 0)
				_, err := audio.Read(make([]byte, 8192))
				require.NoError(t, err)
				audio.Reset()
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func TestAudioSystem_Set(t *testing.T) {
	t.Run("should set sound effect", func(t *testing.T) {
		given := pi.SoundEffect{Speed: 1}
		given.Notes[0] = pi.Note{Pitch: 1, Volume: 2}
		given.Notes[31] = pi.Note{Pitch: 3, Volume: 4}
		// when
		pi.Audio().Set(0, given)
		// then
		actual := pi.Audio().Get(0)
		assert.Equal(t, given, actual)
		assert.NotSame(t, given.Notes, actual.Notes)
	})
}

func TestAudioSystem_Play(t *testing.T) {
	effect := pi.SoundEffect{
		Speed: 1, // speed 1 is ~8.333 ms
		Notes: [32]pi.Note{
			{Pitch: 127, Volume: 255},
		},
	}

	t.Run("should generate audio stream", func(t *testing.T) {
		audio := pi.Audio()
		audio.Reset()
		audio.Set(0, effect)
		// when
		audio.Play(0, 0, 0, 1)
		// then
		buffer := make([]byte, 1500)
		_, err := audio.Read(buffer)
		require.NoError(t, err)
		assert.NotEqual(t, make([]byte, 1500), buffer)
	})

	t.Run("should not generate audio stream when channel is higher than max", func(t *testing.T) {
		audio := pi.Audio()
		audio.Reset()
		audio.Set(0, effect)
		// when
		audio.Play(0, 4, 0, 1)
		// then
		buffer := make([]byte, 1500)
		_, err := audio.Read(buffer)
		require.NoError(t, err)
		assert.Equal(t, make([]byte, 1500), buffer)
	})
}
