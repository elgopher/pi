// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"fmt"
	"time"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/audio/internal"
)

const frameDuration = time.Second / internal.SampleRate

// LiveReader is used by backend to read samples live taking into consideration current time.
// It drops frames when ReadSamples was called too late, and limits how much frames can be read.
type LiveReader struct {
	BufferSize      time.Duration
	ReadSamplesFunc func([]float64)
	Now             func() time.Time

	started      time.Time
	lastFrame    int
	reusedBuffer []float64
}

// ReadSamples reads samples and writes them to buf. Returns how much samples was written to buffer.
func (l *LiveReader) ReadSamples(buf []float64) int {
	now := l.Now()
	if l.started.IsZero() {
		l.started = now
		return 0
	}

	maxRealFrame := int(now.Sub(l.started) / frameDuration)
	maxCallerFrame := l.lastFrame + len(buf)
	maxFrame := pi.MinInt(maxRealFrame, maxCallerFrame)

	framesToDrop := maxRealFrame - maxCallerFrame - int(l.BufferSize/frameDuration) + 1
	if framesToDrop > 0 {
		l.dropFrames(framesToDrop)
	}

	n := maxFrame - l.lastFrame
	l.ReadSamplesFunc(buf[:n])

	l.lastFrame = maxFrame

	return n
}

func (l *LiveReader) dropFrames(framesToDrop int) {
	fmt.Printf("dropping %d audio frames\n", framesToDrop)

	if l.reusedBuffer == nil {
		l.reusedBuffer = make([]float64, 1024)
	}

	for framesToDrop > 0 {
		n := pi.MinInt(framesToDrop, len(l.reusedBuffer))
		l.ReadSamplesFunc(l.reusedBuffer[:n])
		framesToDrop -= n
	}
}
