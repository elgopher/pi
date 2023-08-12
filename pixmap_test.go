// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/snap"
)

func TestNewPixMap(t *testing.T) {
	t.Run("should panic when", func(t *testing.T) {
		tests := map[string]func(){
			"width negative":  func() { pi.NewPixMap(-1, 1) },
			"height negative": func() { pi.NewPixMap(1, -1) },
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Panics(t, test, name)
			})
		}
	})

	t.Run("should not panic", func(t *testing.T) {
		tests := map[string]func(){
			"width 0":  func() { pi.NewPixMap(0, 1) },
			"height 0": func() { pi.NewPixMap(1, 0) },
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.NotPanics(t, test, name)
			})
		}
	})

	t.Run("clipping region should cover entire PixMap", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 3)
		assert.Equal(t, pi.Region{W: 2, H: 3}, pixMap.Clip())
	})

	t.Run("Width() and Height() should return original dimensions", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 3)
		pixMap = pixMap.WithClip(1, 1, 1, 1)
		assert.Equal(t, pixMap.Width(), 2)
		assert.Equal(t, pixMap.Height(), 3)
	})

	t.Run("Pix should return slice having all pixels set to 0", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 3)
		assert.Equal(t, make([]byte, 2*3), pixMap.Pix())
	})
}

func TestNewPixMapWithPixels(t *testing.T) {
	t.Run("should panic when", func(t *testing.T) {
		tests := map[string]func(){
			"length of pixels is not multiple of line width": func() {
				pi.NewPixMapWithPixels(make([]byte, 3), 2)
			},
			"pixels slice length is lower than line width": func() {
				pi.NewPixMapWithPixels(make([]byte, 1), 2)
			},
			"line width is negative": func() {
				pi.NewPixMapWithPixels(make([]byte, 2), -2)
			},
			"line width is zero and pixels slice is not empty": func() {
				pi.NewPixMapWithPixels(make([]byte, 1), 0)
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Panics(t, test, name)
			})
		}
	})

	t.Run("should not panic", func(t *testing.T) {
		tests := map[string]func(){
			"length of pixels is multiple of line width": func() {
				pi.NewPixMapWithPixels(make([]byte, 4), 2)
			},
			"line width is equal to pixels slice length": func() {
				pi.NewPixMapWithPixels(make([]byte, 2), 2)
			},
			"line width is zero and pixels slice length is 0": func() {
				pi.NewPixMapWithPixels(make([]byte, 0), 0)
			},
			"line width is zero and pixels slice is nil": func() {
				pi.NewPixMapWithPixels(nil, 0)
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				test()
			})
		}
	})

	t.Run("should create similar PixMap to NewPixMap", func(t *testing.T) {
		pixMap := pi.NewPixMapWithPixels(make([]byte, 6), 3)
		assert.Equal(t, pi.NewPixMap(3, 2), pixMap)
	})

	t.Run("clipping region should cover entire PixMap", func(t *testing.T) {
		pixMap := pi.NewPixMapWithPixels(make([]byte, 2*3), 2)
		assert.Equal(t, pi.Region{W: 2, H: 3}, pixMap.Clip())
	})

	t.Run("Width() and Height() should return original dimensions", func(t *testing.T) {
		pixMap := pi.NewPixMapWithPixels(make([]byte, 2*3), 2)
		pixMap = pixMap.WithClip(1, 1, 1, 1)
		assert.Equal(t, pixMap.Width(), 2)
		assert.Equal(t, pixMap.Height(), 3)
	})

	t.Run("Pix() should return pixel slice passed to constructor", func(t *testing.T) {
		pix := []byte{1, 2, 3, 4, 5, 6}
		pixMap := pi.NewPixMapWithPixels(pix, 2)
		assert.Equal(t, pix, pixMap.Pix())
	})

	t.Run("should create zero value", func(t *testing.T) {
		tests := map[string]func() pi.PixMap{
			"line width is zero and pixels slice length is 0": func() pi.PixMap {
				return pi.NewPixMapWithPixels(make([]byte, 0), 0)
			},
			"line width is zero and pixels slice is nil": func() pi.PixMap {
				return pi.NewPixMapWithPixels(nil, 0)
			},
		}
		for name, newPixMap := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Zero(t, newPixMap())
			})
		}
	})
}

func TestPixMap_Pointer(t *testing.T) {
	w := 5
	h := 7
	clipX := 1
	clipY := 2
	clipW := w - 2
	clipH := h - 3
	pixMap := pi.NewPixMap(w, h).WithClip(clipX, clipY, clipW, clipH)

	t.Run("should return !ok and zero Pointer", func(t *testing.T) {
		tests := map[string]pointerResult{
			"w=0":        result(pixMap.Pointer(clipX+1, clipY+1, 0, 1)),
			"w=-1":       result(pixMap.Pointer(clipX+1, clipY+1, -1, 1)),
			"h=0":        result(pixMap.Pointer(clipX+1, clipY+1, 1, 0)),
			"h=-1":       result(pixMap.Pointer(clipX+1, clipY+1, 1, -1)),
			"x=clip x+w": result(pixMap.Pointer(clipX+clipW, clipY+1, 1, 1)),
			"x>clip x+w": result(pixMap.Pointer(clipX+clipW+1, clipY+1, 1, 1)),
			"y=clip y+h": result(pixMap.Pointer(clipX+1, clipY+clipH, 1, 1)),
			"y>clip y+h": result(pixMap.Pointer(clipX+1, clipY+clipH+1, 1, 1)),
			"x+w<clip x": result(pixMap.Pointer(clipX-2, clipY+1, 1, 3)),
			"x+w=clip x": result(pixMap.Pointer(clipX-1, clipY+1, 1, 3)),
			"y+h<clip y": result(pixMap.Pointer(clipX+1, clipY-2, 3, 1)),
			"y+h=clip y": result(pixMap.Pointer(clipX+1, clipY-1, 3, 1)),
		}
		for name, res := range tests {
			t.Run(name, func(t *testing.T) {
				assert.False(t, res.ok)
				assert.Zero(t, res.pointer)
			})
		}
	})

	t.Run("should return ok and non zero Pointer", func(t *testing.T) {
		tests := map[string]pointerResult{
			"w=1,h=1":      result(pixMap.Pointer(clipX+1, clipY+1, 1, 1)),
			"x=clip x+w-1": result(pixMap.Pointer(clipX+clipW-1, clipY+1, 1, 1)),
			"y=clip y+h-1": result(pixMap.Pointer(clipX+1, clipY+clipH-1, 1, 1)),
		}
		for name, res := range tests {
			t.Run(name, func(t *testing.T) {
				assert.True(t, res.ok)
				assert.NotZero(t, res.pointer)
			})
		}
	})

	t.Run("should calculate DeltaX", func(t *testing.T) {
		tests := map[string]struct {
			X        int
			Expected int
		}{
			"x=clipX": {
				X:        clipX,
				Expected: 0,
			},
			"x=clipX+1": {
				X:        clipX + 1,
				Expected: 0,
			},
			"x=clipX-1": {
				X:        clipX - 1,
				Expected: 1,
			},
			"x=clipX-2": {
				X:        clipX - 2,
				Expected: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				ptr, _ := pixMap.Pointer(test.X, clipY, 3, 1)
				assert.Equal(t, test.Expected, ptr.DeltaX)
			})
		}
	})

	t.Run("should calculate DeltaY", func(t *testing.T) {
		tests := map[string]struct {
			Y        int
			Expected int
		}{
			"y=clipY": {
				Y:        clipY,
				Expected: 0,
			},
			"y=clipY+1": {
				Y:        clipY + 1,
				Expected: 0,
			},
			"y=clipY-1": {
				Y:        clipY - 1,
				Expected: 1,
			},
			"y=clipY-2": {
				Y:        clipY - 2,
				Expected: 2,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				ptr, _ := pixMap.Pointer(clipX, test.Y, 3, 3)
				assert.Equal(t, test.Expected, ptr.DeltaY)
			})
		}
	})

	t.Run("should return a new slice from Pix slice", func(t *testing.T) {
		tests := map[string]struct {
			Ptr      pi.Pointer
			Expected []byte
		}{
			"entire pixmap": {
				Ptr:      pointer(pixMap.WithClip(0, 0, w, h).Pointer(0, 0, 1, 1)),
				Expected: pixMap.Pix(),
			},
			"1:": {
				Ptr:      pointer(pixMap.WithClip(0, 0, w, h).Pointer(1, 0, 1, 1)),
				Expected: pixMap.Pix()[1:],
			},
			"w:": {
				Ptr:      pointer(pixMap.WithClip(0, 0, w, h).Pointer(0, 1, 1, 1)),
				Expected: pixMap.Pix()[w:],
			},
			"pointer x<clipX": {
				Ptr:      pointer(pixMap.WithClip(1, 0, w, h).Pointer(0, 0, 2, 1)),
				Expected: pixMap.Pix()[1:],
			},
			"pointer y<clipY": {
				Ptr:      pointer(pixMap.WithClip(0, 1, w, h).Pointer(0, 0, 1, 2)),
				Expected: pixMap.Pix()[w:],
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, test.Expected, test.Ptr.Pix)
			})
		}
	})

	t.Run("should calculate RemainingPixels in line", func(t *testing.T) {
		tests := map[string]struct {
			Ptr      pi.Pointer
			Expected int
		}{
			"w=1": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY, 1, clipH)),
				Expected: 1,
			},
			"w>clipW": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY, clipW+1, clipH)),
				Expected: clipW,
			},
			"x<clipX,w=clipW+2": {
				Ptr:      pointer(pixMap.Pointer(clipX-1, clipY, clipW+2, clipH)),
				Expected: clipW,
			},
			"x<clipX,w<clipW": {
				Ptr:      pointer(pixMap.Pointer(clipX-1, clipY, clipW-1, clipH)),
				Expected: clipW - 1 - 1,
			},
			"x=clipX+clipW-1": {
				Ptr:      pointer(pixMap.Pointer(clipX+clipW-1, clipY, 2, clipH)),
				Expected: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, test.Expected, test.Ptr.RemainingPixels)
			})
		}
	})

	t.Run("should calculate RemainingLines", func(t *testing.T) {
		tests := map[string]struct {
			Ptr      pi.Pointer
			Expected int
		}{
			"h=1": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY, clipW, 1)),
				Expected: 1,
			},
			"h>clipH": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY, clipW, clipH+1)),
				Expected: clipH,
			},
			"y<clipY,h=clipH+2": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY-1, clipW, clipH+2)),
				Expected: clipH,
			},
			"y<clipY,h<clipH": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY-1, clipW, clipH-1)),
				Expected: clipH - 1 - 1,
			},
			"y=clipY+clipH-1": {
				Ptr:      pointer(pixMap.Pointer(clipX, clipY+clipH-1, clipW, 2)),
				Expected: 1,
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				assert.Equal(t, test.Expected, test.Ptr.RemainingLines)
			})
		}
	})
}

func TestPixMap_Foreach(t *testing.T) {
	t.Run("should not do anything when update function is nil", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		pixMap.Foreach(0, 0, 2, 2, nil)
	})

	t.Run("should run update line number of times", func(t *testing.T) {
		t.Run("inside clipping range", func(t *testing.T) {
			pixMap := pi.NewPixMap(1, 3)
			timesRun := 0
			pixMap.Foreach(0, 0, 1, 2, func(x, y int, dst []byte) {
				timesRun++
			})
			assert.Equal(t, 2, timesRun)
		})

		t.Run("outside clipping range", func(t *testing.T) {
			pixMap := pi.NewPixMap(1, 3)
			timesRun := 0
			pixMap.Foreach(0, 0, 1, 4, func(x, y int, dst []byte) {
				timesRun++
			})
			assert.Equal(t, 3, timesRun)
		})
	})

	t.Run("should pass trimmed lines", func(t *testing.T) {
		pix := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
		pixMap := pi.NewPixMapWithPixels(pix, 3)

		t.Run("x=0", func(t *testing.T) {
			var lines [][]byte
			pixMap.Foreach(0, 0, 2, 3, func(x, y int, dst []byte) {
				lines = append(lines, dst)
			})
			assert.Equal(t, [][]byte{{0, 1}, {3, 4}, {6, 7}}, lines)
		})

		t.Run("x=1", func(t *testing.T) {
			var lines [][]byte
			pixMap.Foreach(1, 0, 2, 3, func(x, y int, dst []byte) {
				lines = append(lines, dst)
			})
			assert.Equal(t, [][]byte{{1, 2}, {4, 5}, {7, 8}}, lines)
		})
	})

	t.Run("should pass x,y of each line", func(t *testing.T) {
		pix := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
		pixMap := pi.NewPixMapWithPixels(pix, 3)

		t.Run("inside clipping range", func(t *testing.T) {
			var coords []pi.Position
			pixMap.Foreach(1, 1, 2, 2, func(x, y int, dst []byte) {
				coords = append(coords, pi.Position{X: x, Y: y})
			})
			assert.Equal(t, []pi.Position{{X: 1, Y: 1}, {X: 1, Y: 2}}, coords)
		})

		t.Run("outside clipping range", func(t *testing.T) {
			var coords []pi.Position
			pixMap.Foreach(-1, -1, 3, 3, func(x, y int, dst []byte) {
				coords = append(coords, pi.Position{X: x, Y: y})
			})
			assert.Equal(t, []pi.Position{{X: 0, Y: 0}, {X: 0, Y: 1}}, coords)
		})
	})

	t.Run("should not run update when there are no pixels to update", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 2)
		executed := false
		update := func(x, y int, dst []byte) {
			executed = true
		}
		pixMap.Foreach(0, 0, 0, 1, update) // width = 0
		assert.False(t, executed)
	})

	t.Run("should update pixels", func(t *testing.T) {
		pixMap := pi.NewPixMap(2, 3)
		i := byte(1)
		pixMap.Foreach(0, 0, 2, 3, func(x, y int, dst []byte) {
			dst[0] = i
			dst[1] = i + 1
			i += 2
		})
		assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, pixMap.Pix())
	})
}

func TestPixMap_Copy(t *testing.T) {
	testPixMapCopy(t, pi.PixMap.Copy)
}

func testPixMapCopy(t *testing.T, merge func(pi.PixMap, int, int, int, int, pi.PixMap, int, int)) {
	t.Run("src bigger than dst", func(t *testing.T) {
		src := pi.NewPixMapWithPixels([]byte{
			1, 2, 3,
			4, 5, 6,
			7, 8, 9,
			10, 11, 12,
		}, 3)

		dstWidth := 2
		dstHeight := 3
		dstEmpty := []byte{0, 0, 0, 0, 0, 0}

		tests := map[string]struct {
			x, y, w, h int
			dstX, dstY int
			expected   []byte
		}{
			"w=0":      {w: 0, h: 1, expected: dstEmpty},
			"w=1":      {w: 1, h: 1, expected: []byte{1, 0, 0, 0, 0, 0}},
			"w=width":  {w: src.Width(), h: 1, expected: []byte{1, 2, 0, 0, 0, 0}},
			"w>width":  {w: src.Width() + 1, h: 1, expected: []byte{1, 2, 0, 0, 0, 0}},
			"h=0":      {w: 1, h: 0, expected: dstEmpty},
			"h=1":      {w: 1, h: 1, expected: []byte{1, 0, 0, 0, 0, 0}},
			"h=height": {w: 1, h: src.Height(), expected: []byte{1, 0, 4, 0, 7, 0}},
			"h>height": {w: 1, h: src.Height() + 1, expected: []byte{1, 0, 4, 0, 7, 0}},
			"x=-1,w=2": {x: -1, w: 2, h: 1, expected: []byte{0, 1, 0, 0, 0, 0}},
			"x=width":  {x: src.Width(), w: 1, h: 1, expected: dstEmpty},
			"y=-1,h=2": {y: -1, w: 1, h: 2, expected: []byte{0, 0, 1, 0, 0, 0}},
			"y=height": {y: src.Height(), w: 1, h: 1, expected: dstEmpty},
			"dstX=-1":  {w: 2, h: 1, dstX: -1, expected: []byte{2, 0, 0, 0, 0, 0}},
			"dstX=1":   {w: 2, h: 1, dstX: 1, expected: []byte{0, 1, 0, 0, 0, 0}},
			"dstY=-1":  {w: 1, h: 2, dstY: -1, expected: []byte{4, 0, 0, 0, 0, 0}},
			"dstY=1":   {w: 1, h: 2, dstY: 1, expected: []byte{0, 0, 1, 0, 4, 0}},

			"x=-dstWidth":  {x: -dstWidth, w: dstWidth + 1, h: 1, dstX: dstWidth - 1, dstY: dstHeight - 1, expected: dstEmpty},
			"y=-dstHeight": {y: -dstHeight, w: 1, h: dstHeight + 1, dstX: dstWidth - 1, dstY: dstHeight - 1, expected: dstEmpty},
			"dstX=-width":  {x: src.Width() - 1, y: src.Height() - 1, dstX: -src.Width(), w: src.Width() + 1, h: 1, expected: dstEmpty},
			"dstY=-height": {y: src.Height() - 1, dstY: -src.Height(), w: 1, h: src.Height() + 1, expected: dstEmpty},

			"src has more lines to draw than dst": {w: src.Width(), h: src.Height(), expected: []byte{1, 2, 4, 5, 7, 8}},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				dst := pi.NewPixMap(dstWidth, dstHeight)
				merge(src, test.x, test.y, test.w, test.h, dst, test.dstX, test.dstY)
				assert.Equal(t, test.expected, dst.Pix())
			})
		}
	})

	t.Run("src smaller than dst", func(t *testing.T) {
		src := pi.NewPixMapWithPixels([]byte{
			1, 2,
			4, 5,
			7, 8,
		}, 2)

		dstWidth := 3
		dstHeight := 5

		tests := map[string]struct {
			x, y, w, h int
			dstX, dstY int
			expected   []byte
		}{
			"src has less lines than dst":          {w: 2, h: dstHeight, expected: []byte{1, 2, 0, 4, 5, 0, 7, 8, 0, 0, 0, 0, 0, 0, 0}},
			"src has less pixels in line than dst": {w: dstWidth, h: 2, expected: []byte{1, 2, 0, 4, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
			"bottom-right corner of src":           {w: 2, h: 3, dstX: -1, dstY: -2, expected: []byte{8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}},
		}

		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				dst := pi.NewPixMap(dstWidth, dstHeight)
				merge(src, test.x, test.y, test.w, test.h, dst, test.dstX, test.dstY)
				assert.Equal(t, test.expected, dst.Pix())
			})
		}
	})
}

func TestPixMap_Merge(t *testing.T) {
	testPixMapCopy(t, func(src pi.PixMap, x int, y int, w int, h int, dst pi.PixMap, dstX int, dstY int) {
		src.Merge(x, y, w, h, dst, dstX, dstY, func(dst, src []byte) {
			copy(dst, src)
		})
	})
}

func TestPixMap_String(t *testing.T) {
	t.Run("should convert small pixmap to string", func(t *testing.T) {
		pixMap := pi.NewPixMap(1, 2)
		actual := pixMap.String()
		assert.Equal(t, "{width:1, height:2, clip:{X:0 Y:0 W:1 H:2}, pix:[0 0]}", actual)
	})

	t.Run("should convert big pixmap to string", func(t *testing.T) {
		pixMap := pi.NewPixMap(100, 100) // 10K bytes
		actual := pixMap.String()
		assert.True(t, len(actual) < 2500)
	})
}

type pointerResult struct {
	pointer pi.Pointer
	ok      bool
}

func result(pointer pi.Pointer, ok bool) pointerResult {
	return pointerResult{pointer: pointer, ok: ok}
}

func pointer(pointer pi.Pointer, _ bool) pi.Pointer {
	return pointer
}

func assertPixMapEqual(t *testing.T, pixMap pi.PixMap, file string) {
	expected := decodePNG(t, file).Pixels
	if !assert.Equal(t, expected, pixMap.Pix()) {
		screenshot, err := snap.Take()
		require.NoError(t, err)
		fmt.Println("Screenshot taken", screenshot)
	}
}

func assertScreenEqual(t *testing.T, file string) {
	assertPixMapEqual(t, pi.Scr(), file)
}

func assertEmptyPixMap(t *testing.T, pixMap pi.PixMap) {
	emptyPixMap := make([]byte, len(pixMap.Pix()))
	assert.Equal(t, emptyPixMap, pixMap.Pix())
}

func assertNotEmptyPixMap(t *testing.T, pixMap pi.PixMap) {
	emptyPixMap := make([]byte, len(pixMap.Pix()))
	assert.NotEqual(t, emptyPixMap, pixMap.Pix())
}
