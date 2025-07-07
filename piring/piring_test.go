// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piring_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/piring"
)

func TestBuffer(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		buffer := piring.NewBuffer[int](4)
		assert.Equal(t, 0, buffer.Len())
		assert.Zero(t, *buffer.PointerTo(0))
	})

	t.Run("no circle", func(t *testing.T) {
		buffer := piring.NewBuffer[int](3)
		*buffer.NextWritePointer() = 0
		*buffer.NextWritePointer() = 1
		*buffer.NextWritePointer() = 2

		assert.Equal(t, 3, buffer.Len())

		assertHas(t, buffer, 0, 0)
		assertHas(t, buffer, 1, 1)
		assertHas(t, buffer, 2, 2)
	})

	t.Run("circle", func(t *testing.T) {
		buffer := piring.NewBuffer[int](2)
		*buffer.NextWritePointer() = 0
		*buffer.NextWritePointer() = 1
		*buffer.NextWritePointer() = 2

		assert.Equal(t, 2, buffer.Len())

		assertHas(t, buffer, 0, 1)
		assertHas(t, buffer, 1, 2)
	})

	t.Run("two circles", func(t *testing.T) {
		buffer := piring.NewBuffer[int](2)
		*buffer.NextWritePointer() = 0
		*buffer.NextWritePointer() = 1
		*buffer.NextWritePointer() = 2
		*buffer.NextWritePointer() = 3
		*buffer.NextWritePointer() = 4

		assert.Equal(t, 2, buffer.Len())

		assertHas(t, buffer, 0, 3)
		assertHas(t, buffer, 1, 4)
	})

	t.Run("out of range for filled buffer", func(t *testing.T) {
		buffer := piring.NewBuffer[int](3)
		for i := 0; i < buffer.Cap(); i++ {
			*buffer.NextWritePointer() = i
		}

		assertHas(t, buffer, -1, *buffer.PointerTo(2))
		assertHas(t, buffer, 3, *buffer.PointerTo(0))
	})

	t.Run("out of range for half-empty buffer", func(t *testing.T) {
		buffer := piring.NewBuffer[int](4)
		*buffer.NextWritePointer() = 10
		*buffer.NextWritePointer() = 20

		assertHas(t, buffer, -1, *buffer.PointerTo(3))
		assertHas(t, buffer, -3, *buffer.PointerTo(1))
		assertHas(t, buffer, 4, *buffer.PointerTo(0))
	})

	t.Run("out of range for overfilled buffer", func(t *testing.T) {
		buffer := piring.NewBuffer[int](2)
		*buffer.NextWritePointer() = 0
		*buffer.NextWritePointer() = 1
		*buffer.NextWritePointer() = 2

		assertHas(t, buffer, -1, *buffer.PointerTo(1))
		assertHas(t, buffer, 2, *buffer.PointerTo(0))
	})

	t.Run("out of range, far away", func(t *testing.T) {
		buffer := piring.NewBuffer[int](2)
		*buffer.NextWritePointer() = 0
		*buffer.NextWritePointer() = 1

		assertHas(t, buffer, 5, *buffer.PointerTo(1))
		assertHas(t, buffer, -5, *buffer.PointerTo(1))
	})
}

func assertHas[T any](t *testing.T, buffer *piring.Buffer[T], index int, value T) {
	t.Helper()
	assert.Equal(t, value, *buffer.PointerTo(index))
}
