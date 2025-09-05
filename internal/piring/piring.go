// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piring provides a ring buffer implementation.
package piring

func NewBuffer[E any](size int) *Buffer[E] {
	return &Buffer[E]{data: make([]E, size)}
}

// Buffer is a ring buffer (circular queue) that allows reusing
// elements stored in the buffer.
//
// It supports writing new elements. When full,
// the oldest element is overwritten with the new one.
// All elements can be accessed at any time by index.
type Buffer[E any] struct {
	data  []E
	write int
	len   int
	start int
}

// Len returns the number of elements currently stored in the buffer.
func (b *Buffer[E]) Len() int {
	return b.len
}

// Cap returns the maximum number of elements the buffer can hold.
func (b *Buffer[E]) Cap() int {
	return len(b.data)
}

// PointerTo returns a pointer to the element at the given index.
//
// PointerTo(0) returns the oldest element, while PointerTo(Len()-1) returns the newest.
// The index can be out of range — negative or greater than Len or Cap.
// In such cases, the index is wrapped until it fits within the buffer
// and corresponds to a valid element.
//
// PointerTo always returns a pointer to an element in the buffer.
func (b *Buffer[E]) PointerTo(index int) *E {
	idx := index + b.start

	// wrap the index so that both positive and negative values fall
	// into the [0, Cap()-1] range
	if capacity := len(b.data); capacity > 0 {
		idx = ((idx % capacity) + capacity) % capacity
	}

	return &b.data[idx]
}

// NextWritePointer returns a pointer to the next write location.
//
// If there is no space left, the oldest element is returned.
// This allows reusing the oldest element and can help avoid allocations.
func (b *Buffer[E]) NextWritePointer() *E {
	if len(b.data) <= b.len {
		b.start++
		b.len--
		if b.start == len(b.data) {
			b.start = 0
		}
	}
	if b.write >= len(b.data) {
		b.write = 0
	}
	e := &b.data[b.write]
	b.write++
	b.len += 1
	return e
}

// Reset clears the buffer logically for the user,
// but does not actually remove the underlying stored elements.
func (b *Buffer[E]) Reset() {
	b.write = 0
	b.start = 0
	b.len = 0
}
