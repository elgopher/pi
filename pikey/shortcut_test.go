// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pikey_test

import (
	"testing"

	"github.com/elgopher/pi/pikey"
	"github.com/elgopher/pi/piloop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	aDownEvent    = pikey.Event{Type: pikey.EventDown, Key: pikey.A}
	aUpEvent      = pikey.Event{Type: pikey.EventUp, Key: pikey.A}
	ctrlDownEvent = pikey.Event{Type: pikey.EventDown, Key: pikey.Ctrl}
)

func TestRegisterShortcut(t *testing.T) {
	t.Run("single key", func(t *testing.T) {
		executionTimes := 0
		shortcut := pikey.RegisterShortcut(func() {
			executionTimes++
		}, pikey.A)
		require.NotNil(t, shortcut)

		pikey.Target().Publish(aDownEvent)
		piloop.Target().Publish(piloop.EventLateUpdate)

		t.Run("should execute callback when key is pressed", func(t *testing.T) {
			assert.Equal(t, 1, executionTimes)
		})

		executionTimes = 0

		piloop.Target().Publish(piloop.EventLateUpdate)

		t.Run("should not execute callback again on next frame when key is still pressed", func(t *testing.T) {
			assert.Equal(t, 0, executionTimes)
		})

		shortcut.Unregister()

		pikey.Target().Publish(aDownEvent)
		piloop.Target().Publish(piloop.EventLateUpdate)

		t.Run("should not execute callback after shortcut is unregistered", func(t *testing.T) {
			assert.Equal(t, 0, executionTimes)
		})
	})

	t.Run("multiple keys", func(t *testing.T) {
		executionTimes := 0
		shortcut := pikey.RegisterShortcut(func() {
			executionTimes++
		}, pikey.A, pikey.Ctrl)
		require.NotNil(t, shortcut)
		// when
		pikey.Target().Publish(aDownEvent)
		pikey.Target().Publish(ctrlDownEvent)
		piloop.Target().Publish(piloop.EventLateUpdate)
		// then
		assert.Equal(t, 1, executionTimes)
	})

	t.Run("should not run callback when all keys are not pressed simultaneously, but was down before", func(t *testing.T) {
		executionTimes := 0
		shortcut := pikey.RegisterShortcut(func() {
			executionTimes++
		}, pikey.A, pikey.Ctrl)
		require.NotNil(t, shortcut)

		pikey.Target().Publish(aDownEvent)
		piloop.Target().Publish(piloop.EventLateUpdate)
		require.Equal(t, 0, executionTimes)
		// when
		pikey.Target().Publish(ctrlDownEvent)
		pikey.Target().Publish(aUpEvent) // "A" no longer down
		piloop.Target().Publish(piloop.EventLateUpdate)
		// then
		assert.Equal(t, 0, executionTimes)
	})
}
