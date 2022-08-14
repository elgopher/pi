// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"os"
	"testing"
	"testing/fstest"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed internal/testimage/sprite-sheet-16x16.png
	spriteSheet16x16 []byte
	//go:embed internal/testimage/*
	images embed.FS
)

func TestColor(t *testing.T) {
	t.Run("should return default color", func(t *testing.T) {
		pi.BootOrPanic()
		assert.Equal(t, byte(6), pi.Color(2))
	})

	t.Run("should return previous color", func(t *testing.T) {
		pi.BootOrPanic()
		prev := byte(3)
		pi.Color(prev)
		assert.Equal(t, prev, pi.Color(4))
	})
}

func TestColorReset(t *testing.T) {
	t.Run("should reset color to default", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(3)
		pi.ColorReset()
		assert.Equal(t, byte(6), pi.Color(5))
	})

	t.Run("should return previous color", func(t *testing.T) {
		pi.BootOrPanic()
		prev := byte(5)
		pi.Color(prev)
		assert.Equal(t, prev, pi.ColorReset())
	})
}

func TestCls(t *testing.T) {
	pi.Reset()
	pi.ScreenWidth = 2
	pi.ScreenHeight = 2

	t.Run("should clean screen using color 0", func(t *testing.T) {
		pi.BootOrPanic()
		pi.ScreenData = []byte{1, 2, 3, 4}
		// when
		pi.Cls()
		// then
		assert.Equal(t, []byte{0, 0, 0, 0}, pi.ScreenData)
	})

	t.Run("should reset clipping region", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Clip(1, 2, 3, 4)
		// when
		pi.Cls() // clips to 0,0,w,h
		// then
		prevX, prevY, prevW, prevH := pi.ClipReset()
		assert.Zero(t, prevX)
		assert.Zero(t, prevY)
		assert.Equal(t, pi.ScreenWidth, prevW)
		assert.Equal(t, pi.ScreenHeight, prevH)
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
		pi.Color(col)
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
				pi.Color(col)
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
				pi.Color(col)
				pi.Pset(coords.X, coords.Y)
				// then
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})

	t.Run("should not set pixel outside the clipping region (x,y higher than w,h)", func(t *testing.T) {
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		emptyScreen := make([]byte, 8*8)
		tests := []struct{ X, Y int }{
			{1, 2}, {2, 2}, {3, 2},
			{1, 3} /*   */, {3, 3},
			{1, 4} /*   */, {3, 4},
			{1, 5}, {2, 5}, {3, 5},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Clip(2, 3, 1, 2)
				// when
				pi.Color(col)
				pi.Pset(coords.X, coords.Y)
				// then
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})

	t.Run("should set pixel inside the clipping region", func(t *testing.T) {
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		emptyScreen := make([]byte, 8*8)
		pi.BootOrPanic()
		pi.Clip(2, 3, 1, 1)
		// when
		pi.Color(col)
		pi.Pset(2, 3)
		// then
		assert.NotEqual(t, emptyScreen, pi.ScreenData)
	})

	t.Run("should set pixel taking camera position into consideration", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.Camera(1, 2)
		pi.Color(8)
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
				pi.Color(col)
				// when
				pi.Camera(1, 1)
				pi.Pset(coords.X, coords.Y)
				// then
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})

	t.Run("should draw swapped color", func(t *testing.T) {
		pi.ScreenWidth = 1
		pi.ScreenHeight = 1
		pi.BootOrPanic()
		pi.Color(1)
		pi.Pal(1, 2)
		// when
		pi.Pset(0, 0)
		// then
		assert.Equal(t, []byte{2}, pi.ScreenData)
	})

	t.Run("should draw original color after PalReset", func(t *testing.T) {
		pi.ScreenWidth = 1
		pi.ScreenHeight = 1
		pi.BootOrPanic()
		pi.Color(1)
		pi.Pal(1, 2)
		pi.PalReset()
		// when
		pi.Pset(0, 0)
		// then
		assert.Equal(t, []byte{1}, pi.ScreenData)
	})
}

func TestPget(t *testing.T) {
	t.Run("should get color of pixel", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		col := byte(7)
		pi.Color(col)
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

	t.Run("should get color 0 if outside the clipping region (x,y higher than w,h)", func(t *testing.T) {
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		tests := []struct{ X, Y int }{
			{1, 2}, {2, 2}, {3, 2},
			{1, 3} /*   */, {3, 3},
			{1, 4} /*   */, {3, 4},
			{1, 5}, {2, 5}, {3, 5},
		}
		for _, coords := range tests {
			name := fmt.Sprintf("%+v", coords)
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Pset(coords.X, coords.Y)
				pi.Clip(2, 3, 1, 2)
				// when
				actual := pi.Pget(coords.X, coords.Y)
				// then
				assert.Zero(t, actual)
			})
		}
	})

	t.Run("should get pixel inside the clipping region", func(t *testing.T) {
		pi.ScreenWidth = 8
		pi.ScreenHeight = 8
		pi.BootOrPanic()
		col := byte(6)
		pi.Color(col)
		pi.Pset(2, 3)
		pi.Clip(2, 3, 1, 1)
		// when
		actual := pi.Pget(2, 3)
		// then
		assert.Equal(t, col, actual)
	})

	t.Run("should get pixel taking camera position into consideration", func(t *testing.T) {
		pi.ScreenWidth = 2
		pi.ScreenHeight = 2
		pi.BootOrPanic()
		pi.Camera(1, 2)
		pi.Color(8)
		pi.Pset(1, 2)
		// when
		actual := pi.Pget(1, 2)
		// then
		assert.Equal(t, pi.ColorReset(), actual)
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
				boot(8, 8, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
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
				boot(8, 8, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				spr(0, 0, 0)
				// then
				assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
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
		assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
	})

	t.Run("should not draw color 0 after PaltReset", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		spr(2, 0, 0)
		pi.Palt(0, false) // make color 0 opaque
		// when
		pi.PaltReset() // and then make color 0 transparent again
		spr(1, 0, 0)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_0.png")
		assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
	})

	t.Run("should not draw transparent colors", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		spr(2, 0, 0)
		// when
		pi.Palt(0, false)
		pi.Palt(50, true)
		spr(1, 0, 0)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_50.png")
		assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
	})

	t.Run("should swap color", func(t *testing.T) {
		pi.SpriteSheetWidth, pi.SpriteSheetHeight = 8, 8
		pi.ScreenWidth, pi.ScreenHeight = 8, 8
		pi.BootOrPanic()
		const originalColor byte = 7
		const replacementColor byte = 15
		pi.Sset(5, 5, originalColor)
		pi.Pal(originalColor, replacementColor)
		// when
		spr(0, 0, 0)
		// then
		actual := pi.Pget(5, 5)
		assert.Equal(t, replacementColor, actual)
	})

	t.Run("should draw original color after reset", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_0_at_00.png")
		// when
		pi.Pal(1, 3)
		pi.Pal(28, 30)
		pi.PalReset()
		spr(0, 0, 0)
		// then
		assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
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
				boot(8, 8, spriteSheet16x16)
				expectedScreen := decodePNG(t, "internal/testimage/spr/"+test.expectedScreenFile)
				// when
				pi.SprSizeFlip(0, 0, 0, 1.0, test.h, test.flipX, test.flipY)
				// then
				assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
			})
		}
	})

	t.Run("should not draw transparent colors", func(t *testing.T) {
		boot(8, 8, spriteSheet16x16)
		pi.SprSizeFlip(2, 0, 0, 1.0, 1.0, true, false)
		// when
		pi.Palt(0, false)
		pi.Palt(50, true)
		pi.SprSizeFlip(1, 0, 0, 1.0, 1.0, true, false)
		// then
		expectedScreen := decodePNG(t, "internal/testimage/spr/spr_1_on_top_of_2_trans_50_flipx.png")
		assert.Equal(t, expectedScreen.Pixels, pi.ScreenData)
	})
}

func TestSnap(t *testing.T) {
	t.Run("should take screenshot and store it to temp file", func(t *testing.T) {
		pi.Reset()
		pi.BootOrPanic()
		for i := 0; i < len(pi.ScreenData); i++ {
			pi.ScreenData[i] = byte(i % 16) // 16 colors by default
		}
		// when
		screenshot, err := pi.Snap()
		// then
		require.NoError(t, err)
		img := decodeScreenshot(t, screenshot)
		assert.Equal(t, pi.ScreenWidth, img.Width)
		assert.Equal(t, pi.ScreenHeight, img.Height)
		assert.Equal(t, pi.ScreenData, img.Pixels)
		assert.Equal(t, pi.Palette, img.Palette)
	})

	t.Run("should use display palette", func(t *testing.T) {
		pi.Reset()
		pi.ScreenWidth, pi.ScreenHeight = 1, 1
		pi.BootOrPanic()
		original, replacement := byte(1), byte(2)
		pi.PalDisplay(original, replacement) // replace 1 by 2
		pi.Color(original)
		pi.Pset(0, 0)
		screenshot, err := pi.Snap()
		// then
		require.NoError(t, err)
		img := decodeScreenshot(t, screenshot)
		assert.Equal(t, pi.Palette[2], img.Palette[1]) // 1 is replaced by 2
		assert.Equal(t, pi.ScreenData, img.Pixels)
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

func decodeScreenshot(t *testing.T, screenshot string) image.Image {
	file, err := os.Open(screenshot)
	require.NoError(t, err)
	img, err := image.DecodePNG(file)
	require.NoError(t, err)
	return img
}

func boot(screenWidth, screenHeight int, spriteSheet []byte) {
	pi.Reset()
	pi.ScreenWidth = screenWidth
	pi.ScreenHeight = screenHeight
	pi.Resources = fstest.MapFS{
		"sprite-sheet.png": &fstest.MapFile{Data: spriteSheet},
	}
	pi.BootOrPanic()
}
