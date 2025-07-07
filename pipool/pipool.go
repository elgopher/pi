// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pipool provides an extremely simple, non-thread-safe pool
// that can be used to reduce heap memory allocations.
package pipool

// Pool is a very simple, non-thread-safe object pool.
//
// It can only be used from a single goroutine. Since Pi runs on
// a single goroutine, it can safely be used in pi.Update and pi.Draw.
type Pool[T any] struct {
	objects []*T // LIFO
}

// Get returns an object from the pool.
//
// If the pool is empty, it creates a new zero-value object
// and returns a pointer to it.
func (p *Pool[T]) Get() *T {
	n := len(p.objects)
	if n == 0 {
		var t T
		return &t
	}
	last := p.objects[n-1]
	p.objects = p.objects[:n-1] // decrease len, cap will remain the same
	return last
}

// Put returns an object to the pool.
func (p *Pool[T]) Put(obj *T) {
	p.objects = append(p.objects, obj)
}
