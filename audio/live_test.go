// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/audio"
)

const sampleDuration = time.Second / audio.SampleRate

func TestLiveReader_ReadSamples(t *testing.T) {
	t.Run("should not read anything on first call", func(t *testing.T) {
		c := newClock()
		r := audio.LiveReader{
			BufferSize:      time.Second,
			ReadSamplesFunc: readSamplesSeq(),
			Now:             c.Now,
		}
		buf := make([]float64, 1)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Zero(t, n)
		assertSilence(t, buf)
	})

	t.Run("should read one sample", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReader(c)
		c.advance(sampleDuration)
		buf := make([]float64, 1)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{1}, buf)
	})

	t.Run("should read two samples", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReader(c)
		c.advance(2 * sampleDuration)
		buf := make([]float64, 2)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 2, n)
		assert.Equal(t, []float64{1, 2}, buf)
	})

	t.Run("should read only one sample even though buffer is bigger because its to early", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReader(c)
		c.advance(sampleDuration) // only one sample elapsed, so only one sample will be read
		buf := make([]float64, 2)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{1, 0}, buf)
	})

	t.Run("should read samples up to the buffer size", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReader(c)
		c.advance(2 * sampleDuration) // two samples elapsed, but buffer is not big enough to read all
		buf := make([]float64, 1)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{1}, buf)
	})

	t.Run("should read all samples in two calls", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReader(c)
		c.advance(2 * sampleDuration) // two samples elapsed, but buffer is not big enough to read all
		buf := make([]float64, 1)
		r.ReadSamples(buf)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{2}, buf)
	})

	t.Run("should not read anything when its to early and previously samples were read", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReader(c)
		c.advance(sampleDuration)
		r.ReadSamples(make([]float64, 1))
		buf := make([]float64, 1)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Zero(t, n)
		assertSilence(t, buf)
	})

	t.Run("should drop sample when its too late", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReaderWithBufferSize(c, sampleDuration)
		c.advance(2 * sampleDuration) // advance two samples, but because buffer size is one sample then first sample will be dropped
		buf := make([]float64, 1)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{2}, buf)
	})

	t.Run("should continue reading samples after dropping some", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReaderWithBufferSize(c, sampleDuration)
		c.advance(2 * sampleDuration)     // advance two samples, but because buffer size is one sample then first sample will be dropped
		r.ReadSamples(make([]float64, 1)) // sample = 2
		buf := make([]float64, 1)
		c.advance(sampleDuration)
		// when
		n := r.ReadSamples(buf) // should read sample = 3 (should not drop any more samples)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{3}, buf)
	})

	t.Run("should drop huge number of samples", func(t *testing.T) {
		c := newClock()
		r := prepareLiveReaderWithBufferSize(c, sampleDuration)
		c.advance(65537 * sampleDuration) // 65536 samples to drop
		buf := make([]float64, 1)
		// when
		n := r.ReadSamples(buf)
		// then
		assert.Equal(t, 1, n)
		assert.Equal(t, []float64{65537}, buf)
	})
}

type clock struct {
	now time.Time
}

func newClock() *clock {
	return &clock{now: time.UnixMilli(1)}
}

func (c *clock) Now() time.Time {
	return c.now
}

func (c *clock) advance(duration time.Duration) {
	c.now = c.now.Add(duration)
}

func readSamplesSeq() func(buf []float64) {
	n := 1.0
	return func(buf []float64) {
		for i := 0; i < len(buf); i++ {
			buf[i] = n
			n++
		}
	}
}

func prepareLiveReader(c *clock) *audio.LiveReader {
	return prepareLiveReaderWithBufferSize(c, time.Second)
}

func prepareLiveReaderWithBufferSize(c *clock, bufferSize time.Duration) *audio.LiveReader {
	r := &audio.LiveReader{
		BufferSize:      bufferSize,
		ReadSamplesFunc: readSamplesSeq(),
		Now:             c.Now,
	}
	r.ReadSamples(make([]float64, 0)) // first read does not read samples
	return r
}
