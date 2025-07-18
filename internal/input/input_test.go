// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package input_test

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/internal/input"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestState_Duration(t *testing.T) {
	const btn = "btn"

	t.Run("should return 0 when input was never pressed", func(t *testing.T) {
		var i input.State[string]
		assert.Equal(t, 0, i.Duration(btn))
	})

	t.Run("should return 1 when input was pressed and released this frame", func(t *testing.T) {
		var i input.State[string]
		pi.Frame = 1
		i.SetDownFrame(btn, pi.Frame)
		i.SetUpFrame(btn, pi.Frame)
		assert.Equal(t, 1, i.Duration(btn))
	})

	t.Run("should return 0 when input was pressed previous frame and released this frame", func(t *testing.T) {
		var i input.State[string]
		pi.Frame = 0
		i.SetDownFrame(btn, pi.Frame)
		pi.Frame++
		i.SetUpFrame(btn, pi.Frame)
		assert.Equal(t, 0, i.Duration(btn))
	})

	t.Run("should return 2 when input was pressed previous frame but not released this frame", func(t *testing.T) {
		var i input.State[string]
		pi.Frame = 0
		i.SetDownFrame(btn, pi.Frame)
		pi.Frame++
		assert.Equal(t, 2, i.Duration(btn))
	})

	t.Run("should return 1 when input was pressed, released and pressed again this frame", func(t *testing.T) {
		var i input.State[string]
		pi.Frame = 0
		i.SetDownFrame(btn, pi.Frame)
		i.SetUpFrame(btn, pi.Frame)
		i.SetDownFrame(btn, pi.Frame)
		assert.Equal(t, 1, i.Duration(btn))
	})

	t.Run("should return 1 when input was pressed and released previous frame and pressed this frame", func(t *testing.T) {
		var i input.State[string]
		pi.Frame = 0
		i.SetDownFrame(btn, pi.Frame)
		i.SetUpFrame(btn, pi.Frame)
		pi.Frame++
		i.SetDownFrame(btn, pi.Frame)
		assert.Equal(t, 1, i.Duration(btn))
	})
}
