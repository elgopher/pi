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
