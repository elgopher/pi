// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"iter"
	"log"
	"slices"
	"strings"

	"github.com/elgopher/pi/internal"
)

// NewCanvas creates a new Canvas with the specified dimensions.
func NewCanvas(w, h int) Canvas {
	return NewSurface[Color](w, h)
}

// NewSurface creates a new Surface with the specified dimensions.
func NewSurface[T any](w, h int) Surface[T] {
	return Surface[T]{
		width:  w,
		height: h,
		data:   make([]T, w*h),
	}
}

// Canvas is used to store pixels
type Canvas = Surface[Color]

// Surface represents a 2D grid for storing arbitrary data.
//
// Surface is a generic container that maps 2D coordinates (X, Y)
// to values of any type T. It can be used to represent game maps,
// tile layers, simulation grids, or any other 2D spatial data.
//
// Surface is already used in Pi to store pixels in the Canvas type.
//
// Example usage:
//
//	s := NewSurface[int](10, 10)
//	s.Set(3, 4, 42)
//	value := s.Get(3, 4)
type Surface[T any] struct {
	data          []T
	width, height int
}

func (m Surface[T]) String() string {
	var b strings.Builder
	for _, line := range m.LinesIterator(m.EntireArea()) { // escapes to heap
		b.WriteString(fmt.Sprintf("%+v", line))
		b.WriteByte('\n')
	}
	return b.String()
}

func (m Surface[T]) Data() []T {
	return m.data
}

func (m Surface[T]) SetData(data []T) {
	copy(m.data, data)
}

func (m Surface[T]) W() int {
	return m.width
}

func (m Surface[T]) H() int {
	return m.height
}

func (m Surface[T]) Set(x, y int, value T) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	if x >= m.width {
		return
	}
	if y >= m.height {
		return
	}

	idx := y*m.width + x
	m.data[idx] = value
}

func (m Surface[T]) SetMany(x, y int, values ...T) {
	idx := m.FlatIndex(x, y)
	if idx < 0 {
		if -idx > len(values) {
			// nothing to draw
			return
		}
		values = values[-idx:]
		idx = 0
	}
	if idx > len(m.data) {
		// out of bounds
		return
	}

	copy(m.data[idx:], values)
}

func (m Surface[T]) SetAll(values ...T) {
	m.SetMany(0, 0, values...)
}

func (m Surface[T]) SetArea(area IntArea, values ...T) {
	m.setAreaStride(area, values, area.W)
}

func (m Surface[T]) setAreaStride(area IntArea, values []T, stride int) {
	if len(values) < area.Size() {
		err := fmt.Sprintf("invalid values length: %d. Must be %d or more", len(values), area.Size())
		panic(err)
	}

	clippedArea, dx, dy := area.ClippedBy(IntArea{W: m.width, H: m.height})
	if clippedArea.Size() < 0 {
		return
	}

	offset := dx + dy*stride
	if offset >= len(values) {
		// nothing to draw
		return
	}
	values = values[offset:]

	start := 0
	for _, line := range m.LinesIterator(clippedArea) {
		copy(line, values[start:])
		start += stride
	}
}

func (m Surface[T]) SetSurface(x int, y int, src Surface[T]) {
	dstArea := IntArea{X: x, Y: y, W: src.width, H: src.height}
	m.setAreaStride(dstArea, src.data, src.width)
}

func (m Surface[T]) Clone() Surface[T] {
	m.data = slices.Clone(m.data)
	return m
}

// CloneArea clones the specified area of the surface.
//
// If the area extends outside the bounds, the missing regions
// are filled with zero values.
func (m Surface[T]) CloneArea(area IntArea) Surface[T] {
	clone := NewSurface[T](area.W, area.H)
	for y := area.Y; y < area.Y+area.H; y++ {
		for x := area.X; x < area.X+area.W; x++ {
			v := m.Get(x, y)
			clone.Set(x-area.X, y-area.Y, v)
		}
	}
	return clone
}

// FlatIndex returns the index in Surface.Data for the given coordinates.
func (m Surface[T]) FlatIndex(x, y int) int {
	return y*m.width + x
}

// Get returns the value at the given coordinates.
//
// If the coordinates are out of bounds, it returns the zero value of T.
func (m Surface[T]) Get(x, y int) T {
	var zero T
	if x < 0 {
		return zero
	}
	if y < 0 {
		return zero
	}
	if x >= m.width {
		return zero
	}
	if y >= m.height {
		return zero
	}

	return m.data[y*m.width+x]
}

// Get2 returns zero value if there are no more elements.
func (m Surface[T]) Get2(x, y int) (T, T) {
	return m.Get(x, y), m.Get(x+1, y) // optimize someday
}

// Get3 returns zero value if there are no more elements.
func (m Surface[T]) Get3(x, y int) (T, T, T) {
	return m.Get(x, y), m.Get(x+1, y), m.Get(x+2, y) // optimize someday
}

// Returns nil if y is out of bounds
func (m Surface[T]) GetLine(y int) []T {
	if y < 0 {
		return nil
	}
	if y >= m.height {
		return nil
	}

	return m.data[y*m.width : (y+1)*m.width]
}

// LinesIterator returns an iterator for reading and writing data in lines.
//
// This is a very efficient way to process regions in 2D space.
// LinesIterator enables writing advanced, low-level code.
//
// The area must be clipped by m; otherwise, it will panic.
func (m Surface[T]) LinesIterator(area IntArea) iter.Seq2[Position, []T] {
	return func(yield func(pos Position, line []T) bool) {
		i := m.FlatIndex(area.X, area.Y)
		maxY := area.Y + area.H
		for y := area.Y; y < maxY; y++ {
			if !yield(Position{area.X, y}, m.data[i:i+area.W]) {
				return
			}
			i += m.width
		}
	}
}

func (m Surface[T]) EntireArea() IntArea {
	return IntArea{W: m.width, H: m.height}
}

func (m Surface[T]) Clear(v T) {
	// color only the first line
	var firstLine = m.data[:m.width]
	for i := 0; i < m.width; i++ {
		firstLine[i] = v
	}

	// and then copy the first line into rest of the lines
	dst := m.data[m.width:]
	for y := 1; y < m.height; y++ {
		copy(dst, firstLine)
		dst = dst[m.width:]
	}
}

// DecodeCanvas decodes a PNG file into a Canvas using the current Palette.
//
// Must be called after the Palette has been set.
//
// This function can be slow for large images that do not use indexed color mode,
// or that use indexed color mode with a different palette. Whenever possible,
// use images with indexed color mode and the same Palette. Doing so will also
// simplify your workflow, since you can inspect color indexes directly in your
// graphics editor.
func DecodeCanvas(pngFile []byte) Canvas {
	m, err := DecodeCanvasOrErr(pngFile)
	if err != nil {
		panic(err)
	}
	return m
}

// DecodeCanvasOrErr is like DecodeCanvas but returns an error.
//
// This function is useful when you need to validate user-provided PNG data.
// It returns an error if the input is not a valid PNG file or if the image
// contains more than 256 colors.
func DecodeCanvasOrErr(pngFile []byte) (Canvas, error) {
	img, err := png.Decode(bytes.NewReader(pngFile))
	if err != nil {
		return Canvas{}, fmt.Errorf("PNG decoding failed: %w", err)
	}

	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()

	canvas := NewCanvas(dx, dy)

	palettedImage, isIndexedImage := img.(*image.Paletted)
	if isIndexedImage && samePalette(palettedImage.Palette) {
		// fast path
		copy(canvas.data, palettedImage.Pix)
		return canvas, nil
	}

	// slow path
	if dx*dy > 256*1024 {
		log.Println("Decoding big png image which has not indexed color mode. " +
			"The operation will be very slow and can seriously slow down game " +
			"startup time. Consider using png files with indexed color mode, " +
			"all files with the same palette as your game palette (same colors, " +
			"same order)")
	}

	closestColor := internal.ClosestColorPicker[RGB, Color]{
		Palette: Palette,
		Cache:   make(map[color.Color]Color),
	}

	offset := 0

	for y := 0; y < dy; y++ {
		line := canvas.data[offset : offset+dx]
		for x := 0; x < dx; x++ {
			c := img.At(x, y)
			closest, err := closestColor.IndexInPalette(c)
			if err != nil {
				return canvas, err //nolint:wrapcheck
			}
			line[x] = closest
		}
		offset += dx
	}
	return canvas, nil
}

func samePalette(palette color.Palette) bool {
	for i := 0; i < len(palette); i++ {
		r, g, b, _ := palette[i].RGBA()
		c := r&0xff<<16 + g&0xff<<8 + b&0xff
		if c != uint32(Palette[i]) {
			return false
		}
	}

	return true
}
