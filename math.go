// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import "math"

// Sin returns the sine of angle which is in the range of 0.0-1.0 measured clockwise.
//
// Sin returns an inverted result to suit screen space (where Y means "DOWN", as opposed
// to mathematical diagrams where Y typically means "UP"):
//
//	sin(0.25) // returns -1
//
// If you want to use conventional radian-based function without the y inversion, use math.Sin.
func Sin(angle float64) float64 {
	rad := angle * 2 * math.Pi
	return -math.Sin(rad)
}

// Cos returns the cosine of angle which is in the range of 0.0-1.0 measured clockwise.
//
// If you want to use conventional radian-based function use math.Cos.
func Cos(angle float64) float64 {
	rad := angle * 2 * math.Pi
	return math.Cos(rad)
}

// Atan2 converts DX, DY into an angle from 0..1
//
// Similar to Cos and Sin, angle is taken to run anticlockwise in screenspace. For example:
//
//	atan(0,-1)  // returns 0.25
func Atan2(dx, dy float64) float64 {
	v := math.Atan2(dx, dy)
	return math.Mod(0.75+v/(math.Pi*2), 1)
}

// Int is a generic type for all integer types
type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// MinInt returns minimum of two integer numbers.
func MinInt[T Int](x, y T) T {
	if x < y {
		return x
	}

	return y
}

// MaxInt returns maximum of two integer numbers.
func MaxInt[T Int](x, y T) T {
	if x > y {
		return x
	}

	return y
}

// MidInt returns the middle of three integer numbers. Very useful for clamping.
func MidInt[T Int](x, y, z T) T {
	if x > y {
		x, y = y, x
	}

	if y > z {
		y = z
	}

	if x > y {
		y = x
	}

	return y
}

// Mid returns the middle of three float64 numbers. Very useful for clamping.
// NaNs are always put at the beginning (are the smallest ones).
func Mid(x, y, z float64) float64 {
	if x > y || math.IsNaN(y) {
		x, y = y, x
	}

	if y > z || math.IsNaN(z) {
		y = z
	}

	if x > y || math.IsNaN(y) {
		y = x
	}

	return y
}
