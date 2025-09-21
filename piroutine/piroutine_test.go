// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piroutine_test

import (
	"testing"

	"github.com/elgopher/pi/piloop"
	"github.com/elgopher/pi/piroutine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCall(t *testing.T) {
	t.Run("should call callback", func(t *testing.T) {
		var executed = false
		step := piroutine.Call(func() {
			executed = true
		})
		// when
		result := step()
		// then
		assert.True(t, executed)
		assert.True(t, result)
	})
}

func TestSlowDown(t *testing.T) {
	t.Run("should wait n updates before running callback", func(t *testing.T) {
		executionCount := 0
		step := piroutine.SlowDown(2, func() bool {
			executionCount++
			return true
		})
		assert.False(t, step())
		assert.Equal(t, 0, executionCount)
		assert.False(t, step())
		assert.Equal(t, 0, executionCount)
		assert.True(t, step())
		assert.Equal(t, 1, executionCount)
	})

	t.Run("should immediately run callback", func(t *testing.T) {
		executionCount := 0
		step := piroutine.SlowDown(0, func() bool {
			executionCount++
			return true
		})
		assert.True(t, step())
		assert.Equal(t, 1, executionCount)
	})

	t.Run("should wait another n updates after callback returned false", func(t *testing.T) {
		executionCount := 0
		step := piroutine.SlowDown(3, func() bool {
			executionCount++
			return executionCount%2 == 0
		})
		for range 3 {
			assert.False(t, step()) // wait
		}
		assert.False(t, step()) // callback returns false
		for range 3 {
			assert.False(t, step()) // wait
		}
		assert.True(t, step()) // callback returns true this time
		assert.Equal(t, 2, executionCount)
	})
}

func TestRoutine_ScheduleOn(t *testing.T) {
	t.Run("should run callback on event", func(t *testing.T) {
		executionCount := 0
		step := func() bool {
			executionCount++
			return true
		}
		routine := piroutine.New(step)
		// when
		handler := routine.ScheduleOn(piloop.EventDraw)
		// then
		require.True(t, piloop.Target().IsSubscribed(handler))
		piloop.Target().Publish(piloop.EventDraw) // runs callback and unsubscribes handlers
		assert.Equal(t, 1, executionCount)
		assert.False(t, piloop.Target().IsSubscribed(handler))
	})
}
