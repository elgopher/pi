// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"math"

	"github.com/elgopher/pi/vm"
)

// RectFill draws a filled rectangle between points x0,y0 and x1,y1 (inclusive).
//
// RectFill takes into account camera position, clipping region and draw palette.
func RectFill(x0, y0, x1, y1 int, color byte) {
	xmin, xmax := x0-vm.Camera.X, x1-vm.Camera.X
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	ymin, ymax := y0-vm.Camera.Y, y1-vm.Camera.Y
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	if xmin >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		return
	}

	if xmax < vm.ClippingRegion.X {
		return
	}

	if ymin >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		return
	}

	if ymax < vm.ClippingRegion.Y {
		return
	}

	if xmin < vm.ClippingRegion.X {
		xmin = vm.ClippingRegion.X
	}

	if xmax >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		xmax = vm.ClippingRegion.X + vm.ClippingRegion.W - 1
	}

	if ymin < vm.ClippingRegion.Y {
		ymin = vm.ClippingRegion.Y
	}

	if ymax >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		ymax = vm.ClippingRegion.Y + vm.ClippingRegion.H - 1
	}

	w := xmax - xmin + 1
	col := vm.DrawPalette[color]
	for i := 0; i < w; i++ {
		lineOfScreenWidth[i] = col
	}
	line := lineOfScreenWidth[:w]

	for y := ymin; y <= ymax; y++ {
		copy(vm.ScreenData[y*vm.ScreenWidth+xmin:], line)
	}
}

// Rect draws a rectangle between points x0,y0 and x1,y1 (inclusive).
//
// Rect takes into account camera position, clipping region and draw palette.
func Rect(x0, y0, x1, y1 int, color byte) {
	xmin, xmax := x0-vm.Camera.X, x1-vm.Camera.X
	if xmin > xmax {
		xmin, xmax = xmax, xmin
	}
	ymin, ymax := y0-vm.Camera.Y, y1-vm.Camera.Y
	if ymin > ymax {
		ymin, ymax = ymax, ymin
	}

	if xmin >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		return
	}

	if xmax < vm.ClippingRegion.X {
		return
	}

	if ymin >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		return
	}

	if ymax < vm.ClippingRegion.Y {
		return
	}

	drawLeftLine := true
	drawRightLine := true

	if xmin < vm.ClippingRegion.X {
		xmin = vm.ClippingRegion.X
		drawLeftLine = false
	}

	if xmax >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		xmax = vm.ClippingRegion.X + vm.ClippingRegion.W - 1
		drawRightLine = false
	}

	col := vm.DrawPalette[color]

	w := xmax - xmin + 1
	for i := 0; i < w; i++ {
		lineOfScreenWidth[i] = col
	}
	line := lineOfScreenWidth[:w]

	if ymin < vm.ClippingRegion.Y {
		ymin = vm.ClippingRegion.Y
	} else {
		copy(vm.ScreenData[ymin*vm.ScreenWidth+xmin:], line)
	}

	if ymax >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		ymax = vm.ClippingRegion.Y + vm.ClippingRegion.H - 1
	} else {
		copy(vm.ScreenData[ymax*vm.ScreenWidth+xmin:], line)
	}

	if drawLeftLine {
		for y := ymin; y <= ymax; y++ {
			vm.ScreenData[y*vm.ScreenWidth+xmin] = col
		}
	}

	if drawRightLine {
		for y := ymin; y <= ymax; y++ {
			vm.ScreenData[y*vm.ScreenWidth+xmax] = col
		}
	}
}

// Line draws a line between points x0,y0 and x1,y1 (inclusive).
//
// Line takes into account camera position, clipping region and draw palette.
func Line(x0, y0, x1, y1 int, color byte) {
	x0 -= vm.Camera.X
	x1 -= vm.Camera.X
	y0 -= vm.Camera.Y
	y1 -= vm.Camera.Y

	// Bresenham algorithm: https://www.youtube.com/watch?v=IDFB5CDpLDE
	run := float64(x1 - x0)
	if run == 0 {
		verticalLine(x0, y0, y1, color)
		return
	}

	rise := float64(y1 - y0)
	if rise == 0 {
		horizontalLine(y0, x0, x1, color)
		return
	}

	slope := rise / run

	adjust := 1
	if slope < 0 {
		adjust = -1
	}

	offset := 0.0    // performance could be better if offset was an integer instead
	threshold := 0.5 // performance could be better if threshold was an integer instead

	if slope >= -1 && slope <= 1 {
		delta := math.Abs(slope)
		y := y0
		if x1 < x0 {
			x0, x1 = x1, x0
			y = y1
		}

		for x := x0; x <= x1; x++ {
			pset(x, y, color)

			offset += delta
			if offset >= threshold {
				y += adjust
				threshold += 1
			}
		}
	} else {
		delta := math.Abs(run / rise)
		x := x0
		if y0 > y1 {
			y0, y1 = y1, y0
			x = x1
		}

		for y := y0; y <= y1; y++ {
			pset(x, y, color)

			offset += delta
			if offset >= threshold {
				x += adjust
				threshold += 1
			}
		}
	}
}

// verticalLine draws a vertical line between y0-y1 inclusive
func verticalLine(x, y0, y1 int, color byte) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	if x < vm.ClippingRegion.X {
		return
	}

	if x >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		return
	}

	if y0 < vm.ClippingRegion.Y {
		y0 = vm.ClippingRegion.Y
	}

	if y1 >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		y1 = vm.ClippingRegion.Y + vm.ClippingRegion.H - 1
	}

	for y := y0; y <= y1; y++ {
		vm.ScreenData[y*vm.ScreenWidth+x] = vm.DrawPalette[color]
	}
}

// horizontalLine draws a vertical line between x0-x1 inclusive
func horizontalLine(y, x0, x1 int, color byte) {
	if y < vm.ClippingRegion.Y {
		return
	}

	if y >= vm.ClippingRegion.Y+vm.ClippingRegion.H {
		return
	}

	if x0 > x1 {
		x0, x1 = x1, x0
	}

	if x0 < vm.ClippingRegion.X {
		x0 = vm.ClippingRegion.X
	}

	if x1 >= vm.ClippingRegion.X+vm.ClippingRegion.W {
		x1 = vm.ClippingRegion.X + vm.ClippingRegion.W - 1
	}

	offset := y * vm.ScreenWidth

	for x := x0; x <= x1; x++ {
		vm.ScreenData[offset+x] = vm.DrawPalette[color]
	}
}

// Circ draws a circle
//
// Circ takes into account camera position, clipping region and draw palette.
func Circ(centerX, centerY, radius int, color byte) {
	centerX = centerX - vm.Camera.X
	centerY = centerY - vm.Camera.Y
	// Code based on Frédéric Goset work: http://fredericgoset.ovh/mathematiques/courbes/en/bresenham_circle.html
	x := 0
	y := radius
	m := 5 - 4*radius

	for x <= y {
		pset(centerX+x, centerY+y, color)
		pset(centerX+x, centerY-y, color)
		pset(centerX-x, centerY+y, color)
		pset(centerX-x, centerY-y, color)
		pset(centerX+y, centerY+x, color)
		pset(centerX+y, centerY-x, color)
		pset(centerX-y, centerY+x, color)
		pset(centerX-y, centerY-x, color)

		if m > 0 {
			y--
			m -= 8 * y
		}

		x++

		m += 8*x + 4
	}
}

// CircFill draws a filled circle
//
// CircFill takes into account camera position, clipping region and draw palette.
func CircFill(centerX, centerY, radius int, color byte) {
	// Code based on Frédéric Goset work: http://fredericgoset.ovh/mathematiques/courbes/en/filled_circle.html
	x := 0
	y := radius
	m := 5 - 4*radius

	for x <= y {
		RectFill(centerX-y, centerY-x, centerX+y, centerY-x, color)
		RectFill(centerX-y, centerY+x, centerX+y, centerY+x, color)

		if m > 0 {
			RectFill(centerX-x, centerY-y, centerX+x, centerY-y, color)
			RectFill(centerX-x, centerY+y, centerX+x, centerY+y, color)
			y--
			m -= 8 * y
		}

		x++

		m += 8*x + 4
	}
}
