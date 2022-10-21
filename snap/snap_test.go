// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package snap_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
	"github.com/elgopher/pi/snap"
	"github.com/elgopher/pi/vm"
)

func TestSnap(t *testing.T) {
	if runtime.GOOS == "js" {
		t.Skip("storing files does not work on js")
		return
	}

	t.Run("should take screenshot and store it to temp file", func(t *testing.T) {
		pi.Reset()
		pi.MustBoot()
		for i := 0; i < len(vm.ScreenData); i++ {
			vm.ScreenData[i] = byte(i % 16) // 16 colors by default
		}
		// when
		screenshot, err := snap.Take()
		// then
		require.NoError(t, err)
		img := decodeScreenshot(t, screenshot)
		assert.Equal(t, pi.ScreenWidth, img.Width)
		assert.Equal(t, pi.ScreenHeight, img.Height)
		assert.Equal(t, vm.ScreenData, img.Pixels)
		assert.Equal(t, vm.Palette, img.Palette)
	})

	t.Run("should use display palette", func(t *testing.T) {
		pi.Reset()
		pi.ScreenWidth, pi.ScreenHeight = 1, 1
		pi.MustBoot()
		original, replacement := byte(1), byte(2)
		pi.PalDisplay(original, replacement) // replace 1 by 2
		pi.Pset(0, 0, original)
		screenshot, err := snap.Take()
		// then
		require.NoError(t, err)
		img := decodeScreenshot(t, screenshot)
		assert.Equal(t, vm.Palette[2], img.Palette[1]) // 1 is replaced by 2
		assert.Equal(t, vm.ScreenData, img.Pixels)
	})
}

func decodeScreenshot(t *testing.T, screenshot string) image.Image {
	file, err := os.Open(screenshot)
	require.NoError(t, err)
	img, err := image.DecodePNG(file)
	require.NoError(t, err)
	return img
}
