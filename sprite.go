// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"fmt"
)

// Spr draws the given sprite at (dx, dy) on the current draw target.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func Spr(sprite Sprite, dx, dy int) {
	Stretch(sprite, dx, dy, sprite.W, sprite.H)
}

// Blit draws the contents of the src Canvas at (x, y) on the current draw target.
//
// It takes into account the camera position, clipping region,
// color tables, and masks.
func Blit(src Canvas, x int, y int) {
	Spr(CanvasSprite(src), x, y)
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
	// Stretch is fast, so there is no need to implement dedicated Spr or Blit
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

	if sprite.FlipY {
		src.Y += float64(dh-1) * stepY
		stepY *= -1
	}

	if sprite.FlipX {
		src.X += float64(dw-1) * stepX
		stepX *= -1
	}

	srcX, srcY := src.X, src.Y

	for line := 0.0; line < dst.H; line++ {
		srcLineIdx := int(srcY) * srcSource.width // multiplication, but only once per line, so it's not a performance problem

		for cell := 0.0; cell < dst.W; cell++ {
			sourceColor := srcSource.data[srcLineIdx+int(srcX)] & ReadMask
			targetColor := drawTarget.data[targetIdx] & TargetMask
			drawTarget.data[targetIdx] =
				ColorTables[(sourceColor|targetColor)>>6][sourceColor&(MaxColors-1)][targetColor&(MaxColors-1)]
			srcX += stepX
			targetIdx++
		}
		srcX = src.X
		srcY += stepY
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
