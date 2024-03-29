// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
)

var (
	//go:embed internal/testimage/sprite-sheet-16x16.png
	spriteSheet16x16 []byte
	//go:embed internal/testimage/invalid-sheet-width.png
	invalidSpriteSheetWidth []byte
	//go:embed internal/testimage/invalid-sheet-height.png
	invalidSpriteSheetHeight []byte
	//go:embed internal/testimage/*
	images embed.FS
)

func TestPixMap_Clear(t *testing.T) {
	t.Run("should clear PixMap using color 0", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		copy(pixMap.Pix(), []byte{1, 2, 3, 4})
		// when
		pixMap.Clear()
		// then
		assert.Equal(t, []byte{0, 0, 0, 0}, pixMap.Pix())
	})
}

func TestPixMap_ClearCol(t *testing.T) {
	t.Run("should clean PixMap using color 7", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		copy(pixMap.Pix(), []byte{1, 2, 3, 4})
		// when
		pixMap.ClearCol(7)
		// then
		assert.Equal(t, []byte{7, 7, 7, 7}, pixMap.Pix())
	})
}

func TestPixMap_Set(t *testing.T) {
	const col byte = 2

	t.Run("should set color of pixel", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		// when
		pixMap.Set(1, 1, col)
		// then
		assert.Equal(t, col, pixMap.Pix()[3])
	})

	t.Run("should not set pixel outside the screen", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)

		tests := []struct{ X, Y int }{
			{-1, 0},
			{0, -1},
			{2, 0},
			{0, 2},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				// when
				pixMap.Set(coords.X, coords.Y, col)
				// then
				assertEmptyPixMap(t, pixMap)
			})
		}
	})

	t.Run("should not set pixel outside the clipping region", func(t *testing.T) {
		tests := []struct{ X, Y int }{
			{0, 1},
			{1, 0},
			{2, 1},
			{1, 2},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pixMap := pi.NewPixMap(4, 4).WithClip(1, 1, 1, 1)
				// when
				pixMap.Set(coords.X, coords.Y, col)
				// then
				assertEmptyPixMap(t, pixMap)
			})
		}
	})

	t.Run("should not set pixel outside the clipping region (x,y higher than w,h)", func(t *testing.T) {
		tests := []struct{ X, Y int }{
			{1, 2}, {2, 2}, {3, 2},
			{1, 3} /*   */, {3, 3},
			{1, 4} /*   */, {3, 4},
			{1, 5}, {2, 5}, {3, 5},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pixMap := pi.NewPixMap(8, 8).WithClip(2, 3, 1, 2)
				// when
				pixMap.Set(coords.X, coords.Y, col)
				// then
				assertEmptyPixMap(t, pixMap)
			})
		}
	})

	t.Run("should set pixel inside the clipping region", func(t *testing.T) {
		pixMap := pi.NewPixMap(8, 8).WithClip(2, 3, 1, 1)
		// when
		pixMap.Set(2, 3, col)
		// then
		assertNotEmptyPixMap(t, pixMap)
	})
}

func TestSet(t *testing.T) {
	const col byte = 2

	t.Run("should set pixel taking camera position into account", func(t *testing.T) {
		pi.Reset()
		pi.SetScreenSize(2, 2)
		pi.Camera.Set(1, 2)
		// when
		pi.Set(1, 2, 8)
		// then
		expected := make([]byte, 4)
		expected[0] = 8
		assert.Equal(t, expected, pi.Scr().Pix())
	})

	t.Run("should not set pixel outside the screen when camera is set", func(t *testing.T) {
		emptyScreen := make([]byte, 4)
		tests := []struct{ X, Y int }{
			{0, 0},
			{1, 0},
			{0, 1},
			{2, 0},
			{0, 2},
			{3, 0},
			{3, 1},
			{3, 2},
			{0, 3},
			{1, 3},
			{2, 3},
			{3, 3},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pi.Reset()
				pi.SetScreenSize(2, 2)
				// when
				pi.Camera.Set(1, 1)
				pi.Set(coords.X, coords.Y, col)
				// then
				assert.Equal(t, emptyScreen, pi.Scr().Pix())
			})
		}
	})

	t.Run("should draw swapped color", func(t *testing.T) {
		pi.Reset()
		pi.SetScreenSize(1, 1)
		pi.Pal[1] = 2
		// when
		pi.Set(0, 0, 1)
		// then
		assert.Equal(t, []byte{2}, pi.Scr().Pix())
	})

	t.Run("should draw original color after resetting draw palette", func(t *testing.T) {
		pi.Reset()
		pi.SetScreenSize(1, 1)
		pi.Pal[1] = 2
		pi.Pal.Reset()
		// when
		pi.Set(0, 0, 1)
		// then
		assert.Equal(t, []byte{1}, pi.Scr().Pix())
	})
}

func TestPixMap_Get(t *testing.T) {
	t.Run("should get color of pixel", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		col := byte(7)
		pixMap.Set(1, 1, col)
		// expect
		assert.Equal(t, col, pixMap.Get(1, 1))
	})

	t.Run("should get color 0 if outside the screen", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		pixMap.ClearCol(7)

		tests := []struct{ X, Y int }{
			{-1, 0},
			{0, -1},
			{2, 0},
			{0, 2},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				// when
				actual := pixMap.Get(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})

	t.Run("should get color 0 if outside the clipping region", func(t *testing.T) {
		tests := []struct{ X, Y int }{
			{0, 1},
			{1, 0},
			{2, 1},
			{1, 2},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pixMap := pi.NewPixMap(4, 4).WithClip(1, 1, 1, 1)
				pixMap.Set(coords.X, coords.Y, 7)
				// when
				actual := pixMap.Get(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})

	t.Run("should get color 0 if outside the clipping region (x,y higher than w,h)", func(t *testing.T) {
		tests := []struct{ X, Y int }{
			{1, 2}, {2, 2}, {3, 2},
			{1, 3} /*   */, {3, 3},
			{1, 4} /*   */, {3, 4},
			{1, 5}, {2, 5}, {3, 5},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pixMap := pi.NewPixMap(8, 8).WithClip(2, 3, 1, 2)
				pixMap.Set(coords.X, coords.Y, 7)
				// when
				actual := pixMap.Get(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})

	t.Run("should get pixel inside the clipping region", func(t *testing.T) {
		pixMap := pi.NewPixMap(8, 8).WithClip(2, 3, 1, 1)
		col := byte(6)
		pixMap.Set(2, 3, col)
		// when
		actual := pixMap.Get(2, 3)
		// then
		assert.Equal(t, col, actual)
	})
}

func TestGet(t *testing.T) {
	t.Run("should get pixel taking camera position into consideration", func(t *testing.T) {
		pi.Reset()
		pi.SetScreenSize(2, 2)
		pi.Camera.Set(1, 2)
		const color byte = 8
		pi.Set(1, 2, color)
		// when
		actual := pi.Get(1, 2)
		// then
		assert.Equal(t, color, actual)
	})

	t.Run("should get color 0 for pixels outside the screen when camera is set", func(t *testing.T) {
		tests := []struct{ X, Y int }{
			{0, 0},
			{1, 0},
			{0, 1},
			{2, 0},
			{0, 2},
			{3, 0},
			{3, 1},
			{3, 2},
			{0, 3},
			{1, 3},
			{2, 3},
			{3, 3},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pi.Reset()
				pi.SetScreenSize(2, 2)
				pi.ClsCol(7)
				pi.Camera.Set(1, 1)
				// when
				actual := pi.Get(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})
}

func TestClip(t *testing.T) {
	t.Run("should return entire screen by default", func(t *testing.T) {
		pi.Reset()
		pi.SetScreenSize(8, 8)
		x, y, w, h := pi.Clip(1, 2, 3, 4)
		assert.Zero(t, x)
		assert.Zero(t, y)
		assert.Equal(t, pi.Scr().Width(), w)
		assert.Equal(t, pi.Scr().Height(), h)
	})

	t.Run("should return previous clipping region", func(t *testing.T) {
		pi.Reset()
		pi.SetScreenSize(8, 8)
		pi.Clip(1, 2, 3, 4)
		x, y, w, h := pi.Clip(5, 6, 7, 8)
		assert.Equal(t, 1, x)
		assert.Equal(t, 2, y)
		assert.Equal(t, 3, w)
		assert.Equal(t, 4, h)
	})
}

func TestPixMap_WithClip(t *testing.T) {
	t.Run("should clip with entire pixmap", func(t *testing.T) {
		tests := map[pi.Region]pi.Region{
			{-1, 0, 7, 7}: {0, 0, 6, 7},
			{0, -1, 7, 7}: {0, 0, 7, 6},
			{0, 0, 9, 8}:  {0, 0, 8, 8},
			{0, 0, 8, 9}:  {0, 0, 8, 8},
			{1, 0, 8, 8}:  {1, 0, 7, 8},
			{0, 1, 8, 8}:  {0, 1, 8, 7},
		}
		for given, expected := range tests {
			t.Run(fmt.Sprintf("%+v", given), func(t *testing.T) {
				pixMap := pi.NewPixMap(8, 8)
				// when
				pixMap = pixMap.WithClip(given.X, given.Y, given.W, given.H)
				// then
				assert.Equal(t, pi.Region{X: expected.X, Y: expected.Y, W: expected.W, H: expected.H}, pixMap.Clip())
			})
		}
	})
}

func TestClipReset(t *testing.T) {
	t.Run("should return previous clip", func(t *testing.T) {
		pi.Reset()
		pi.Clip(1, 2, 3, 4)
		// when
		x, y, w, h := pi.ClipReset()
		// then
		assert.Equal(t, x, 1)
		assert.Equal(t, y, 2)
		assert.Equal(t, w, 3)
		assert.Equal(t, h, 4)
	})

	t.Run("should reset clip to full screen size", func(t *testing.T) {
		pi.Reset()
		pi.Clip(1, 2, 3, 4)
		// when
		pi.ClipReset()
		// then
		x, y, w, h := pi.ClipReset()
		assert.Equal(t, x, 0)
		assert.Equal(t, y, 0)
		assert.Equal(t, w, pi.Scr().Width())
		assert.Equal(t, h, pi.Scr().Height())
	})
}

func TestCamera(t *testing.T) {
	t.Run("should return initial camera", func(t *testing.T) {
		pi.Reset()
		assert.Equal(t, 0, pi.Camera.X)
		assert.Equal(t, 0, pi.Camera.Y)
	})
}

func TestSpr(t *testing.T) {
	testSpr(t, pi.Spr)
}

func testSpr(t *testing.T, spr func(spriteNo int, x int, y int)) {
	t.Run("should not draw not existing sprite", func(t *testing.T) {
		tests := map[string]int{
			"negative sprite":        -1,
			"sprite number too high": 4,
		}

		for name, spriteNo := range tests {
			t.Run(name, func(t *testing.T) {
				pi.Reset()
				pi.SetScreenSize(8, 8)
				pi.UseEmptySpriteSheet(16, 16)
				pi.ClsCol(7)
				snapshot := clone(pi.Scr().Pix())
				// when
				spr(spriteNo, 0, 0)
				// then
				assert.Equal(t, snapshot, pi.Scr().Pix())
			})
		}
	})

	t.Run("should draw sprite", func(t *testing.T) {
		tests := map[string]struct {
			spriteNo           int
			x, y               int
			cameraX, cameraY   int
			expectedScreenFile string
		}{
			"sprite 0 at (0,0)":               {spriteNo: 0, expectedScreenFile: "spr_0_at_00.png"},
			"sprite 1 at (0,0)":               {spriteNo: 1, expectedScreenFile: "spr_1_at_00.png"},
			"sprite 3 at (0,0)":               {spriteNo: 3, expectedScreenFile: "spr_3_at_00.png"},
			"sprite 0 at (1,0)":               {x: 1, expectedScreenFile: "spr_0_at_10.png"},
			"sprite 0 at (0,1)":               {y: 1, expectedScreenFile: "spr_0_at_01.png"},
			"sprite 0 at (-1,0)":              {x: -1, expectedScreenFile: "spr_0_at_-10.png"},
			"sprite 0 at (0,-1)":              {y: -1, expectedScreenFile: "spr_0_at_0-1.png"},
			"sprite 0 at (8,0)":               {x: 8, expectedScreenFile: "zeros_8x8.png"},
			"sprite 0 at (0,8)":               {y: 8, expectedScreenFile: "zeros_8x8.png"},
			"sprite 0 at (-8,0)":              {x: -8, expectedScreenFile: "zeros_8x8.png"},
			"sprite 0 at (0,-8)":              {y: -8, expectedScreenFile: "zeros_8x8.png"},
			"sprite 0 at (0,0), camera (1,0)": {cameraX: 1, expectedScreenFile: "spr_0_at_00_cam_10.png"},
			"sprite 0 at (0,0), camera (0,1)": {cameraY: 1, expectedScreenFile: "spr_0_at_00_cam_01.png"},
			"sprite 0 at (9,0)":               {x: 9, expectedScreenFile: "zeros_8x8.png"},
			"sprite 0 at (0,9)":               {y: 9, expectedScreenFile: "zeros_8x8.png"},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				boot(8, 8, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
				// when
				pi.Camera.Set(test.cameraX, test.cameraY)
				spr(test.spriteNo, test.x, test.y)
				// then
				assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
			})
		}
	})

	t.Run("should draw sprite with clipping", func(t *testing.T) {
		tests := map[string]struct {
			clipX, clipY, clipW, clipH int
			expectedScreenFile         string
		}{
			"sprite 0 at (0,0), clip (1,0,7,8)": {clipX: 1, clipW: 7, clipH: 8, expectedScreenFile: "spr_0_at_00_clip1078.png"},
			"sprite 0 at (0,0), clip (0,1,8,7)": {clipY: 1, clipW: 8, clipH: 7, expectedScreenFile: "spr_0_at_00_clip0187.png"},
			"sprite 0 at (0,0), clip (0,0,7,8)": {clipY: 0, clipW: 7, clipH: 8, expectedScreenFile: "spr_0_at_00_clip0078.png"},
			"sprite 0 at (0,0), clip (0,0,8,7)": {clipY: 0, clipW: 8, clipH: 7, expectedScreenFile: "spr_0_at_00_clip0087.png"},
			"sprite 0 at (0,0), clip (2,0,7,8)": {clipX: 2, clipW: 7, clipH: 8, expectedScreenFile: "spr_0_at_00_clip2078.png"},
			"sprite 0 at (0,0), clip (0,2,8,7)": {clipY: 2, clipW: 8, clipH: 7, expectedScreenFile: "spr_0_at_00_clip0287.png"},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				boot(8, 8, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				spr(0, 0, 0)
				// then
				assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
			})
		}
	})

	t.Run("should not draw color 0 by default", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		spr(2, 0, 0)
		// when
		spr(1, 0, 0)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_0.png")
		assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
	})

	t.Run("should not draw color 0 after reset", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		spr(2, 0, 0)
		pi.Palt[0] = false // make color 0 opaque
		// when
		pi.Palt.Reset() // and then make color 0 transparent again
		spr(1, 0, 0)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_0.png")
		assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
	})

	t.Run("should not draw transparent colors", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		spr(2, 0, 0)
		// when
		pi.Palt[0] = false
		pi.Palt[50] = true
		spr(1, 0, 0)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_50.png")
		assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
	})

	t.Run("should swap color", func(t *testing.T) {
		pi.Reset()
		pi.UseEmptySpriteSheet(8, 8)
		pi.SetScreenSize(8, 8)
		const originalColor byte = 7
		const replacementColor byte = 15
		pi.SprSheet().Set(5, 5, originalColor)
		pi.Pal[originalColor] = replacementColor
		// when
		spr(0, 0, 0)
		// then
		actual := pi.Get(5, 5)
		assert.Equal(t, replacementColor, actual)
	})

	t.Run("should draw original color after reset", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_0_at_00.png")
		// when
		pi.Pal[1] = 3
		pi.Pal[28] = 30
		pi.Pal.Reset()
		spr(0, 0, 0)
		// then
		assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
	})
}

func TestSprSize(t *testing.T) {
	testSprSize(t, pi.SprSize)
}

func testSprSize(t *testing.T, sprSize func(spriteNo int, x, y int, w, h float64)) {
	testSpr(t, func(spriteNo int, x int, y int) {
		sprSize(spriteNo, x, y, 1.0, 1.0)
	})

	t.Run("should not draw sprite", func(t *testing.T) {
		tests := map[string]struct {
			w, h float64
		}{
			"w=0":  {w: 0, h: 1},
			"h=0":  {w: 1, h: 0},
			"w=-1": {w: -1, h: 1},
			"h=-1": {w: 1, h: -1},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.SetScreenSize(8, 8)
				pi.UseEmptySpriteSheet(16, 16)
				pi.ClsCol(7)
				snapshot := clone(pi.Scr().Pix())
				// when
				sprSize(0, 0, 0, test.w, test.h)
				// then
				assert.Equal(t, snapshot, pi.Scr().Pix())
			})
		}
	})

	t.Run("should draw sprite", func(t *testing.T) {
		tests := map[string]struct {
			spriteNo           int
			w, h               float64
			x, y               int
			expectedScreenFile string
		}{
			"sprite at (0,0,0.5,1.0)":   {w: 0.5, h: 1.0, expectedScreenFile: "spr_0_0_0.5_1.0.png"},
			"sprite at (0,0,1.0,0.5)":   {w: 1.0, h: 0.5, expectedScreenFile: "spr_0_0_1.0_0.5.png"},
			"sprite at (0,0,2.0,1.0)":   {w: 2.0, h: 1.0, expectedScreenFile: "spr_0_0_2.0_1.0.png"},
			"sprite at (0,0,1.0,2.0)":   {w: 1.0, h: 2.0, expectedScreenFile: "spr_0_0_1.0_2.0.png"},
			"sprite at (0,0,2.5,1.0)":   {w: 2.5, h: 1.0, expectedScreenFile: "spr_0_0_2.0_1.0.png"},
			"sprite at (0,0,1.0,2.5)":   {w: 1.0, h: 2.5, expectedScreenFile: "spr_0_0_1.0_2.0.png"},
			"sprite at (8,0,2.0,1.0)":   {x: 8, w: 2.0, h: 1.0, expectedScreenFile: "spr_8_0_2.0_1.0.png"},
			"sprite at (0,8,1.0,2.0)":   {y: 8, w: 1.0, h: 2.0, expectedScreenFile: "spr_0_8_1.0_2.0.png"},
			"sprite at (0,0,1.1,0.5)":   {w: 1.1, h: 0.5, expectedScreenFile: "spr_0_0_1.0_0.5.png"}, // should floor(w*8)
			"sprite at (0,0,0.5,1.1)":   {w: 0.5, h: 1.1, expectedScreenFile: "spr_0_0_0.5_1.0.png"}, // should floor(h*8)
			"sprite 1 at (0,0,2.0,1.0)": {spriteNo: 1, w: 2.0, h: 1.0, expectedScreenFile: "spr_1_at_0_0_2.0_1.0.png"},
			"sprite 1 at (0,0,1.9,1.0)": {spriteNo: 1, w: 1.9, h: 1.0, expectedScreenFile: "spr_1_at_0_0_2.0_1.0.png"},
			"sprite 2 at (0,0,1.0,2.0)": {spriteNo: 2, w: 1.0, h: 2.0, expectedScreenFile: "spr_2_at_0_0_1.0_2.0.png"},
			"sprite 2 at (0,0,1.0,1.9)": {spriteNo: 2, w: 1.0, h: 1.9, expectedScreenFile: "spr_2_at_0_0_1.0_2.0.png"},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				boot(16, 16, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
				// when
				sprSize(test.spriteNo, test.x, test.y, test.w, test.h)
				// then
				assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
			})
		}
	})
}

func TestSprSizeFlip(t *testing.T) {
	testSprSize(t, func(spriteNo int, x, y int, w, h float64) {
		pi.SprSizeFlip(spriteNo, x, y, w, h, false, false)
	})

	t.Run("should flip", func(t *testing.T) {
		tests := map[string]struct {
			flipX, flipY       bool
			h                  float64
			expectedScreenFile string
		}{
			"sprite 0 at (0,0), flip y":           {flipY: true, h: 1, expectedScreenFile: "spr_0_at_00_flipy.png"},
			"sprite 0 at (0,0), flip x":           {flipX: true, h: 1, expectedScreenFile: "spr_0_at_00_flipx.png"},
			"sprite 0 at (0,0), flip xy":          {flipX: true, flipY: true, h: 1, expectedScreenFile: "spr_0_at_00_flipxy.png"},
			"sprite 0 at (0,0), no flip":          {h: 1, expectedScreenFile: "spr_0_at_00.png"},
			"sprite 0 at (0,0), height 0, flip y": {flipY: true, h: 0, expectedScreenFile: "zeros_8x8.png"},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				boot(8, 8, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
				// when
				pi.SprSizeFlip(0, 0, 0, 1.0, test.h, test.flipX, test.flipY)
				// then
				assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
			})
		}
	})

	t.Run("should not draw transparent colors", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		pi.SprSizeFlip(2, 0, 0, 1.0, 1.0, true, false)
		// when
		pi.Palt[0] = false
		pi.Palt[50] = true
		pi.SprSizeFlip(1, 0, 0, 1.0, 1.0, true, false)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_50_flipx.png")
		assert.Equal(t, expectedScreen.Pix, pi.Scr().Pix())
	})
}

func TestPosition_Reset(t *testing.T) {
	t.Run("should reset x,y to origin", func(t *testing.T) {
		p := pi.Position{X: 1, Y: 2}
		p.Reset()
		assert.Zero(t, p)
	})
}

func TestPosition_Set(t *testing.T) {
	t.Run("should update x and y simultaneously", func(t *testing.T) {
		var p pi.Position
		p.Set(1, 2)
		assert.Equal(t, pi.Position{X: 1, Y: 2}, p)
	})
}

func clone(s []byte) []byte {
	cloned := make([]byte, len(s))
	copy(cloned, s)
	return cloned
}

func decodePNG(t *testing.T, file string) image.Image {
	payload, err := images.ReadFile(file)
	require.NoError(t, err)
	data, err := image.DecodePNG(bytes.NewReader(payload))
	require.NoError(t, err)
	return data
}

func boot(screenWidth, screenHeight int, spriteSheet []byte) {
	pi.Reset()
	pi.SetScreenSize(screenWidth, screenHeight)
	pi.Load(fstest.MapFS{
		"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet},
	})
}
