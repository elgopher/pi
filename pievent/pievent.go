// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pievent provides a simple observer pattern implementation for decoupled communication.
//
// The observer (or publish/subscribe) pattern lets different parts of your program
// communicate without knowing about each other's concrete types. Instead of
// directly calling methods on another component, a publisher broadcasts an event,
// and any interested subscribers listening on the same target will receive it.
//
// This approach is useful for building flexible, decoupled systems where
// components shouldn't depend tightly on each other. For example, a player movement
// system can publish "Moved" events, and other systems (like animations, sound,
// or networking) can subscribe to those events and react appropriately
// without being directly linked to the movement logic.
//
// Events in pievent are synchronous: when an event is published, all subscribers
// receive it immediately on the same goroutine. This design keeps things
// simple and predictable (no concurrency issues) but means event handlers
// should be fast to avoid blocking the publisher.
//
// The API is also optimized to avoid allocations in typical usage. Registering
// and publishing events generally does not allocate on the heap,
// making it suitable even for high-frequency game loops.
//
// However, note that pievent is not a silver bullet for all communication needs:
//   - It's best used when multiple independent systems need to react to the same event.
//   - If you only ever have one listener for a signal, it's often better to just
//     call a function or assign a callback directly, which is simpler, easier to debug,
//     and easier to maintain.
//
// Also remember that pievent-based listeners often need careful design:
// they usually use type-switches or if-statements to identify the event type,
// which can make the code more complex if overused.
//
// In short, use pievent when you need many decoupled listeners for an event,
// but prefer simpler direct calls when only one consumer exists.
package pievent

import (
	"log"
	"slices"

	"github.com/elgopher/pi/internal/fileloc"
)

var GlobalTracingOff bool // when true, tracing is globally disabled (zero-allocation mode)

func NewTarget[T comparable]() Target[T] {
	return &target[T]{}
}

// Track creates a new target that tracks all new subscriptions,
// allowing them to be unsubscribed later by calling UnsubscribeAll
// on the returned TrackingTarget.
func Track[T comparable](wrappedTarget Target[T]) *TrackingTarget[T] {
	return &TrackingTarget[T]{wrappedTarget: wrappedTarget}
}

// Target defines an interface for publishing and subscribing to events.
//
// A Target manages event listeners and allows components to publish events
// without knowing who will receive them.
type Target[T comparable] interface {
	// Publish sends an event to all matching subscribers.
	Publish(event T)

	// SubscribeAll registers a listener that receives all events.
	SubscribeAll(listener func(event T, handler Handler)) Handler

	// Subscribe registers a listener for a specific event.
	// If the event is the zero value (for example, an empty string),
	// the listener receives all events (same as SubscribeAll).
	Subscribe(event T, f func(T, Handler)) Handler

	// Unsubscribe removes a previously registered listener.
	Unsubscribe(Handler)

	// IsSubscribed reports whether the given Handler is currently subscribed.
	IsSubscribed(Handler) bool

	// SetTracing enables or disables tracing for this Target.
	SetTracing(enabled bool)
}

// Handler is an identifier returned when subscribing to a Target.
//
// It can be used to later unsubscribe or check subscription status.
type Handler int

type target[T comparable] struct {
	handlers []eventHandler[T]
	tracing  bool
	lastID   Handler
}

var logger = func() *log.Logger {
	return log.New(log.Default().Writer(), "[pievent] ", 0)
}()

func (t *target[T]) Publish(event T) {
	var zeroEvent T

	if t.tracing {
		logger.Printf("Publishing %+v", event)
	}
	for i := 0; i < len(t.handlers); i++ {
		handler := &t.handlers[i]
		if handler.event == zeroEvent || handler.event == event {
			if t.tracing {
				logger.Println("Calling event handler", handler.fileLoc)
			}
			handler.f(event, handler.id)
		}
	}
}

func (t *target[T]) SubscribeAll(f func(T, Handler)) Handler {
	t.lastID++

	var event T
	return t.Subscribe(event, f)
}

func (t *target[T]) Subscribe(event T, f func(T, Handler)) Handler {
	t.lastID++

	t.handlers = append(t.handlers, eventHandler[T]{})
	h := &t.handlers[len(t.handlers)-1]

	h.f = f
	h.id = t.lastID
	h.event = event
	if !GlobalTracingOff {
		h.fileLoc = fileloc.Get("github.com/elgopher/pi/pievent.") // skip pievent wrapper created in Subscribe function
	}

	return t.lastID
}

func (t *target[T]) SetTracing(enabled bool) {
	t.tracing = enabled
}

func (t *target[T]) Unsubscribe(handlerID Handler) {
	for i := 0; i < len(t.handlers); i++ {
		handler := &t.handlers[i]
		if handler.id == handlerID {
			if i < len(t.handlers)-1 {
				copy(t.handlers[i:], t.handlers[i+1:])
			}
			t.handlers = t.handlers[:len(t.handlers)-1]
			return
		}
	}
}

func (t *target[T]) IsSubscribed(handlerID Handler) bool {
	for _, handler := range t.handlers {
		if handler.id == handlerID {
			return true
		}
	}
	return false
}

type eventHandler[T comparable] struct {
	id      Handler
	fileLoc string
	f       func(T, Handler)
	event   T
}

// TrackingTarget wraps a Target and tracks all subscriptions.
//
// It allows you to later unsubscribe all tracked handlers at once.
type TrackingTarget[T comparable] struct {
	wrappedTarget Target[T]
	handlers      []Handler
}

func (t *TrackingTarget[T]) Publish(event T) {
	t.wrappedTarget.Publish(event)
}

func (t *TrackingTarget[T]) Subscribe(event T, listener func(T, Handler)) Handler {
	handler := t.wrappedTarget.Subscribe(event, listener)
	t.handlers = append(t.handlers, handler)
	return handler
}

func (t *TrackingTarget[T]) SubscribeAll(listener func(T, Handler)) Handler {
	handler := t.wrappedTarget.SubscribeAll(listener)
	t.handlers = append(t.handlers, handler)
	return handler
}

func (t *TrackingTarget[T]) SetTracing(enabled bool) {
	t.wrappedTarget.SetTracing(enabled)
}

func (t *TrackingTarget[T]) Handlers() []Handler {
	return slices.Clone(t.handlers)
}

func (t *TrackingTarget[T]) Unsubscribe(handler Handler) {
	t.wrappedTarget.Unsubscribe(handler)
	t.handlers = slices.DeleteFunc(t.handlers, func(i Handler) bool {
		return handler == i
	})
}

func (t *TrackingTarget[T]) IsSubscribed(h Handler) bool {
	return t.wrappedTarget.IsSubscribed(h)
}

func (t *TrackingTarget[T]) UnsubscribeAll() {
	for _, handler := range t.handlers {
		t.wrappedTarget.Unsubscribe(handler)
	}
}
