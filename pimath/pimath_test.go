// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pimath_test

import (
	"testing"

	"github.com/elgopher/pi/pimath"
	"github.com/stretchr/testify/assert"
)

func TestClamp(t *testing.T) {
	tests := map[string]struct {
		x, min, max int
		expected    int
	}{
		"lower than min": {
			x: -1, min: 0, max: 1,
			expected: 0,
		},
		"equal to min": {
			x: 0, min: 0, max: 1,
			expected: 0,
		},
		"greater than max": {
			x: 2, min: 0, max: 1,
			expected: 1,
		},
		"equal to max": {
			x: 1, min: 0, max: 1,
			expected: 1,
		},
		"between min and max": {
			x: 1, min: 0, max: 2,
			expected: 1,
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			actual := pimath.Clamp(testCase.x, testCase.min, testCase.max)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestLerp(t *testing.T) {
	tests := map[string]struct {
		a, b, t  float64
		expected float64
	}{
		"in the middle": {
			a: 1, b: 2, t: 0.5,
			expected: 1.5,
		},
		"beginning": {
			a: 1, b: 2, t: 0,
			expected: 1,
		},
		"end": {
			a: 1, b: 2, t: 1,
			expected: 2,
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			actual := pimath.Lerp(testCase.a, testCase.b, testCase.t)
			assert.InDelta(t, testCase.expected, actual, 0.01)
		})
	}
}

func TestDistance(t *testing.T) {
	tests := map[string]struct {
		x1, y1   float64
		x2, y2   float64
		expected float64
	}{
		"(0,0) to (2,2)": {
			x1: 0, y1: 0,
			x2: 2, y2: 2,
			expected: 2.83,
		},
		"(2,2) to (0,0)": {
			x1: 2, y1: 2,
			x2: 0, y2: 0,
			expected: 2.83,
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			actual := pimath.Distance(testCase.x1, testCase.y1, testCase.x2, testCase.y2)
			assert.InDelta(t, testCase.expected, actual, 0.01)
		})
	}
}
