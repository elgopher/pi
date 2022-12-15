// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAudioStreamReader_Read(t *testing.T) {
	t.Run("should read 0 samples when buffer is empty", func(t *testing.T) {
		reader := &audioStreamReader{}
		n, err := reader.Read([]byte{})
		require.NoError(t, err)
		assert.Equal(t, n, 0)
	})

	t.Run("should convert floats to linear PCM (signed 16bits little endian, 2 channel stereo).", func(t *testing.T) {
		AudioStream = &fakeAudioStream{buffer: []float64{-1, 0, 1, -0.5, 0.5}}
		reader := &audioStreamReader{}
		actual := make([]byte, 20)
		// when
		n, err := reader.Read(actual)
		// then
		require.NoError(t, err)
		assert.Equal(t, 20, n)
		assert.Equal(t, []byte{
			1, 0x80, // -1, left channel, second byte, 7 bit has sign bit
			1, 0x80, // right channel - copy of left channel
			0, 0, // 0
			0, 0, // 0
			0xFF, 0x7F, // 1
			0xFF, 0x7F, // 1
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
		}, actual)
	})

	t.Run("should continue reading stream using bigger buffer than before", func(t *testing.T) {
		AudioStream = &fakeAudioStream{buffer: []float64{1, -0.5, 0.5}}
		reader := &audioStreamReader{}
		smallBuffer := make([]byte, 4)
		n, err := reader.Read(smallBuffer)
		require.Equal(t, 4, n)
		require.NoError(t, err)
		biggerBuffer := make([]byte, 8)
		// when
		n, err = reader.Read(biggerBuffer)
		// then
		assert.Equal(t, 8, n)
		require.NoError(t, err)
		assert.Equal(t, []byte{
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
		}, biggerBuffer)
	})

	t.Run("should clamp float values to [-1,1]", func(t *testing.T) {
		AudioStream = &fakeAudioStream{buffer: []float64{-2, 2}}
		reader := &audioStreamReader{}
		actual := make([]byte, 8)
		// when
		n, err := reader.Read(actual)
		// then
		require.NoError(t, err)
		assert.Equal(t, 8, n)
		assert.Equal(t, []byte{
			1, 0x80, // -2 clamped to -1
			1, 0x80,
			0xFF, 0x7F, // 2 clamped to 1
			0xFF, 0x7F,
		}, actual)
	})

	t.Run("should convert floats even when buffer is too small", func(t *testing.T) {
		AudioStream = &fakeAudioStream{buffer: []float64{1, -1}}
		reader := &audioStreamReader{}
		actual := make([]byte, 8)

		for i := 0; i < 8; i += 2 {
			n, err := reader.Read(actual[i : i+2])
			require.NoError(t, err)
			assert.Equal(t, 2, n)
		}

		assert.Equal(t, []byte{
			0xFF, 0x7F, // 1
			0xFF, 0x7F,
			1, 0x80, // -1
			1, 0x80,
		}, actual)
	})

	t.Run("should only read the minimum number of floats", func(t *testing.T) {
		fakeStream := fakeAudioStream{buffer: []float64{0, 1, -1, 0}}
		AudioStream = &fakeStream
		reader := &audioStreamReader{}
		n, err := reader.Read(make([]byte, 8)) // read 2 samples first
		require.NoError(t, err)
		require.Equal(t, 8, n)
		n, err = reader.Read(make([]byte, 4)) // read 1 sample only
		require.NoError(t, err)
		require.Equal(t, 4, n)
		assert.Len(t, fakeStream.buffer, 1, "one float in the buffer should still be available for reading")
	})
}

type fakeAudioStream struct {
	buffer []float64
}

func (s *fakeAudioStream) Read(p []float64) (n int, err error) {
	n = copy(p, s.buffer)
	s.buffer = s.buffer[n:]
	return n, nil
}
