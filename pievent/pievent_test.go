// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pievent_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elgopher/pi/pievent"
)

const event = "event"

func TestTarget_Publish(t *testing.T) {
	t.Run("should publish event to subscriber synchronously", func(t *testing.T) {
		target := pievent.NewTarget[string]()
		var subscriber Subscriber
		target.Subscribe(event, subscriber.Handler)
		// when
		target.Publish(event)
		// then
		subscriber.AssertEventReceived(t, event)
	})

	t.Run("should publish event to all subscribers", func(t *testing.T) {
		var subscriber1, subscriber2 Subscriber

		target := pievent.NewTarget[string]()
		target.Subscribe(event, subscriber1.Handler)
		target.Subscribe(event, subscriber2.Handler)
		// when
		target.Publish(event)
		// then
		subscriber1.AssertEventReceived(t, event)
		subscriber2.AssertEventReceived(t, event)
	})

	t.Run("should not publish event to subscriber listening to a different event", func(t *testing.T) {
		var subscriber1, subscriber2 Subscriber

		target := pievent.NewTarget[string]()
		target.Subscribe(event, subscriber1.Handler)
		target.Subscribe("other", subscriber2.Handler)
		// when
		target.Publish(event)
		// then
		subscriber2.AssertNoEventReceived(t)
	})

	t.Run("should publish event to subscriber listening for all events", func(t *testing.T) {
		t.Run("subscribe with zero-value event", func(t *testing.T) {
			var subscriber Subscriber

			target := pievent.NewTarget[string]()
			target.Subscribe("", subscriber.Handler) // zero-value string register subscriber for all events
			// when
			target.Publish(event)
			// then
			subscriber.AssertEventReceived(t, event)
		})

		t.Run("subscribe all", func(t *testing.T) {
			var subscriber Subscriber

			target := pievent.NewTarget[string]()
			target.SubscribeAll(subscriber.Handler)
			// when
			target.Publish(event)
			// then
			subscriber.AssertEventReceived(t, event)
		})
	})
}

func TestTarget_IsSubscribed(t *testing.T) {
	t.Run("should return true if subscriber is subscribed", func(t *testing.T) {
		target := pievent.NewTarget[string]()
		var subscriber Subscriber
		handler := target.Subscribe("", subscriber.Handler)
		assert.True(t, target.IsSubscribed(handler))
	})

	t.Run("should return false if subscriber is no longer subscribed", func(t *testing.T) {
		target := pievent.NewTarget[string]()
		var subscriber Subscriber
		handler := target.Subscribe("", subscriber.Handler)
		target.Unsubscribe(handler)
		assert.False(t, target.IsSubscribed(handler))
	})
}

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

func TestTrackingTarget_Subscribe(t *testing.T) {
	t.Run("should return all subscription handlers", func(t *testing.T) {
		var subscriber1, subscriber2 Subscriber

		target := pievent.NewTarget[string]()
		trackedTarget := pievent.Track(target)
		// when
		handler1 := trackedTarget.Subscribe("1", subscriber1.Handler)
		handler2 := trackedTarget.Subscribe("2", subscriber2.Handler)
		// then
		handlers := trackedTarget.Handlers()
		assert.Equal(t, []pievent.Handler{handler1, handler2}, handlers)
		// and
		assert.True(t, trackedTarget.IsSubscribed(handler1))
		assert.True(t, trackedTarget.IsSubscribed(handler2))
	})
}

func TestTrackingTarget_UnsubscribeAll(t *testing.T) {
	t.Run("should unsubscribe all handlers", func(t *testing.T) {
		var subscriber1, subscriber2 Subscriber

		target := pievent.NewTarget[string]()
		trackedTarget := pievent.Track(target)
		handler1 := trackedTarget.Subscribe(event, subscriber1.Handler)
		handler2 := trackedTarget.Subscribe(event, subscriber2.Handler)
		// when
		trackedTarget.UnsubscribeAll()
		// then
		assert.Empty(t, trackedTarget.Handlers())
		// and
		assert.False(t, trackedTarget.IsSubscribed(handler1))
		assert.False(t, trackedTarget.IsSubscribed(handler2))
		// and
		target.Publish(event)
		subscriber1.AssertNoEventReceived(t)
		subscriber2.AssertNoEventReceived(t)
	})
}

func TestTrackingTarget_Unsubscribe(t *testing.T) {
	t.Run("should unsubscribe handler", func(t *testing.T) {
		var subscriber Subscriber

		target := pievent.NewTarget[string]()
		trackedTarget := pievent.Track(target)
		handler := trackedTarget.Subscribe(event, subscriber.Handler)
		// when
		trackedTarget.Unsubscribe(handler)
		// then
		assert.Empty(t, trackedTarget.Handlers())
		// and
		target.Publish(event)
		subscriber.AssertNoEventReceived(t)
		assert.False(t, trackedTarget.IsSubscribed(handler))
	})
}

type Subscriber struct {
	eventReceived []string
}

func (s *Subscriber) Handler(e string, handler pievent.Handler) {
	s.eventReceived = append(s.eventReceived, e)
}

func (s *Subscriber) AssertEventReceived(t *testing.T, event string) {
	t.Helper()
	require.Len(t, s.eventReceived, 1)
	assert.Equal(t, event, s.eventReceived[0])
}

func (s *Subscriber) AssertNoEventReceived(t *testing.T) {
	t.Helper()
	require.Empty(t, s.eventReceived)
}
