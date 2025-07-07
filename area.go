// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Area specifies rectangular boundaries on a [Surface].
//
// It is used by [Sprite], [SetClip], [Surface.CloneArea], [Surface.LinesIterator], and others.
type Area[T Number] struct {
	X, Y, W, H T
}

// IntArea is an Area with integer coordinates.
type IntArea = Area[int]

// Size returns the number of values in the area.
func (a Area[T]) Size() T {
	return a.W * a.H
}

// WithX returns a new Area with the X value updated.
func (a Area[T]) WithX(x T) Area[T] {
	a.X = x
	return a
}

// WithY returns a new Area with the Y value updated.
func (a Area[T]) WithY(y T) Area[T] {
	a.Y = y
	return a
}

// WithW returns a new Area with the W value updated.
func (a Area[T]) WithW(w T) Area[T] {
	a.W = w
	return a
}

// WithH returns a new Area with the H value updated.
func (a Area[T]) WithH(h T) Area[T] {
	a.H = h
	return a
}

// MovedBy returns a new Area moved by dx and dy.
func (a Area[T]) MovedBy(dx, dy T) Area[T] {
	a.X += dx
	a.Y += dy
	return a
}

// ClippedBy clips the area so it does not extend beyond the given clip region.
//
// In addition to the new area, it also returns how much it was shifted.
// The output parameters dx and dy are always positive.
func (a Area[T]) ClippedBy(clip Area[T]) (_ Area[T], dx, dy T) {
	if a.X < clip.X {
		dx = clip.X - a.X
		a.W -= dx
		a.X = clip.X
	}
	if a.Y < clip.Y {
		dy = clip.Y - a.Y
		a.H -= dy
		a.Y = clip.Y
	}
	if a.X+a.W > clip.X+clip.W {
		a.W = clip.X + clip.W - a.X
	}
	if a.W < 0 {
		a.W = 0
	}
	if a.Y+a.H > clip.Y+clip.H {
		a.H = clip.Y + clip.H - a.Y
	}
	if a.H < 0 {
		a.H = 0
	}
	return a, dx, dy
}

// Contains reports whether the point (x, y) is inside the Area.
func (a Area[T]) Contains(x, y T) bool {
	return x >= a.X && x < a.X+a.W && y >= a.Y && y < a.Y+a.H
}
