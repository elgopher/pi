// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine //nolint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/audio"
)

func TestEbitenPlayerSource_Read(t *testing.T) {
	t.Run("should read 0 samples when buffer is empty", func(t *testing.T) {
		reader := &ebitenPlayerSource{}
		n, err := reader.Read([]byte{})
		require.NoError(t, err)
		assert.Equal(t, n, 0)
	})

	t.Run("should convert floats to linear PCM (signed 16bits little endian, 2 channel stereo).", func(t *testing.T) {
		reader := &ebitenPlayerSource{
			audioSystem: &audioSystemMock{buffer: []float64{-1, 0, 1, -0.5, 0.5}},
		}
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
		reader := &ebitenPlayerSource{
			audioSystem: &audioSystemMock{buffer: []float64{1, -0.5, 0.5}},
		}
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
		reader := &ebitenPlayerSource{
			audioSystem: &audioSystemMock{buffer: []float64{-2, 2}},
		}
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
		reader := &ebitenPlayerSource{
			audioSystem: &audioSystemMock{buffer: []float64{1, -1}},
		}
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
		mock := &audioSystemMock{buffer: []float64{0, 1, -1, 0}}
		reader := &ebitenPlayerSource{
			audioSystem: mock,
		}
		n, err := reader.Read(make([]byte, 8)) // read 2 samples first
		require.NoError(t, err)
		require.Equal(t, 8, n)
		n, err = reader.Read(make([]byte, 4)) // read 1 sample only
		require.NoError(t, err)
		require.Equal(t, 4, n)
		assert.Len(t, mock.buffer, 1, "one float in the buffer should still be available for reading")
	})
}

type audioSystemMock struct {
	buffer []float64
}

func (m *audioSystemMock) ReadSamples(buffer []float64) {
	n := copy(buffer, m.buffer)
	m.buffer = m.buffer[n:]
}

func (m *audioSystemMock) Sfx(sfxNo int, channel audio.Channel, offset, length int) {}
func (m *audioSystemMock) Music(patterNo int, fadeMs int, channelMask byte)         {}

func (m *audioSystemMock) Stat() audio.Stat {
	return audio.Stat{}
}

func (m *audioSystemMock) SetSfx(sfxNo int, e audio.SoundEffect) {}
func (m *audioSystemMock) GetSfx(sfxNo int) audio.SoundEffect {
	return audio.SoundEffect{}
}

func (m *audioSystemMock) SetMusic(patternNo int, _ audio.Pattern) {}
func (m *audioSystemMock) GetMusic(patterNo int) audio.Pattern {
	return audio.Pattern{}
}

func (m *audioSystemMock) Save() ([]byte, error) {
	return nil, nil
}

func (m *audioSystemMock) Load(bytes []byte) error {
	return nil
}
