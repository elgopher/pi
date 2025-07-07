// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

// Masks used by drawing operations.
//
// These define which bits of the color value are considered
// when reading or writing pixels.
var (
	// ReadMask is the mask applied when reading the draw color.
	ReadMask Color = MaxColors - 1

	// TargetMask is the mask used when reading the target color in sprite operations.
	TargetMask Color = MaxColors - 1

	// ShapeTargetMask is the mask used when reading the target color in shape operations.
	ShapeTargetMask Color = MaxColors - 1
)
