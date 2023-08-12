// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "math"

// RectFill draws a filled rectangle on screen between points x0,y0 and x1,y1 (inclusive).
//
// RectFill takes into account camera position, clipping region and draw palette.
func RectFill(x0, y0, x1, y1 int, color byte) {
	col := DrawPalette[color]

	x0 -= Camera.X
	x1 -= Camera.X
	y0 -= Camera.Y
	y1 -= Camera.Y

	screen.RectFill(x0, y0, x1, y1, col)
}

// RectFill draws a filled rectangle between points x0,y0 and x1,y1 (inclusive).
func (p PixMap) RectFill(x0 int, y0 int, x1 int, y1 int, col byte) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	ptr, ok := p.Pointer(x0, y0, x1-x0+1, y1-y0+1)
	if !ok {
		return
	}

	line := p.lineOfColor(col, ptr.RemainingPixels)

	copy(ptr.Pix, line)
	for y := 1; y < ptr.RemainingLines; y++ {
		ptr.Pix = ptr.Pix[p.width:]
		copy(ptr.Pix, line)
	}
}

// Rect draws a rectangle on screen between points x0,y0 and x1,y1 (inclusive).
//
// Rect takes into account camera position, clipping region and draw palette.
func Rect(x0, y0, x1, y1 int, color byte) {
	color = DrawPalette[color]

	x0, x1 = x0-Camera.X, x1-Camera.X
	y0, y1 = y0-Camera.Y, y1-Camera.Y

	screen.Rect(x0, y0, x1, y1, color)
}

// Rect draws a rectangle between points x0,y0 and x1,y1 (inclusive).
func (p PixMap) Rect(x0, y0, x1, y1 int, col byte) {
	p.horizontalLine(y0, x0, x1, col)
	p.horizontalLine(y1, x0, x1, col)
	p.verticalLine(x0, y0, y1, col)
	p.verticalLine(x1, y0, y1, col)
}

// Line draws a line on screen between points x0,y0 and x1,y1 (inclusive).
//
// Line takes into account camera position, clipping region and draw palette.
func Line(x0, y0, x1, y1 int, color byte) {
	color = DrawPalette[color]

	x0 -= Camera.X
	x1 -= Camera.X
	y0 -= Camera.Y
	y1 -= Camera.Y

	screen.Line(x0, y0, x1, y1, color)
}

// Line draws a line between points x0,y0 and x1,y1 (inclusive).
func (p PixMap) Line(x0, y0, x1, y1 int, color byte) {
	// Bresenham algorithm: https://www.youtube.com/watch?v=IDFB5CDpLDE
	run := float64(x1 - x0)
	if run == 0 {
		p.verticalLine(x0, y0, y1, color)
		return
	}

	rise := float64(y1 - y0)
	if rise == 0 {
		p.horizontalLine(y0, x0, x1, color)
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
			p.Set(x, y, color)

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
			p.Set(x, y, color)

			offset += delta
			if offset >= threshold {
				x += adjust
				threshold += 1
			}
		}
	}
}

// verticalLine draws a vertical line between y0-y1 inclusive
func (p PixMap) verticalLine(x, y0, y1 int, color byte) {
	if y0 > y1 {
		y0, y1 = y1, y0
	}

	ptr, ok := p.Pointer(x, y0, 1, y1-y0+1)
	if !ok {
		return
	}

	index := 0
	for i := 0; i < ptr.RemainingLines; i++ {
		ptr.Pix[index] = color
		index += p.width
	}
}

// horizontalLine draws a vertical line between x0-x1 inclusive
func (p PixMap) horizontalLine(y, x0, x1 int, color byte) {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	ptr, ok := p.Pointer(x0, y, x1-x0+1, 1)
	if !ok {
		return
	}

	pix := ptr.Pix[:ptr.RemainingPixels]
	for i := 0; i < len(pix); i++ {
		pix[i] = color
	}
}

// Circ draws a circle on screen.
//
// Circ takes into account camera position, clipping region and draw palette.
func Circ(centerX, centerY, radius int, color byte) {
	color = DrawPalette[color]

	centerX = centerX - Camera.X
	centerY = centerY - Camera.Y

	screen.Circ(centerX, centerY, radius, color)
}

// Circ draws a circle.
func (p PixMap) Circ(centerX int, centerY int, radius int, color byte) {
	// Code based on Frédéric Goset work: http://fredericgoset.ovh/mathematiques/courbes/en/bresenham_circle.html
	x := 0
	y := radius
	m := 5 - 4*radius

	for x <= y {
		p.Set(centerX+x, centerY+y, color)
		p.Set(centerX+x, centerY-y, color)
		p.Set(centerX-x, centerY+y, color)
		p.Set(centerX-x, centerY-y, color)
		p.Set(centerX+y, centerY+x, color)
		p.Set(centerX+y, centerY-x, color)
		p.Set(centerX-y, centerY+x, color)
		p.Set(centerX-y, centerY-x, color)

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
	color = DrawPalette[color]

	centerX = centerX - Camera.X
	centerY = centerY - Camera.Y

	// Code based on Frédéric Goset work: http://fredericgoset.ovh/mathematiques/courbes/en/filled_circle.html
	screen.CircFill(centerX, centerY, radius, color)
}

func (p PixMap) CircFill(centerX int, centerY int, radius int, color byte) {
	x := 0
	y := radius
	m := 5 - 4*radius

	for x <= y {
		p.RectFill(centerX-y, centerY-x, centerX+y, centerY-x, color)
		p.RectFill(centerX-y, centerY+x, centerX+y, centerY+x, color)

		if m > 0 {
			p.RectFill(centerX-x, centerY-y, centerX+x, centerY-y, color)
			p.RectFill(centerX-x, centerY+y, centerX+x, centerY+y, color)
			y--
			m -= 8 * y
		}

		x++

		m += 8*x + 4
	}
}
