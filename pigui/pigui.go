// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package pigui offers a minimal API for building GUIs.
package pigui

import (
	"slices"

	"github.com/elgopher/pi"
	"github.com/elgopher/pi/pimouse"
)

// New creates a new GUI root element.
//
// To update and draw the element along with its children,
// add it to your game loop by calling Element.Update and Element.Draw.
func New() *Element {
	return &Element{
		area: pi.IntArea{
			W: pi.Screen().W(), // TODO Should not be fixed
			H: pi.Screen().H(),
		},
	}
}

// Attach attaches a new element with the specified size to the parent.
//
// It returns the newly created element.
func Attach(parent *Element, x, y, w, h int) *Element {
	ch := &Element{
		area: pi.IntArea{X: x, Y: y, W: w, H: h},
	}
	parent.Attach(ch)
	return ch
}

type Element struct {
	OnDraw     func(DrawEvent)
	OnUpdate   func(UpdateEvent)
	OnPressed  func(Event)
	OnReleased func(Event)
	OnTapped   func(Event)
	area       pi.Area[int]
	children   []*Element
	pressed    bool
}

// Attach re-attaches an existing element to the parent e.
func (e *Element) Attach(child *Element) {
	e.children = append(e.children, child)
}

// Detach detaches the specified child element from e.
func (e *Element) Detach(child *Element) {
	e.children = slices.DeleteFunc(e.children, func(element *Element) bool {
		return child == element
	})
}

// Update should be called from within pi.Update
// or by any subscriber listening to piloop.EventUpdate events
func (e *Element) Update() {
	prevCamera := pi.Camera
	defer func() {
		pi.Camera = prevCamera
	}()

	pi.Camera.X -= e.area.X // musze przesuwac kamere, zeby dzieci dzieciow zbieraly eventy
	pi.Camera.Y -= e.area.Y

	mousePosition := pimouse.Position.Add(pi.Camera)

	hasPointer := mousePosition.X >= 0 && mousePosition.X < e.area.W &&
		mousePosition.Y >= 0 && mousePosition.Y < e.area.H

	propagate := getPropagateToChildrenFromThePool()

	updateEvent := UpdateEvent{
		Element:                  e,
		HasPointer:               hasPointer,
		propagateToChildren:      propagate,
		propagateToChildrenToken: propagateToChildrenToken,
	}
	if e.OnUpdate != nil {
		e.OnUpdate(updateEvent)
	}

	mouseLeft := pimouse.Duration(pimouse.Left)

	if hasPointer && mouseLeft == 1 {
		e.pressed = true
		if e.OnPressed != nil {
			e.OnPressed(Event{
				Element:    e,
				HasPointer: true,
			})
		}
	} else if e.pressed && mouseLeft == 0 {
		e.pressed = false
		if e.OnReleased != nil {
			e.OnReleased(Event{
				Element:    e,
				HasPointer: hasPointer,
			})
		}
		if hasPointer {
			if e.OnTapped != nil {
				e.OnTapped(Event{
					Element:    e,
					HasPointer: true,
				})
			}
		}
	}

	childrenPropagation := propagate.value
	propagateToChildrenPool.Put(propagate)

	if childrenPropagation {
		for _, child := range e.children {
			child.Update()
		}
	}
}

// Draw should be called from within pi.Draw
// or by any subscriber listening to piloop.EventDraw events
func (e *Element) Draw() {
	prevCamera := pi.Camera
	defer func() {
		pi.Camera = prevCamera
	}()

	pi.Camera.X -= e.area.X
	pi.Camera.Y -= e.area.Y

	prevClip := pi.SetClip(pi.IntArea{
		X: -pi.Camera.X, Y: -pi.Camera.Y,
		W: e.area.W, H: e.area.H,
	})
	defer func() {
		pi.SetClip(prevClip)
	}()

	propagate := getPropagateToChildrenFromThePool()

	mousePosition := pimouse.Position.Add(pi.Camera)
	hasPointer := mousePosition.X >= 0 && mousePosition.X < e.area.W &&
		mousePosition.Y >= 0 && mousePosition.Y < e.area.H

	drawEvent := DrawEvent{
		Element:                  e,
		HasPointer:               hasPointer,
		Pressed:                  e.pressed,
		propagateToChildren:      propagate,
		propagateToChildrenToken: propagateToChildrenToken,
	}
	if e.OnDraw != nil {
		e.OnDraw(drawEvent)
	}

	childrenPropagation := propagate.value
	propagateToChildrenPool.Put(propagate)

	if childrenPropagation {
		for _, child := range e.children {
			child.Draw()
		}
	}
}
func (e *Element) Area() pi.IntArea {
	return e.area
}
