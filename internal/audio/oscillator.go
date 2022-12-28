// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"io"
	"math"
	"time"

	"github.com/elgopher/pi"
)

const (
	sampleRate = 44100
)

// Oscillator generates finite PCM signal.
type Oscillator struct {
	durationFrames int
	waveForm       WaveForm
	currentFrame   int
	currentPhase   float64
	phaseStep      float64
}

// Read returns [io.EOF] when tone stopped playing
func (o *Oscillator) Read(out []float64) (n int, err error) {
	if len(out) == 0 {
		return 0, nil
	}

	remaining := o.durationFrames - o.currentFrame
	if remaining == 0 {
		return 0, io.EOF
	}

	length := pi.MinInt(len(out), remaining)
	for i := 0; i < length; i++ {
		out[i] = o.waveForm.F(o.currentPhase)
		o.currentPhase += o.phaseStep
		if o.currentPhase >= 2*math.Pi {
			o.currentPhase -= 2 * math.Pi
		}
	}

	o.currentFrame += length

	return length, nil
}

// SetFrequency changes the frequency if needed. Duration or wave form is not altered.
func (o *Oscillator) SetFrequency(freq uint16) {
	o.phaseStep = (2 * math.Pi * float64(freq)) / sampleRate
}

// SetDuration changes the duration if needed. Frequency or wave form is not altered.
func (o *Oscillator) SetDuration(duration time.Duration) {
	o.durationFrames = int(duration*sampleRate) / int(time.Second)
}

// SetWaveForm changes the wave form if needed. Frequency or duration is not altered.
//
// On next Read new wave will be generated.
//
// Please note that to check if wave form was actually changed only wave form names
// are compared (in Go you can't compare functions).
func (o *Oscillator) SetWaveForm(f WaveForm) {
	o.waveForm = f
}

type WaveForm struct {
	Name string
	F    func(time float64) float64
}

// Reset resets the object to initial state (no frequency and duration set, Read will return EOF)
//
// Use this method to avoid creating a new Oscillator object each time you want to play a note.
func (o *Oscillator) Reset() {
	o.durationFrames = 0
	o.waveForm = WaveForm{}
	o.currentFrame = 0
	o.currentPhase = 0
	o.phaseStep = 0
}
