// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pimath

import "math"

// Lerp computes the linear interpolation between a and b.
//
// The parameter t should be in the range 0 to 1.
func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// Number describes any numeric type in Go.
//
// Includes signed integers, unsigned integers, and floating-point types.
type Number interface {
	~int | ~float64 |
		~int8 | ~int16 | ~int32 | ~int64 |
		~float32 |
		~uint | ~byte | ~uint16 | ~uint32 | ~uint64
}

// Clamp limits the value x to the range [min, max].
func Clamp[T Number](x, min, max T) T {
	if x < min {
		return min
	} else if x > max {
		return max
	}
	return x
}

// Distance returns the distance between (x1, y1) and (x2, y2).
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}
