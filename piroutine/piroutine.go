// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package piroutine provides functionality similar to coroutines
// found in other languages (e.g., Lua). It allows you to create Routines,
// which are programs composed of a sequence of steps,
// typically executed over many game frames.
//
// This package helps you write more readable code, but at the cost of
// more difficult debugging and higher resource usage.
// Recommended for cases where your logic can be expressed as
// a sequence of steps, such as creating animations.
package piroutine

import (
	"log"

	"github.com/elgopher/pi/pievent"
	"github.com/elgopher/pi/piloop"
)

// New creates a new Routine composed of the provided steps.
func New(steps ...Step) *Routine {
	if len(steps) == 0 {
		return &Routine{
			stopped: true,
		}
	}

	return &Routine{steps: steps}
}

// Step is a single step in a Routine.
//
// It returns true when the step has finished
// and the next step in the Routine can be executed
type Step func() bool

// Call creates a Routine step that executes the given function
// and immediately advances to the next step
func Call(f func()) Step {
	return func() bool {
		f()
		return true
	}
}

// Printf logs the text and immediately advances to the next Routine step.
func Printf(format string, v ...any) Step {
	return func() bool {
		log.Printf(format, v...)
		return true
	}
}

// Wait creates a Routine step that waits for n resumes
// before advancing to the next step.
func Wait(n int) Step {
	toGo := n + 1
	return func() bool {
		toGo -= 1
		return toGo == 0
	}
}

// SlowDown returns a Routine step that executes f every n updates.
func SlowDown(n int, f func() bool) Step {
	remaining := n
	return func() bool {
		if remaining <= 0 {
			if f() {
				return true
			}

			remaining = n
		} else {
			remaining -= 1
		}

		return false
	}
}

// Routine is a program composed of a sequence of steps,
// typically executed over multiple game frames.
type Routine struct {
	steps       []Step
	currentStep int
	tracing     bool
	stopped     bool
	log         *log.Logger
	name        string
}

func (r *Routine) initLogger() {
	prefix := "[piroutine] "
	if r.name != "" {
		prefix = "[piroutine] " + "[" + r.name + "]"
	}
	r.log = log.New(log.Default().Writer(), prefix, 0)
}

// SetTracing enables or disables logging of each step execution.
//
// Very useful for debugging a Routine.
func (r *Routine) SetTracing(tracing bool) {
	r.initLogger()
	if r.tracing && !tracing {
		r.log.Println("Tracing stopped")
	}
	r.tracing = tracing
	if r.tracing {
		r.log.Println("Tracing started")
	}
}

// Stop terminates the execution of the Routine.
func (r *Routine) Stop() {
	if r.tracing {
		r.log.Println("Stopped on step", r.currentStep)
	}
	r.stopped = true
}

func (r *Routine) Stopped() bool {
	return r.stopped
}

// Resume resumes the execution of the Routine.
//
// It returns true if the Routine has not yet finished
// and can be continued by calling Resume again.
func (r *Routine) Resume() bool {
	for !r.stopped {
		if r.tracing {
			r.log.Println("Calling step", r.currentStep)
		}
		finished := r.steps[r.currentStep]()
		if !finished {
			if r.tracing {
				r.log.Println("Routine suspended")
			}
			return true
		}

		r.currentStep += 1
		if r.currentStep >= len(r.steps) {
			r.stopped = true
			if r.tracing {
				r.log.Println("Routine finished")
			}
			return false
		}
	}
	return false
}

// SetName sets the name of the Routine.
//
// Used for tracing.
func (r *Routine) SetName(name string) {
	r.name = name
	if r.tracing {
		r.initLogger()
	}
}

// ScheduleOn schedules the Routine to resume on the given event.
//
// Each publication of the event will call Resume.
//
// Returns a handler that can be unregistered from piloop.Target
// to stop further resuming.
func (r *Routine) ScheduleOn(event piloop.Event) pievent.Handler {
	return piloop.Target().Subscribe(event, func(_ piloop.Event, handler pievent.Handler) {
		if r.tracing {
			r.log.Printf("Event %s received", event)
		}
		if !r.Resume() {
			piloop.Target().Unsubscribe(handler)
		}
	})
}
