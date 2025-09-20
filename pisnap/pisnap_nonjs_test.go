// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package pisnap_test

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pisnap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaptureOrErr(t *testing.T) {
	// when
	file, err := pisnap.CaptureOrErr()
	// then
	require.NoError(t, err)
	assert.NotEmpty(t, file)
	// and file exists
	f, err := os.ReadFile(file)
	require.NoError(t, err, "cannot read PNG file")
	// and
	img, err := png.Decode(bytes.NewReader(f))
	require.NoError(t, err, "file is not a valid PNG")
	// and
	palettedImage, ok := img.(*image.Paletted)
	require.True(t, ok, "image is not a Paletted")
	assertPalettedImage(t, palettedImage, pi.Screen())
}
