// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pikey_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pikey"
)

func TestDuration(t *testing.T) {
	t.Run("should return 0 by default", func(t *testing.T) {
		assert.Equal(t, 0, pikey.Duration(pikey.A))
	})

	pikey.Target().Publish(aDownEvent)

	t.Run("should return 1 when key was down in the current frame", func(t *testing.T) {
		assert.Equal(t, 1, pikey.Duration(pikey.A))
	})

	pi.Frame++

	t.Run("should return 2 when key has been down since the previous frame", func(t *testing.T) {
		assert.Equal(t, 2, pikey.Duration(pikey.A))
	})

	pikey.Target().Publish(aUpEvent)

	t.Run("should return 0 when key is up", func(t *testing.T) {
		assert.Equal(t, 0, pikey.Duration(pikey.A))
	})
}
