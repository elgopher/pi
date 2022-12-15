// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package audio

import (
	"io"
	"time"
)

// Oscillator generates finite PCM signal.
type Oscillator struct {
	frequency    uint16
	duration     time.Duration
	waveForm     WaveForm
	currentFrame int
}

// Read returns [io.EOF] when tone stopped playing
func (o *Oscillator) Read(p []float64) (n int, err error) {
	return 0, io.EOF
}

// SetFrequency changes the frequency if needed. Duration or wave form is not altered.
func (o *Oscillator) SetFrequency(freq uint16) {}

// SetDuration changes the duration if needed. Frequency or wave form is not altered.
func (o *Oscillator) SetDuration(duration time.Duration) {}

// SetWaveForm changes the wave form if needed. Frequency or duration is not altered.
//
// On next Read new wave will be generated.
//
// Please note that to check if wave form was actually changed only wave form names
// are compared (in Go you can't compare functions).
func (o *Oscillator) SetWaveForm(f WaveForm) {
}

type WaveForm struct {
	Name string
	F    func(time float64) float64
}

// Reset resets the object to initial state (no frequency and duration set, Read will return EOF)
//
// Use this method to avoid creating a new Oscillator object each time you want to play a note.
func (o *Oscillator) Reset() {
	o.frequency = 0
	o.duration = 0
	o.currentFrame = 0
	o.waveForm = WaveForm{}
}
