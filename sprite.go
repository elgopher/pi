// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"fmt"
	"github.com/elgopher/pi/pimath"
)

// DrawSprite draws the given sprite at (dx, dy) on the current draw target.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func DrawSprite(sprite Sprite, dx, dy int) {
	Stretch(sprite, dx, dy, sprite.W, sprite.H)
}

// DrawCanvas draws the contents of the src Canvas at (x, y) on the current draw target.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func DrawCanvas(src Canvas, x int, y int) {
	DrawSprite(CanvasSprite(src), x, y)
}

// CanvasSprite returns a Sprite covering the entire canvas.
func CanvasSprite(canvas Canvas) Sprite {
	return Sprite{Area: canvas.EntireArea(), Source: canvas}
}

// SpriteFrom creates a Sprite from the specified region of the canvas.
func SpriteFrom(canvas Canvas, x, y, w, h int) Sprite {
	return Sprite{Area: IntArea{X: x, Y: y, W: w, H: h}, Source: canvas}
}

// Stretch draws the given sprite stretched to the specified size.
//
// dx, dy specify the position on the current draw target.
// dw, dh specify the width and height on the current draw target.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func Stretch(sprite Sprite, dx, dy, dw, dh int) {
	// Stretch is fast, so there is no need to implement dedicated DrawSprite or DrawCanvas
	// functions.
	dx -= Camera.X
	dy -= Camera.Y

	stepX := float64(sprite.W) / float64(dw)
	stepY := float64(sprite.H) / float64(dh)

	sourceClip := sprite.Source.EntireArea()
	sourceClipf := Area[float64]{W: float64(sourceClip.W), H: float64(sourceClip.H)}

	src := Area[float64]{X: float64(sprite.X), Y: float64(sprite.Y), W: float64(sprite.W), H: float64(sprite.H)}
	// First, clip the coordinates to ensure they never go outside the source or destination bounds:
	src, sdx, sdy := src.ClippedBy(sourceClipf)

	dst := Area[float64]{X: float64(dx), Y: float64(dy), W: float64(dw), H: float64(dh)}
	clipf := Area[float64]{X: float64(clip.X), Y: float64(clip.Y), W: float64(clip.W), H: float64(clip.H)}

	// Move the destination if the source was offset by clipping.
	dst, ddx, ddy := dst.MovedBy(sdx/stepX, sdy/stepY).ClippedBy(clipf)
	// Shift the source again by the destination normalization offset.
	src, _, _ = src.MovedBy(ddx*stepX, ddy*stepY).ClippedBy(sourceClipf)

	if src.W == 0 || src.H == 0 {
		return
	}

	targetIdx := drawTarget.FlatIndex(int(dst.X), int(dst.Y))
	targetStride := drawTarget.width - int(dst.W)
	srcSource := sprite.Source

	stepXAbs := src.W / dst.W
	stepYAbs := src.H / dst.H

	// start sampling from half-step offsets
	if sprite.FlipX {
		src.X += src.W - stepXAbs/2
		stepXAbs = -stepXAbs
	} else {
		src.X += stepXAbs / 2
	}

	if sprite.FlipY {
		src.Y += src.H - stepYAbs/2
		stepYAbs = -stepYAbs
	} else {
		src.Y += stepYAbs / 2
	}

	srcY := src.Y

	srcMaxX := int(src.X + src.W)
	srcMaxY := int(src.Y + src.H)

	for line := 0.0; line < dst.H; line++ {
		syIndex := pimath.Clamp(int(srcY), 0, srcMaxY-1)
		srcLineIdx := syIndex * srcSource.width

		srcX := src.X
		for cell := 0; cell < int(dst.W); cell++ {
			sxIndex := pimath.Clamp(int(srcX), 0, srcMaxX-1)

			sourceColor := srcSource.data[srcLineIdx+sxIndex] & ReadMask
			targetColor := drawTarget.data[targetIdx] & TargetMask
			drawTarget.data[targetIdx] =
				ColorTables[(sourceColor|targetColor)>>6][sourceColor&(MaxColors-1)][targetColor&(MaxColors-1)]
			srcX += stepXAbs
			targetIdx++
		}
		srcY += stepYAbs
		targetIdx += targetStride
	}
}

// Sprite represents a portion of a Canvas.
type Sprite struct {
	Area[int]

	Source Canvas
	FlipX  bool
	FlipY  bool
}

// WithFlipX returns a new Sprite with the FlipX value updated.
func (s Sprite) WithFlipX(flip bool) Sprite {
	s.FlipX = flip
	return s
}

// WithFlipY returns a new Sprite with the FlipY value updated.
func (s Sprite) WithFlipY(flip bool) Sprite {
	s.FlipY = flip
	return s
}

// WithSizeScaled returns a new Sprite scaled by w and h.
func (s Sprite) WithSizeScaled(w, h float64) Sprite {
	s.W = int(float64(s.W) * w)
	s.H = int(float64(s.H) * h)
	return s
}

// WithSize returns a new Sprite with its size set to the given w and h.
func (s Sprite) WithSize(w, h int) Sprite {
	s.W = w
	s.H = h
	return s
}

// WithSource returns a new Sprite with its Source set to the given Canvas.
func (s Sprite) WithSource(source Canvas) Sprite {
	s.Source = source
	return s
}

func (s Sprite) String() string {
	return fmt.Sprintf("x=%d y=%d w=%d h=%d flipX=%v flipY=%v",
		s.X, s.Y, s.W, s.H, s.FlipX, s.FlipY)
}
