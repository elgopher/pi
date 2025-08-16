// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piroutine_test

import (
	"testing"

	"github.com/elgopher/pi/piroutine"
	"github.com/stretchr/testify/assert"
)

func TestPiroutine(t *testing.T) {
	r := piroutine.New(
		piroutine.Wait(2),
		piroutine.Printf("abc"),
		piroutine.Wait(5),
	)
	r.SetTracing(true)
	for r.Resume() {
	}
}

func TestRoutine_Resume(t *testing.T) {
	t.Run("empty routine", func(t *testing.T) {
		routine := piroutine.New()
		assert.False(t, routine.Resume())
	})

	t.Run("single step", func(t *testing.T) {
		var stepExecuted bool
		step := func() bool {
			stepExecuted = true
			return true // finish immediately
		}
		routine := piroutine.New(step)
		assert.False(t, routine.Resume()) // finished
		assert.True(t, stepExecuted)
	})

	t.Run("wait step", func(t *testing.T) {
		routine := piroutine.New(piroutine.Wait(1))
		assert.True(t, routine.Resume())  // not yet finished
		assert.False(t, routine.Resume()) // finished
		assert.False(t, routine.Resume()) // nothing changed
	})
}
