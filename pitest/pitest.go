// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pitest provides helper functions for writing unit tests.
package pitest

import (
	"testing"
	"unsafe"

	"github.com/elgopher/pi"
)

// AssertSurfaceEqual asserts that two surfaces are equal:
// they have the same size and identical data.
func AssertSurfaceEqual[T comparable](t *testing.T, expected, actual pi.Surface[T]) (equal bool) {
	t.Helper()

	equal = true

	if expected.W() != actual.W() {
		equal = false
		t.Errorf("expected width: %v, actual width: %v", expected.W(), actual.W())
	}
	if expected.H() != actual.H() {
		equal = false
		t.Errorf("expected height: %v, actual height: %v", expected.H(), actual.H())
	}

	var hasDifferentData bool

	for y := 0; y < expected.H(); y++ {
		for x := 0; x < expected.W(); x++ {
			e := expected.Get(x, y)
			a := actual.Get(x, y)
			if e != a {
				if !hasDifferentData {
					t.Errorf("\nexpected:\n%s\nactual:\n%s", expected, actual)
				}
				t.Errorf("expected at (%d,%d): %v, actual: %v", x, y, e, a)
				hasDifferentData = true
				equal = false
			}
		}
	}

	return
}

// AssertSpriteEqual asserts that two sprites have equal fields
// and reference the same source
func AssertSpriteEqual(t *testing.T, expected, actual pi.Sprite) {
	t.Helper()

	expectedSource := expected.Source
	actualSource := actual.Source
	expectedArrayPointer := unsafe.SliceData(expectedSource.Data())
	actualArrayPointer := unsafe.SliceData(actualSource.Data())
	sameArrayPointer := actualArrayPointer == expectedArrayPointer

	if !sameArrayPointer {
		t.Errorf("sprites are using different sources, "+
			"pointers to Canvas data are not the same, expected: %+v, actual: %+v",
			expectedArrayPointer, actualArrayPointer)
	}

	if expected.Area != actual.Area {
		t.Errorf("expected sprite area: %+v, actual area: %+v", expected.Area, actual.Area)
	}
	if expected.FlipX != actual.FlipX {
		t.Errorf("expected sprite flipX: %v, actual flipX: %v", expected.FlipX, actual.FlipX)
	}
	if expected.FlipY != actual.FlipY {
		t.Errorf("expected sprite flipY: %v, actual flipY: %v", expected.FlipY, actual.FlipY)
	}
}
