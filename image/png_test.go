// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package image_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/image"
	"github.com/elgopher/pi/vm"
)

//go:embed not-indexed.png
var notIndexedFile []byte

//go:embed indexed.png
var indexedFile []byte

func TestDecodePNG(t *testing.T) {
	t.Run("should return error when reader is nil", func(t *testing.T) {
		img, err := image.DecodePNG(nil)
		require.Error(t, err)
		assert.Empty(t, img)
	})

	t.Run("should return error when png file is corrupted", func(t *testing.T) {
		var corruptedData []byte
		reader := bytes.NewReader(corruptedData)
		img, err := image.DecodePNG(reader)
		require.Error(t, err)
		assert.Empty(t, img)
	})

	t.Run("should return error when png file hasn't got indexed color mode", func(t *testing.T) {
		reader := bytes.NewReader(notIndexedFile)
		img, err := image.DecodePNG(reader)
		require.Error(t, err)
		assert.Empty(t, img)
	})

	t.Run("should decode file", func(t *testing.T) {
		reader := bytes.NewReader(indexedFile)
		img, err := image.DecodePNG(reader)
		require.NoError(t, err)
		assert.Equal(t, 4, img.Width)
		assert.Equal(t, 1, img.Height)
		assert.Equal(t,
			//nolint:govet
			[256]vm.RGB{{0, 0, 0}, {0xff, 0xd4, 0x53}, {0xed, 0x45, 0x9c}, {0x6b, 0xd4, 0x7f}},
			img.Palette)
		assert.Equal(t, []byte{0, 1, 2, 3}, img.Pixels)
	})
}
