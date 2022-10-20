// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/vm"
)

func TestRectFill(t *testing.T) {
	testRect(t, pi.RectFill, "rectfill")
}

func testRect(t *testing.T, rect func(x0, y0, x1, y1 int, color byte), dir string) {
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
				pi.MustBoot()
				pi.ClsCol(5)
				rect(test.x0, test.y0, test.x1, test.y1, test.color)
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
				pi.MustBoot()
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				rect(test.x0, test.y0, test.x1, test.y1, white)
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
				pi.MustBoot()
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				rect(test.x0, test.y0, test.x1, test.y1, white)
				// then
				emptyScreen := make([]byte, len(vm.ScreenData))
				assert.Equal(t, emptyScreen, vm.ScreenData)
			})
		}
	})

	t.Run("should move by camera position", func(t *testing.T) {
		pi.MustBoot()
		pi.Camera(-2, -3)
		rect(0, 1, 2, 4, white)
		assertScreenEqual(t, "internal/testimage/"+dir+"/camera_0,1,2,4.png")
	})

	t.Run("should replace color from draw palette", func(t *testing.T) {
		pi.MustBoot()
		pi.Pal(white, 3)
		rect(5, 5, 10, 10, white)
		assertScreenEqual(t, "internal/testimage/"+dir+"/pal_5,5,10,10.png")
	})
}

func TestRect(t *testing.T) {
	testRect(t, pi.Rect, "rect")
}

func TestLine(t *testing.T) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	const white, red = 7, 8

	t.Run("should not draw anything outside clipping region", func(t *testing.T) {
		tests := map[string]struct {
			clipX, clipY, clipW, clipH int
			x0, y0, x1, y1             int
		}{
			"both x<clipX": {
				clipW: 16, clipH: 16,
				x0: -1, x1: -1,
			},
			"both x==clipX+W": {
				clipX: 1, clipW: 15, clipH: 16,
				x0: 16, x1: 16,
			},
			"both x>clipX+W": {
				clipX: 1, clipW: 15, clipH: 16,
				x0: 17, x1: 17,
			},
			"both y<clipY": {
				clipW: 16, clipH: 16,
				y0: -1, y1: -1,
			},
			"both y==clipY+H": {
				clipY: 1, clipW: 15, clipH: 16,
				y0: 16, y1: 16,
			},
			"both y>clipY+H": {
				clipY: 1, clipW: 15, clipH: 16,
				y0: 17, y1: 17,
			},
			"horizontal,both x<clipX": {
				clipW: 16, clipH: 16,
				x0: -2, x1: -1,
			},
			"horizontal,both x==clipX+W": {
				clipX: 1, clipW: 14, clipH: 16,
				x0: 15, y0: 0, x1: 16, y1: 0,
			},
			"horizontal,both x>clipX+W": {
				clipX: 1, clipW: 14, clipH: 16,
				x0: 16, y0: 0, x1: 17, y1: 0,
			},
			"horizontal line,y<clipY": {
				clipY: 1, clipW: 16, clipH: 16,
				y0: -1, y1: -1, x0: 0, x1: 2,
			},
			"horizontal line,y==clipY+H": {
				clipY: 1, clipW: 16, clipH: 14,
				y0: 15, y1: 15, x0: 0, x1: 2,
			},
			"horizontal line,y>clipY+H": {
				clipY: 1, clipW: 16, clipH: 14,
				y0: 16, y1: 16, x0: 0, x1: 2,
			},
			"slope 1,both x<clipX": {
				clipW: 16, clipH: 16,
				x0: -2, y0: 0, x1: -1, y1: 1,
			},
			"slope 1,y<clipY": {
				clipY: 1, clipW: 16, clipH: 16,
				y0: -3, y1: -1, x0: 0, x1: 2,
			},
			"slope 1,y0==clipY+H": {
				clipY: 1, clipW: 16, clipH: 14,
				y0: 15, y1: 17, x0: 0, x1: 2,
			},
			"slope 1,y0>clipY+H": {
				clipY: 1, clipW: 16, clipH: 14,
				y0: 16, y1: 18, x0: 0, x1: 2,
			},
			"slope 1,x0==clipX+W": {
				clipX: 1, clipW: 14, clipH: 16,
				x0: 15, y0: 0, x1: 17, y1: 2,
			},
			"slope 1,x0>clipX+W": {
				clipX: 1, clipW: 14, clipH: 16,
				x0: 16, y0: 0, x1: 18, y1: 2,
			},
			"slope 2,both x<clipX": {
				clipW: 16, clipH: 16,
				x0: -2, y0: 0, x1: -1, y1: 2,
			},
			"slope 2,y<clipY": {
				clipY: 1, clipW: 16, clipH: 16,
				y0: -5, y1: -1, x0: 0, x1: 2,
			},
			"slope 2,y0==clipY+H": {
				clipY: 1, clipW: 16, clipH: 14,
				y0: 15, y1: 19, x0: 0, x1: 2,
			},
			"slope 2,y0>clipY+H": {
				clipY: 1, clipW: 16, clipH: 14,
				y0: 16, y1: 20, x0: 0, x1: 2,
			},
			"slope 2,x0==clipX+W": {
				clipX: 1, clipW: 14, clipH: 16,
				x0: 15, y0: 0, x1: 17, y1: 4,
			},
			"slope 2,x0>clipX+W": {
				clipX: 1, clipW: 14, clipH: 16,
				x0: 16, y0: 0, x1: 18, y1: 4,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				pi.Line(test.x0, test.y0, test.x1, test.y1, white)
				// then
				emptyScreen := make([]byte, len(vm.ScreenData))
				assert.Equal(t, emptyScreen, vm.ScreenData)
			})
		}
	})

	t.Run("should draw line", func(t *testing.T) {
		tests := map[string]struct {
			x0, y0, x1, y1 int
			color          byte
		}{
			"slope 0, white":          {0, 0, 0, 0, white},
			"slope 0, red":            {0, 0, 0, 0, red},
			"vertical line":           {0, 0, 0, 1, white},
			"vertical line inverse":   {0, 1, 0, 0, white},
			"vertical line red":       {0, 0, 0, 1, red},
			"horizontal line":         {0, 0, 1, 0, white},
			"horizontal line inverse": {1, 0, 0, 0, white},
			"horizontal line red":     {0, 0, 1, 0, red},
			"slope 1":                 {0, 0, 2, 2, white},
			"slope -1":                {0, 2, 2, 0, white},
			"slope 0.5":               {0, 0, 2, 1, white},
			"slope -0.5":              {0, 1, 2, 0, white},
			"slope 2":                 {0, 0, 2, 4, white},
			"slope -2":                {0, 4, 2, 0, white},
			"slope 2.5":               {0, 0, 2, 5, white},
			"slope -2.5":              {0, 5, 2, 0, white},
			"slope 1.5":               {0, 3, 15, 13, white},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				pi.ClsCol(5)
				// when
				pi.Line(test.x0, test.y0, test.x1, test.y1, test.color)
				assertScreenEqual(t, "internal/testimage/line/draw/"+name+".png")
			})
		}
	})

	t.Run("should draw inside clipping region", func(t *testing.T) {
		tests := map[string]struct {
			clipX, clipY, clipW, clipH int
			x0, y0, x1, y1             int
		}{
			"top-left": {
				clipX: 1, clipW: 15, clipH: 16,
				x0: 1, x1: 1,
			},
			"top-right": {
				clipX: 1, clipW: 15, clipH: 16,
				x0: 15, x1: 15,
			},
			"top": {
				clipY: 1, clipW: 16, clipH: 15,
				y0: 1, y1: 1,
			},
			"bottom": {
				clipY: -1, clipW: 16, clipH: 16,
				y0: 14, y1: 14,
			},
			"vertical line": {
				clipW: 16, clipH: 16,
				y0: -1, y1: 16,
			},
			"horizontal line": {
				clipW: 16, clipH: 16,
				x0: -1, x1: 17,
			},
			"horizontal line, clipx=1": {
				clipX: 1, clipW: 2, clipH: 16,
				x0: 1, x1: 2,
			},
			"horizontal line, bottom": {
				clipW: 16, clipH: 16,
				x0: 0, y0: 15, x1: 15, y1: 15,
			},
			"slope 1": {
				clipW: 16, clipH: 16,
				x0: -1, y0: -1, x1: 16, y1: 16,
			},
			"slope -1": {
				clipW: 16, clipH: 16,
				x0: 16, y0: -1, x1: -1, y1: 16,
			},
			"slope 2": { // different from Pico-8
				clipW: 16, clipH: 16,
				x0: -1, y0: -1, x1: 8, y1: 17,
			},
			"slope -2": { // different from Pico-8
				clipW: 16, clipH: 16,
				x0: 8, y0: -1, x1: -1, y1: 17,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				pi.ClsCol(5)
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				// when
				pi.Line(test.x0, test.y0, test.x1, test.y1, white)
				assertScreenEqual(t, "internal/testimage/line/clip/"+name+".png")
			})
		}
	})

	t.Run("should use draw palette", func(t *testing.T) {
		tests := map[string]struct {
			x0, y0, x1, y1 int
		}{
			"horizontal line": {
				x1: 15,
			},
			"vertical line": {
				y1: 15,
			},
			"slope 1": {
				x1: 15, y1: 15,
			},
			"slope 2": {
				x1: 15, y1: 30,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				pi.ClsCol(5)
				pi.Pal(white, red)
				// when
				pi.Line(test.x0, test.y0, test.x1, test.y1, white)
				assertScreenEqual(t, "internal/testimage/line/pal/"+name+".png")
			})
		}
	})

	t.Run("should move by camera position", func(t *testing.T) {
		tests := map[string]struct {
			x0, y0, x1, y1 int
		}{
			"horizontal line": {
				x1: 15,
			},
			"vertical line": {
				y1: 15,
			},
			"slope 1": {
				x1: 15, y1: 15,
			},
			"slope 2": {
				x1: 15, y1: 30,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				pi.ClsCol(5)
				pi.Camera(-1, -2)
				// when
				pi.Line(test.x0, test.y0, test.x1, test.y1, red)
				assertScreenEqual(t, "internal/testimage/line/camera/"+name+".png")
			})
		}
	})
}

func TestCirc(t *testing.T) {
	testCirc(t, pi.Circ, "circ")
}

func testCirc(t *testing.T, circ func(x, y, r int, color byte), dir string) {
	pi.ScreenWidth = 16
	pi.ScreenHeight = 16
	const white, red = 7, 8

	t.Run("should draw", func(t *testing.T) {
		tests := map[string]struct {
			x, y, r int
			color   byte
		}{
			"radius -1":       {8, 8, -1, white},
			"radius 0, white": {8, 8, 0, white},
			"radius 1, white": {8, 8, 1, white},
			"radius 1, red":   {8, 8, 1, red},
			"radius 2":        {8, 8, 2, white},
			"radius 3":        {8, 8, 3, white},
			"radius 4":        {8, 8, 4, white},
			"radius 5":        {8, 8, 5, white},
			"radius 6":        {8, 8, 6, white},
			"radius 7":        {8, 8, 7, white},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				pi.ClsCol(5)
				// when
				circ(test.x, test.y, test.r, test.color)
				assertScreenEqual(t, "internal/testimage/"+dir+"/draw/"+name+".png")
			})
		}
	})

	t.Run("should use draw palette", func(t *testing.T) {
		pi.MustBoot()
		pi.ClsCol(5)
		pi.Pal(white, red)
		// when
		circ(8, 8, 1, white)
		assertScreenEqual(t, "internal/testimage/"+dir+"/draw/radius 1, red.png")
	})

	t.Run("should move by camera position", func(t *testing.T) {
		pi.MustBoot()
		pi.ClsCol(5)
		pi.Camera(2, 1)
		// when
		circ(8, 8, 5, white)
		assertScreenEqual(t, "internal/testimage/"+dir+"/camera.png")
	})

	t.Run("should draw inside clipping region", func(t *testing.T) {
		tests := map[string]struct {
			clipX, clipY, clipW, clipH int
			x, y, r                    int
		}{
			"left": {
				clipW: 16, clipH: 16,
				x: 4, y: 5, r: 5,
			},
			"top": {
				clipW: 16, clipH: 16,
				x: 5, y: 4, r: 5,
			},
			"right": {
				clipW: 16, clipH: 16,
				x: 11, y: 5, r: 5,
			},
			"bottom": {
				clipW: 16, clipH: 16,
				x: 5, y: 11, r: 5,
			},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				pi.MustBoot()
				pi.ClsCol(5)
				// when
				pi.Clip(test.clipX, test.clipY, test.clipW, test.clipH)
				circ(test.x, test.y, test.r, white)
				assertScreenEqual(t, "internal/testimage/"+dir+"/clip/"+name+".png")
			})
		}
	})
}

func TestCircFill(t *testing.T) {
	testCirc(t, pi.CircFill, "circfill")
}
