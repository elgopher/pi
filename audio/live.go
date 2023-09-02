// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"fmt"
	"time"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio/internal"
)

const sampleDuration = time.Second / internal.SampleRate

// LiveReader is used by backend to read samples live taking into consideration current time.
// It drops samples when ReadSamples was called too late, and limits how much samples can be read.
type LiveReader struct {
	BufferSize      time.Duration
	ReadSamplesFunc func([]float64)
	Now             func() time.Time

	started      time.Time
	lastSample   int
	reusedBuffer []float64
}

// ReadSamples reads samples and writes them to buf. Returns how much samples was written to buffer.
func (l *LiveReader) ReadSamples(buf []float64) int {
	now := l.Now()
	if l.started.IsZero() {
		l.started = now
		return 0
	}

	maxRealSample := int(now.Sub(l.started) / sampleDuration)
	maxCallerSample := l.lastSample + len(buf)

	maxSample := pi.MinInt(maxRealSample, maxCallerSample)
	samplesToRead := maxSample - l.lastSample

	samplesToDrop := maxRealSample - maxCallerSample - int(l.BufferSize/sampleDuration) + 1
	if samplesToDrop > 0 {
		l.dropSamples(samplesToDrop)
		maxSample += samplesToDrop
	}

	l.ReadSamplesFunc(buf[:samplesToRead])

	l.lastSample = maxSample

	return samplesToRead
}

func (l *LiveReader) dropSamples(samplesToDrop int) {
	fmt.Printf("dropping %d audio samples\n", samplesToDrop)

	if l.reusedBuffer == nil {
		l.reusedBuffer = make([]float64, 1024)
	}

	for samplesToDrop > 0 {
		n := pi.MinInt(samplesToDrop, len(l.reusedBuffer))
		l.ReadSamplesFunc(l.reusedBuffer[:n])
		samplesToDrop -= n
	}
}
