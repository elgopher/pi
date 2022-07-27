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

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed internal/testimage/sprite-sheet-16x16.png
var spriteSheet16x16 []byte

//go:embed internal/testimage/*.png
var images embed.FS

func TestCls(t *testing.T) {
	t.Run("should clean screen using color 0", func(t *testing.T) {
		pi.Reset()
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.ScreenData = []byte{1, 2, 3, 4}
		// when
		pi.Cls()
		// then
		assert.Equal(t, []byte{0, 0, 0, 0}, pi.ScreenData)
	})
}

func TestClsCol(t *testing.T) {
	t.Run("should clean screen using color 0", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.ScreenData = []byte{1, 2, 3, 4}
		// when
		pi.ClsCol(7)
		// then
		assert.Equal(t, []byte{7, 7, 7, 7}, pi.ScreenData)
	})
}

func TestPset(t *testing.T) {
	col := byte(2)

	t.Run("should set color of pixel", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		// when
		pi.Color = col
		pi.Pset(1, 1)
		// then
		assert.Equal(t, col, pi.ScreenData[3])
	})

	t.Run("should not set pixel outside the screen", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()

		emptyScreen := make([]byte, len(pi.ScreenData))

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
				pi.Color = col
				pi.Pset(coords.X, coords.Y)
				// then
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})

	t.Run("should not set pixel outside the clipping region", func(t *testing.T) {
		pi.ScreenWidth = 4
		pi.ScreenHeight = 4
		emptyScreen := make([]byte, 16)
		tests := []struct{ X, Y int }{
			{0, 1},
			{1, 0},
			{2, 1},
			{1, 2},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Clip(1, 1, 1, 1)
				// when
				pi.Color = col
				pi.Pset(coords.X, coords.Y)
				// then
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})

	t.Run("should set pixel taking camera position into consideration", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.Camera(1, 2)
		pi.Color = 8
		// when
		pi.Pset(1, 2)
		// then
		expected := make([]byte, 4)
		expected[0] = 8
		assert.Equal(t, expected, pi.ScreenData)
	})

	t.Run("should not set pixel outside the screen when camera is set", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
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
				pi.BootOrPanic()
				pi.Color = col
				// when
				pi.Camera(1, 1)
				pi.Pset(coords.X, coords.Y)
				// then
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})
}

func TestPget(t *testing.T) {
	t.Run("should get color of pixel", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		col := byte(7)
		pi.Color = col
		pi.Pset(1, 1)
		// expect
		assert.Equal(t, col, pi.Pget(1, 1))
	})

	t.Run("should get color 0 if outside the screen", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.ClsCol(7)

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
				actual := pi.Pget(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})

	t.Run("should get color 0 if outside the clipping region", func(t *testing.T) {
		pi.ScreenWidth = 4
		pi.ScreenHeight = 4

		tests := []struct{ X, Y int }{
			{0, 1},
			{1, 0},
			{2, 1},
			{1, 2},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Pset(coords.X, coords.Y)
				pi.Clip(1, 1, 1, 1)
				// when
				actual := pi.Pget(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})

	t.Run("should get pixel taking camera position into consideration", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.Camera(1, 2)
		pi.Color = 8
		pi.Pset(1, 2)
		// when
		actual := pi.Pget(1, 2)
		// then
		assert.Equal(t, pi.Color, actual)
	})

	t.Run("should get color 0 for pixels outside the screen when camera is set", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
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
				pi.BootOrPanic()
				pi.ClsCol(7)
				pi.Camera(1, 1)
				// when
				actual := pi.Pget(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})
}
func TestClip(t *testing.T) {
	t.Run("should return entire screen by default", func(t *testing.T) {
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		pi.BootOrPanic()
		x, y, w, h := pi.Clip(1, 2, 3, 4)
		assert.Zero(t, x)
		assert.Zero(t, y)
		assert.Equal(t, pi.ScreenWidth, w)
		assert.Equal(t, pi.ScreenHeight, h)
	})

	t.Run("should return previous clipping region", func(t *testing.T) {
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		pi.BootOrPanic()
		pi.Clip(1, 2, 3, 4)
		x, y, w, h := pi.Clip(5, 6, 7, 8)
		assert.Equal(t, 1, x)
		assert.Equal(t, 2, y)
		assert.Equal(t, 3, w)
		assert.Equal(t, 4, h)
	})

	t.Run("should clip with entire screen", func(t *testing.T) {
		tests := map[clippingRegion]clippingRegion{
			{-1, 0, 7, 7}: {0, 0, 6, 7},
			{0, -1, 7, 7}: {0, 0, 7, 6},
			{0, 0, 9, 8}:  {0, 0, 8, 8},
			{0, 0, 8, 9}:  {0, 0, 8, 8},
			{1, 0, 8, 8}:  {1, 0, 7, 8},
			{0, 1, 8, 8}:  {0, 1, 8, 7},
		}
		for given, expected := range tests {
			t.Run(fmt.Sprintf("%+v", given), func(t *testing.T) {
				pi.ScreenWidth = 8
				pi.ScreenHeight = 8
				pi.BootOrPanic()
				pi.Clip(given.x, given.y, given.w, given.h)
				x, y, w, h := pi.ClipReset()
				assert.Equal(t, expected.x, x)
				assert.Equal(t, expected.y, y)
				assert.Equal(t, expected.w, w)
				assert.Equal(t, expected.h, h)
			})
		}
	})
}

func TestClipReset(t *testing.T) {
	t.Run("should return previous clip", func(t *testing.T) {
		pi.BootOrPanic()
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
		pi.BootOrPanic()
		pi.Clip(1, 2, 3, 4)
		// when
		pi.ClipReset()
		// then
		x, y, w, h := pi.ClipReset()
		assert.Equal(t, x, 0)
		assert.Equal(t, y, 0)
		assert.Equal(t, w, pi.ScreenWidth)
		assert.Equal(t, h, pi.ScreenHeight)
	})
}

func TestCamera(t *testing.T) {
	t.Run("should return initial camera", func(t *testing.T) {
		pi.BootOrPanic()
		initialX, initialY := pi.Camera(1, 2)
		assert.Equal(t, 0, initialX)
		assert.Equal(t, 0, initialY)
	})

	t.Run("should return previous camera", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Camera(1, 2)
		x, y := pi.Camera(2, 3)
		assert.Equal(t, 1, x)
		assert.Equal(t, 2, y)
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
				pi.ScreenWidth = 8
				pi.ScreenHeight = 8
				pi.SpriteSheetWidth = 16
				pi.SpriteSheetHeight = 16
				pi.BootOrPanic()
				pi.ClsCol(7)
				snapshot := clone(pi.ScreenData)
				// when
				spr(spriteNo, 0, 0)
				// then
				assert.Equal(t, snapshot, pi.ScreenData)
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
				pi.ScreenWidth = 8
				pi.ScreenHeight = 8
				pi.Resources = fstest.MapFS{
					"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
				}
				pi.BootOrPanic()
				expectedScreen := decodePNG(t, "internal/testimage/"+test.expectedScreenFile)
				// when
				pi.Camera(test.cameraX, test.cameraY)
				spr(test.spriteNo, test.x, test.y)
				// then
				assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
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
				pi.ScreenWidth = 8
				pi.ScreenHeight = 8
				pi.Resources = fstest.MapFS{
					"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
				}
				pi.BootOrPanic()
				expectedScreen := decodePNG(t, "internal/testimage/"+test.expectedScreenFile)
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				spr(0, 0, 0)
				// then
				assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
			})
		}
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
				pi.ScreenWidth = 8
				pi.ScreenHeight = 8
				pi.SpriteSheetWidth = 16
				pi.SpriteSheetHeight = 16
				pi.BootOrPanic()
				pi.ClsCol(7)
				snapshot := clone(pi.ScreenData)
				// when
				sprSize(0, 0, 0, test.w, test.h)
				// then
				assert.Equal(t, snapshot, pi.ScreenData)
			})
		}
	})

	t.Run("should draw sprite", func(t *testing.T) {
		tests := map[string]struct {
			w, h               float64
			x, y               int
			expectedScreenFile string
		}{
			"sprite at (0,0,0.5,1.0)": {w: 0.5, h: 1.0, expectedScreenFile: "spr_0_0_0.5_1.0.png"},
			"sprite at (0,0,1.0,0.5)": {w: 1.0, h: 0.5, expectedScreenFile: "spr_0_0_1.0_0.5.png"},
			"sprite at (0,0,2.0,1.0)": {w: 2.0, h: 1.0, expectedScreenFile: "spr_0_0_2.0_1.0.png"},
			"sprite at (0,0,1.0,2.0)": {w: 1.0, h: 2.0, expectedScreenFile: "spr_0_0_1.0_2.0.png"},
			"sprite at (0,0,2.5,1.0)": {w: 2.5, h: 1.0, expectedScreenFile: "spr_0_0_2.0_1.0.png"},
			"sprite at (0,0,1.0,2.5)": {w: 1.0, h: 2.5, expectedScreenFile: "spr_0_0_1.0_2.0.png"},
			"sprite at (8,0,2.0,1.0)": {x: 8, w: 2.0, h: 1.0, expectedScreenFile: "spr_8_0_2.0_1.0.png"},
			"sprite at (0,8,1.0,2.0)": {y: 8, w: 1.0, h: 2.0, expectedScreenFile: "spr_0_8_1.0_2.0.png"},
			"sprite at (0,0,1.1,0.5)": {w: 1.1, h: 0.5, expectedScreenFile: "spr_0_0_1.0_0.5.png"}, // should floor(w*8)
			"sprite at (0,0,0.5,1.1)": {w: 0.5, h: 1.1, expectedScreenFile: "spr_0_0_0.5_1.0.png"}, // should floor(h*8)
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.ScreenWidth = 16
				pi.ScreenHeight = 16
				pi.Resources = fstest.MapFS{
					"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
				}
				pi.BootOrPanic()
				expectedScreen := decodePNG(t, "internal/testimage/"+test.expectedScreenFile)
				// when
				sprSize(0, test.x, test.y, test.w, test.h)
				// then
				assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
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
				pi.ScreenWidth = 8
				pi.ScreenHeight = 8
				pi.Resources = fstest.MapFS{
					"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet16x16},
				}
				pi.BootOrPanic()
				expectedScreen := decodePNG(t, "internal/testimage/"+test.expectedScreenFile)
				// when
				pi.SprSizeFlip(0, 0, 0, 1.0, test.h, test.flipX, test.flipY)
				// then
				assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
			})
		}
	})
}

func clone(s []byte) []byte {
	cloned := make([]byte, len(s))
	copy(cloned, s)
	return cloned
}

type clippingRegion struct {
	x, y, w, h int
}

func decodePNG(t *testing.T, file string) image.Image {
	payload, err := images.ReadFile(file)
	require.NoError(t, err)
	data, err := image.DecodePNG(bytes.NewReader(payload))
	require.NoError(t, err)
	return data
}
