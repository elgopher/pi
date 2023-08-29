package audio_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/audio"
)

func TestBuffer_Write(t *testing.T) {
	t.Run("cant write to buffer with zero length", func(t *testing.T) {
		buffer := audio.NewBuffer(0)
		// when
		buffer.Write([]float64{1})
		// then
		out := make([]float64, 1)
		n := buffer.Read(out)
		assert.Equal(t, n, 0)
		assert.Equal(t, []float64{0}, out)
	})

	t.Run("should write to buffer once", func(t *testing.T) {
		tests := map[string]struct {
			bufferSize  uint
			in          []float64
			out         []float64
			expectedN   int
			expectedOut []float64
		}{
			"single sample": {
				bufferSize:  1,
				in:          []float64{1},
				out:         make([]float64, 1),
				expectedN:   1,
				expectedOut: []float64{1},
			},
			"single sample, out bigger than available samples": {
				bufferSize:  1,
				in:          []float64{1},
				out:         make([]float64, 2),
				expectedN:   1,
				expectedOut: []float64{1, 0},
			},
			"single sample, buffer size 2": {
				bufferSize:  2,
				in:          []float64{1},
				out:         make([]float64, 1),
				expectedN:   1,
				expectedOut: []float64{1},
			},
			"two samples, buffer big enough": {
				bufferSize:  2,
				in:          []float64{1, 2},
				out:         make([]float64, 2),
				expectedN:   2,
				expectedOut: []float64{1, 2},
			},
			"two samples, buffer smaller": {
				bufferSize:  1,
				in:          []float64{1, 2},
				out:         make([]float64, 1),
				expectedN:   1,
				expectedOut: []float64{2},
			},
		}

		for name, testCase := range tests {
			t.Run(name, func(t *testing.T) {
				buffer := audio.NewBuffer(testCase.bufferSize)
				// when
				buffer.Write(testCase.in)
				// then
				n := buffer.Read(testCase.out)
				assert.Equal(t, n, testCase.expectedN)
				assert.Equal(t, testCase.expectedOut, testCase.out)
			})
		}
	})

	t.Run("should write to buffer after previous write", func(t *testing.T) {
		tests := map[string]struct {
			bufferSize  uint
			in1         []float64
			in2         []float64
			out         []float64
			expectedN   int
			expectedOut []float64
		}{
			"single sample": {
				bufferSize:  2,
				in1:         []float64{1},
				in2:         []float64{2},
				out:         make([]float64, 2),
				expectedN:   2,
				expectedOut: []float64{1, 2},
			},
		}

		for name, testCase := range tests {
			t.Run(name, func(t *testing.T) {

				buffer := audio.NewBuffer(testCase.bufferSize)
				buffer.Write(testCase.in1)
				// when
				buffer.Write(testCase.in2)
				// then
				n := buffer.Read(testCase.out)
				assert.Equal(t, testCase.expectedN, n)
				assert.Equal(t, testCase.expectedOut, testCase.out)
			})
		}

	})

}
