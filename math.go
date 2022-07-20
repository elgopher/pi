package pi

import "math"

// Sin returns the sine of angle which is in the range of 0.0-1.0 measured clockwise.
//
// Sin returns an inverted result to suit screen space (where Y means "DOWN", as opposed
// to mathematical diagrams where Y typically means "UP"):
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
//	atan(0,-1)  // returns 0.25
func Atan2(dx, dy float64) float64 {
	v := math.Atan2(dx, dy)
	return math.Mod(0.75+v/(math.Pi*2), 1)
}
