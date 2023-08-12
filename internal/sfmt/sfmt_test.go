// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package sfmt_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/internal/sfmt"
)

func TestFormatBigSlice(t *testing.T) {
	const sliceLen = 10
	slice := make([]int, sliceLen)
	for i := range slice {
		slice[i] = i
	}

	tests := []struct {
		maxSize  int
		expected string
	}{
		{
			maxSize:  -1,
			expected: "(10)[...]",
		},
		{
			maxSize:  0,
			expected: "(10)[...]",
		},
		{
			maxSize:  1,
			expected: "(10)[0 ...]",
		},
		{
			maxSize:  sliceLen,
			expected: "[0 1 2 3 4 5 6 7 8 9]",
		},
		{
			maxSize:  sliceLen + 1,
			expected: "[0 1 2 3 4 5 6 7 8 9]",
		},
	}
	for _, testCase := range tests {
		testName := fmt.Sprintf("%d", testCase.maxSize)
		t.Run(testName, func(t *testing.T) {
			// when
			s := sfmt.FormatBigSlice(slice, testCase.maxSize)
			assert.Equal(t, testCase.expected, s)
		})
	}
}
