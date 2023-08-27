// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/audio/internal"
)

const (
	freq = 440.0
	step = freq / internal.SampleRate
)

func TestTriangle(t *testing.T) {
	for i := float64(0); i < 3; i += step {
		fmt.Printf("%f\t%f\n", i, internal.Triangle(i))
	}
}

func TestOrgan(t *testing.T) {
	for i := float64(0); i < 3; i += step {
		fmt.Printf("%f\t%f\n", i, internal.Organ(i))
	}
}

func TestOscillator_Next(t *testing.T) {
	var o internal.Oscillator
	o.FreqHz = 440.0
	o.Func = internal.Organ
	// when
	v := o.NextSample()
	// then
	assert.InDelta(t, internal.Organ(0), v, 0.000000001)
	assert.InDelta(t, 0.019954648526077097, o.Pos, 0.000000001)
}
