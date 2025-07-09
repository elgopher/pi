// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package internal

import (
	"github.com/elgopher/pi"
	"github.com/elgopher/pi/piring"
)

var theScreenRecorder = newScreenRecorder(128)

// SetScreenSnapshotHistorySize sets the new number of screen snapshots held in memory
// Calling this method clears the screen snapshots history
func SetScreenSnapshotHistorySize(historySize int) {
	theScreenRecorder = newScreenRecorder(historySize)
}

func newScreenRecorder(historySize int) *screenRecorder {
	buffer := piring.NewBuffer[screenSnapshot](historySize)

	return &screenRecorder{snapshots: buffer}
}

type screenRecorder struct {
	snapshots *piring.Buffer[screenSnapshot]
	shift     int // which element from the end is currently selected
}

type screenSnapshot struct {
	canvas         pi.Canvas
	paletteMapping pi.PaletteMap
	palette        pi.PaletteArray
}

func (s *screenRecorder) Save() {
	snapshot := s.snapshots.NextWritePointer()
	screen := pi.Screen()
	// reuse canvas if possible
	if snapshot.canvas.W() == screen.W() && snapshot.canvas.H() == screen.H() {
		snapshot.canvas.SetData(screen.Data())
	} else {
		snapshot.canvas = screen.Clone()
	}
	snapshot.palette = pi.Palette
	snapshot.paletteMapping = pi.PaletteMapping

	s.shift = 0
}

func (s *screenRecorder) HasPrev() bool {
	return -s.shift+1 <= s.snapshots.Len()
}

func (s *screenRecorder) ShowPrev() bool {
	if -s.shift+1 > s.snapshots.Len() {
		return false
	}
	s.shift -= 1
	s.showCurrent()
	return true
}

func (s *screenRecorder) ShowNext() bool {
	if s.shift >= -1 {
		return false
	}
	s.shift += 1
	s.showCurrent()
	return true
}

func (s *screenRecorder) showCurrent() {
	snapshot := s.snapshots.PointerTo(s.snapshots.Len() + s.shift)
	pi.Screen().SetData(snapshot.canvas.Data())
	pi.Palette = snapshot.palette
	pi.PaletteMapping = snapshot.paletteMapping
}

func (s *screenRecorder) Reset() {
	s.snapshots.Reset()
	s.shift = 0
}

func (s *screenRecorder) GoToLast() {
	s.shift = 0
}
