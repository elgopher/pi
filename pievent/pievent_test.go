// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pievent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/pievent"
)

func TestEventHandler_Unsubscribe(t *testing.T) {
	t.Run("should not unsubscribe other handler", func(t *testing.T) {
		target := pievent.NewTarget[string]()
		first := target.SubscribeAll(func(s string, _ pievent.Handler) {})
		target.Unsubscribe(first)
		eventReceived := false
		target.SubscribeAll(func(string, pievent.Handler) {
			eventReceived = true
		})
		target.Unsubscribe(first) // remove first handler again
		target.Publish("test")
		assert.True(t, eventReceived)
	})
}
