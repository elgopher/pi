// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/elgopher/pi"
	"github.com/stretchr/testify/assert"
)

func TestRectFill(t *testing.T) {
	testRect(t, pi.RectFill, "rectfill")
}

func testRect(t *testing.T, rect func(x0, y0, x1, y1 int), dir string) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16

	const white = 7

	t.Run("should draw rectangle", func(t *testing.T) {
		tests := map[string]struct {
			x0, y0, x1, y1 int
			color          byte
		}{
			"0,0,0,0":         {0, 0, 0, 0, white},
			"0,0,0,0 color 8": {0, 0, 0, 0, 8},
			"1,0,1,0":         {1, 0, 1, 0, white},
			"0,1,0,1":         {0, 1, 0, 1, white},
			"1,1,1,1":         {1, 1, 1, 1, white},
			"15,15,15,15":     {15, 15, 15, 15, white},
			"0,0,1,0":         {0, 0, 1, 0, white},
			"0,0,0,1":         {0, 0, 0, 1, white},
			"0,0,1,1":         {0, 0, 1, 1, white},
			"3,4,5,6":         {3, 4, 5, 6, white},
			"1,0,0,0":         {1, 0, 0, 0, white},
			"0,1,0,0":         {0, 1, 0, 0, white},
			"1,1,0,0":         {1, 1, 0, 0, white},
			"6,5,4,3":         {6, 5, 4, 3, white},
			"-1,0,0,0":        {-1, 0, 0, 0, white},
			"0,-1,0,0":        {0, -1, 0, 0, white},
			"16,0,15,0":       {16, 0, 15, 0, white},
			"17,0,15,0":       {17, 0, 15, 0, white},
			"0,16,0,15":       {0, 16, 0, 15, white},
			"0,17,0,15":       {0, 17, 0, 15, white},
			"15,16,15,15":     {15, 16, 15, 15, white},
			"15,17,15,15":     {15, 17, 15, 15, white},
			"-1,-1,16,16":     {-1, -1, 16, 16, white},
			"0,0,15,15":       {0, 0, 15, 15, white},
			"0,-1,15,15":      {0, -1, 15, 15, white}, // Rect specific - missing top line
			"0,0,15,16":       {0, 0, 15, 16, white},  // Rect specific - missing bottom line
			"-1,0,15,15":      {-1, 0, 15, 15, white}, // Rect specific - missing left line
			"0,0,16,15":       {0, 0, 16, 15, white},  // Rect specific - missing right line
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.ClsCol(5)
				pi.Color(test.color)
				rect(test.x0, test.y0, test.x1, test.y1)
				assertScreenEqual(t, "internal/testimage/"+dir+"/draw/"+name+".png")
			})
		}
	})

	t.Run("should draw inside clipping region", func(t *testing.T) {
		tests := map[string]struct {
			x0, y0, x1, y1             int
			clipX, clipY, clipW, clipH int
		}{
			"clipx_0,0,2,1": {
				x1: 2, y1: 1,
				clipX: 1, clipW: 16, clipH: 16,
			},
			"clipw_0,0,15,1": {
				x1: 15, y1: 1,
				clipW: 15, clipH: 16,
			},
			"clipxw_0,0,15,1": {
				x1: 15, y1: 1,
				clipX: 1, clipW: 14, clipH: 16,
			},
			"clipy_0,0,15,1": {
				x1: 15, y1: 1,
				clipY: 1, clipW: 16, clipH: 16,
			},
			"cliph_0,0,15,15": {
				x1: 15, y1: 15,
				clipW: 16, clipH: 15,
			},
			"clipyh_0,0,15,15": {
				x1: 15, y1: 15,
				clipY: 1, clipW: 16, clipH: 14,
			},
			"clipall_0,0,15,15": {
				x1: 15, y1: 15,
				clipX: 1, clipY: 1, clipW: 14, clipH: 14,
			},
			"clipxw_22_3,0,3,0": {
				x0: 3, y0: 0, x1: 3, y1: 0, // X > clip width, but still should be drawn (because clip X is>0)
				clipX: 2, clipW: 2, clipH: 16,
			},
			"clipyh_22_0,3,0,3": {
				x0: 0, y0: 3, x1: 0, y1: 3, // Y > clip height, but still should be drawn (because clip Y is>0)
				clipY: 2, clipW: 16, clipH: 2,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				pi.Color(white)
				rect(test.x0, test.y0, test.x1, test.y1)
				assertScreenEqual(t, "internal/testimage/"+dir+"/clip/"+name+".png")
			})
		}
	})

	t.Run("should not draw anything outside clipping region", func(t *testing.T) {
		tests := map[string]struct {
			clipX, clipY, clipW, clipH int
			x0, y0, x1, y1             int
		}{
			"both x < clip x": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: -1, y0: 0, x1: -1, y1: 0,
			},
			"both x < clip x - 1": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: -2, y0: 0, x1: -2, y1: 0,
			},
			"both x > clip width": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: 17, y0: 0, x1: 17, y1: 0,
			},
			"both x==clip width": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: 16, y0: 0, x1: 16, y1: 0,
			},
			"both y < clip y": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: 0, y0: -1, x1: 0, y1: -1,
			},
			"both y > clip height": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: 0, y0: 17, x1: 0, y1: 17,
			},
			"both y==clip height": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: 0, y0: 16, x1: 0, y1: 16,
			},
			"both y==clip height,both x==1": {
				clipX: 0, clipY: 0, clipW: 16, clipH: 16,
				x0: 1, y0: 16, x1: 1, y1: 16,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.BootOrPanic()
				pi.Color(white)
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				rect(test.x0, test.y0, test.x1, test.y1)
				// then
				emptyScreen := make([]byte, len(pi.ScreenData))
				assert.Equal(t, emptyScreen, pi.ScreenData)
			})
		}
	})

	t.Run("should move by camera position", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(white)
		pi.Camera(-2, -3)
		rect(0, 1, 2, 4)
		assertScreenEqual(t, "internal/testimage/"+dir+"/camera_0,1,2,4.png")
	})

	t.Run("should replace color from draw palette", func(t *testing.T) {
		pi.BootOrPanic()
		pi.Color(white)
		pi.Pal(white, 3)
		rect(5, 5, 10, 10)
		assertScreenEqual(t, "internal/testimage/"+dir+"/pal_5,5,10,10.png")
	})
}

func TestRect(t *testing.T) {
	testRect(t, pi.Rect, "rect")
}
