// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio_test

import (
	_ "embed"
	"github.com/elgopher/pi/piaudio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeRaw(t *testing.T) {
	t.Run("should decode raw sample", func(t *testing.T) {
		raw := []byte{0, 1, 255}
		sample := piaudio.DecodeRaw(raw, 22050)

		assert.Equal(t, []int8{0, 1, -1}, sample.Data())
		assert.Equal(t, uint16(22050), sample.SampleRate())
	})
}

var (
	//go:embed internal/test/valid.wav
	validWAV []byte
	//go:embed internal/test/invalid_16bit.wav
	invalid16bitWAV []byte
	//go:embed internal/test/invalid_stereo.wav
	invalidStereoWAV []byte
	//go:embed internal/test/invalid-sample-rate.wav
	invalidSampleRateWAV []byte
	//go:embed internal/test/invalid-not-pcm.wav
	invalidNotPCM []byte
)

func TestDecodeWavOrErr(t *testing.T) {
	t.Run("should decode valid wav file", func(t *testing.T) {
		sample, err := piaudio.DecodeWavOrErr(validWAV)
		require.NoError(t, err)
		assert.Equal(t, []int8{0, 90, 127, 90, 0, -91, -128, -91}, sample.Data())
		assert.Equal(t, uint16(8363), sample.SampleRate())
	})

	t.Run("should return error when wav is invalid", func(t *testing.T) {
		tests := map[string]struct {
			file        []byte
			expectedErr string
		}{
			"16-bit": {
				file:        invalid16bitWAV,
				expectedErr: "only 8-bit PCM supported, got 16 bits",
			},
			"stereo": {
				file:        invalidStereoWAV,
				expectedErr: "only mono supported, got 2 channels",
			},
			"88200 sample rate": {
				file:        invalidSampleRateWAV,
				expectedErr: "sample rate is too high. Max 48kHz supported, got 88200",
			},
			"not PCM": {
				file:        invalidNotPCM,
				expectedErr: "only PCM supported",
			},
		}

		for testName, testCase := range tests {
			t.Run(testName, func(t *testing.T) {
				sample, err := piaudio.DecodeWavOrErr(testCase.file)
				assert.Nil(t, sample)
				assert.ErrorContains(t, err, testCase.expectedErr)
			})
		}
	})
}
