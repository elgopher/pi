// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitengine //nolint

import (
	"sync"
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
		audioSystem := &audioSystemMock{buffer: []float64{-1, 0, 1, -0.5, 0.5}}
		reader := &ebitenPlayerSource{
			audioSystem: audioSystem,
			readSamples: audioSystem.ReadSamples,
		}
		actual := make([]byte, 40)
		// when
		n, err := reader.Read(actual)
		// then
		require.NoError(t, err)
		assert.Equal(t, 40, n)
		assert.Equal(t, []byte{
			1, 0x80, // -1, left channel, second byte, 7 bit has sign bit
			1, 0x80, // right channel - copy of left channel
			1, 0x80, // copy of left channel (because Ebitengine is using 44100 but Pi 22050)
			1, 0x80, // copy of right channel
			0, 0, // 0
			0, 0, // 0
			0, 0, // 0
			0, 0, // 0
			0xFF, 0x7F, // 1
			0xFF, 0x7F, // 1
			0xFF, 0x7F, // 1
			0xFF, 0x7F, // 1
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
		}, actual)
	})

	t.Run("should continue reading stream using bigger buffer than before", func(t *testing.T) {
		audioSystem := &audioSystemMock{buffer: []float64{1, -0.5, 0.5}}
		reader := &ebitenPlayerSource{
			audioSystem: audioSystem,
			readSamples: audioSystem.ReadSamples,
		}
		smallBuffer := make([]byte, 8)
		n, err := reader.Read(smallBuffer)
		require.Equal(t, 8, n)
		require.NoError(t, err)
		biggerBuffer := make([]byte, 16)
		// when
		n, err = reader.Read(biggerBuffer)
		// then
		assert.Equal(t, 16, n)
		require.NoError(t, err)
		assert.Equal(t, []byte{
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			1, 0xC0, // -0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
			0xFF, 0x3F, // 0.5
		}, biggerBuffer)
	})

	t.Run("should clamp float values to [-1,1]", func(t *testing.T) {
		audioSystem := &audioSystemMock{buffer: []float64{-2, 2}}
		reader := &ebitenPlayerSource{
			audioSystem: audioSystem,
			readSamples: audioSystem.ReadSamples,
		}
		actual := make([]byte, 16)
		// when
		n, err := reader.Read(actual)
		// then
		require.NoError(t, err)
		assert.Equal(t, 16, n)
		assert.Equal(t, []byte{
			1, 0x80, // -2 clamped to -1
			1, 0x80,
			1, 0x80,
			1, 0x80,
			0xFF, 0x7F, // 2 clamped to 1
			0xFF, 0x7F,
			0xFF, 0x7F,
			0xFF, 0x7F,
		}, actual)
	})

	t.Run("should convert floats even when buffer is too small", func(t *testing.T) {
		audioSystem := &audioSystemMock{buffer: []float64{1, -1}}
		reader := &ebitenPlayerSource{
			audioSystem: audioSystem,
			readSamples: audioSystem.ReadSamples,
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
		audioSystem := &audioSystemMock{buffer: []float64{0, 1, -1, 0}}
		reader := &ebitenPlayerSource{
			audioSystem: audioSystem,
			readSamples: audioSystem.ReadSamples,
		}
		n, err := reader.Read(make([]byte, 16)) // read 2 samples first
		require.NoError(t, err)
		require.Equal(t, 16, n)
		n, err = reader.Read(make([]byte, 8)) // read 1 sample only
		require.NoError(t, err)
		require.Equal(t, 8, n)
		assert.Len(t, audioSystem.buffer, 1, "one float in the buffer should still be available for reading")
	})
}

func TestEbitenPlayerSource_ThreadSafety(t *testing.T) {
	t.Run("should not fail when run with -race flag", func(t *testing.T) {
		reader := &ebitenPlayerSource{
			audioSystem: &audio.Synthesizer{},
			readSamples: func(buf []float64) int {
				return len(buf)
			},
		}

		const goroutines = 100

		var group sync.WaitGroup
		group.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer group.Done()

				_, err := reader.Read(make([]byte, 128))
				require.NoError(t, err)
				reader.Play(0, 0, 0, 0)
				reader.Stop(0)
				reader.StopLoop(0)
				reader.StopChan(0)
				reader.Music(0, 0, 0)
				reader.Stat()
				reader.SetSfx(0, audio.SoundEffect{})
				reader.SetMusic(0, audio.Pattern{})
			}()
		}

		group.Wait()
	})
}

type audioSystemMock struct {
	buffer []float64
}

func (m *audioSystemMock) ReadSamples(buffer []float64) int {
	n := copy(buffer, m.buffer)
	m.buffer = m.buffer[n:]
	return n
}

func (m *audioSystemMock) Play(sfxNo, channel, offset, length int)          {}
func (m *audioSystemMock) Stop(sfxNo int)                                   {}
func (m *audioSystemMock) StopLoop(channel int)                             {}
func (m *audioSystemMock) StopChan(channel int)                             {}
func (m *audioSystemMock) Music(patterNo int, fadeMs int, channelMask byte) {}

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
