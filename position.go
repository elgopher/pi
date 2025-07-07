// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Position represents a 2D integer coordinate.
//
// It stores X and Y values in a 2D grid.
type Position struct{ X, Y int }

// Add returns a new Position after adding coordinates with other
func (p Position) Add(other Position) Position {
	return Position{p.X + other.X, p.Y + other.Y}
}

// Subtract returns a new Position after subtracting coordinates with other
func (p Position) Subtract(other Position) Position {
	return Position{p.X - other.X, p.Y - other.Y}
}

// WithX returns a new Position with updated X.
func (p Position) WithX(x int) Position {
	p.X = x
	return p
}

// WithY returns a new Position with updated Y.
func (p Position) WithY(y int) Position {
	p.Y = y
	return p
}
