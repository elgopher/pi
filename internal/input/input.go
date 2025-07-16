// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package input

import "github.com/elgopher/pi"

type State[T comparable] struct {
	pressedInputs map[T]*pressedInput
}

func (s *State[T]) Duration(input T) int {
	return s.pressedInput(input).duration()
}

func (s *State[T]) SetStartFrame(input T, frame int) {
	s.pressedInput(input).startFrame = frame
}

func (s *State[T]) SetStopFrame(input T, frame int) {
	s.pressedInput(input).stopFrame = frame
}

func (s *State[T]) pressedInput(input T) *pressedInput {
	if s.pressedInputs == nil {
		s.pressedInputs = map[T]*pressedInput{}
	}
	p, ok := s.pressedInputs[input]
	if !ok {
		p = &pressedInput{startFrame: -1, stopFrame: -1}
		s.pressedInputs[input] = p
	}
	return p
}

type pressedInput struct {
	startFrame, stopFrame int
}

func (p pressedInput) duration() int {
	if p.startFrame < 0 {
		return 0
	}
	if p.startFrame > p.stopFrame {
		return pi.Frame - p.startFrame + 1
	}
	if p.stopFrame == pi.Frame {
		return 1
	}
	return 0
}
