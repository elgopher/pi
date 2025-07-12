// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piaudio

type LoopType string

const (
	LoopNone    LoopType = "none"    // no loop.
	LoopForward LoopType = "forward" // infinite forward loop.
)
