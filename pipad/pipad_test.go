// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pipad_test

import (
	"testing"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pipad"
	"github.com/stretchr/testify/assert"
)

var (
	player0connected    = pipad.EventConnection{Type: pipad.EventConnect, Player: 0}
	player1connected    = pipad.EventConnection{Type: pipad.EventConnect, Player: 1}
	player0disconnected = pipad.EventConnection{Type: pipad.EventDisconnect, Player: 0}
	player1disconnected = pipad.EventConnection{Type: pipad.EventDisconnect, Player: 1}
)

func TestPlayerCount(t *testing.T) {
	assert.Equal(t, 0, pipad.PlayerCount())
	pipad.ConnectionTarget().Publish(player0connected)
	pipad.ConnectionTarget().Publish(player1connected)
	assert.Equal(t, 2, pipad.PlayerCount())
	pipad.ConnectionTarget().Publish(player0disconnected)
	assert.Equal(t, 1, pipad.PlayerCount())
	pipad.ConnectionTarget().Publish(player1disconnected)
}

func TestDuration(t *testing.T) {
	{
		pipad.ConnectionTarget().Publish(player0connected)
		pipad.ButtonTarget().Publish(
			pipad.EventButton{Type: pipad.EventDown, Button: pipad.A, Player: 0},
		)

		t.Run("should return duration when button was pressed", func(t *testing.T) {
			duration := pipad.Duration(pipad.A)
			assert.Equal(t, 1, duration)

			playerDuration := pipad.PlayerDuration(pipad.A, 0)
			assert.Equal(t, 1, playerDuration)
		})

		pi.Frame++

		t.Run("should take into account how many frames passed", func(t *testing.T) {
			assert.Equal(t, 2, pipad.Duration(pipad.A))
			assert.Equal(t, 2, pipad.PlayerDuration(pipad.A, 0))
		})

		pipad.ConnectionTarget().Publish(player0disconnected)
	}

	t.Run("should return 0 after controller was disconnected", func(t *testing.T) {
		assert.Equal(t, 0, pipad.Duration(pipad.A))
		assert.Equal(t, 0, pipad.PlayerDuration(pipad.A, 0))
	})

	t.Run("should return the longest duration when two controllers are pressed simultaneously", func(t *testing.T) {
		pipad.ConnectionTarget().Publish(player0connected)
		defer pipad.ConnectionTarget().Publish(player0disconnected)
		pipad.ConnectionTarget().Publish(player1connected)
		defer pipad.ConnectionTarget().Publish(player1disconnected)

		pipad.ButtonTarget().Publish(
			pipad.EventButton{Type: pipad.EventDown, Button: pipad.A, Player: 0},
		)
		pi.Frame++
		pipad.ButtonTarget().Publish(
			pipad.EventButton{Type: pipad.EventUp, Button: pipad.A, Player: 1},
		)
		assert.Equal(t, 2, pipad.Duration(pipad.A))
		assert.Equal(t, 2, pipad.PlayerDuration(pipad.A, 0))
	})

	t.Run("should return 0 when player was never connected", func(t *testing.T) {
		assert.Equal(t, 0, pipad.PlayerDuration(pipad.A, 2))
	})
}
