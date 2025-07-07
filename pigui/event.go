// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pigui

type Event struct {
	Element    *Element
	HasPointer bool
}

type UpdateEvent struct {
	Element                  *Element
	HasPointer               bool
	propagateToChildren      *propagateToChildren
	propagateToChildrenToken int
}

// call it if you don't want the event to be propagated to children
func (e UpdateEvent) StopPropagation() {
	e.propagateToChildren.set(false, e.propagateToChildrenToken)
}

type DrawEvent struct {
	Element                  *Element
	HasPointer               bool
	propagateToChildren      *propagateToChildren
	propagateToChildrenToken int
	Pressed                  bool
}

// call it if you don't want the event to be propagated to children
func (e DrawEvent) StopPropagation() {
	e.propagateToChildren.set(false, e.propagateToChildrenToken)
}
