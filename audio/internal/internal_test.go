// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/audio/internal"
)

func TestSamples_Sum(t *testing.T) {
	tests := map[string]struct {
		samples  internal.Samples
		expected float64
	}{
		"all zeroes": {
			samples:  internal.Samples{0, 0, 0, 0},
			expected: 0.0,
		},
		"more than 1.0": {
			samples:  internal.Samples{1.0, 0.5, 0, 0},
			expected: 1.5,
		},
		"less than -1.0": {
			samples:  internal.Samples{-1.0, 0, 0, -0.5},
			expected: -1.5,
		},
		"all different values": {
			samples:  internal.Samples{0.1, 0.2, 0.3, 0.4},
			expected: 1.0,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// when
			sum := test.samples.Sum()
			// then
			assert.InDeltaf(t, test.expected, sum, 0.000000001, "sum should be %f but was %f", test.expected, sum)
		})
	}
}

func TestPitchToFreq(t *testing.T) {
	tests := map[int]float64{
		33: 440,             // a2
		0:  65.40639132515,  // c0
		63: 2489.0158697766, // d#5
		12: 130.8127826503,  // c1
		24: 261.6255653006,  // c2
	}
	for pitch, expectedFreq := range tests {
		actualFreq := internal.PitchToFreq(pitch)
		assert.InDelta(t, expectedFreq, actualFreq, 0.000000001)
	}
}
