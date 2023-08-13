// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pi

import (
	"fmt"

	"github.com/elgopher/pi/internal/sfmt"
)

// PixMap is a generic data structure for manipulating any kind of pixel data - screen, sprite-sheet etc.
// PixMap uses a single byte (8 bits) for storing single color/pixel. This means that max 256 colors
// can be used. PixMap can also be used for maps which not necessary contain pixel colors, such as world map
// (as long as only 256 different tiles are used).
//
// To create PixMap please use either NewPixMap or NewPixMapWithPix function.
//
// All PixMap functions (besides Clear and ClearCol) take into account the clipping region.
type PixMap struct {
	pix    []byte
	width  int
	height int
	clip   Region

	zeroPix      []byte
	wholeLinePix []byte
}

// NewPixMap creates new instance of PixMap with specified size.
// Width and height cannot be negative.
func NewPixMap(width, height int) PixMap {
	if width < 0 {
		panic("negative PixMap with")
	}
	if height < 0 {
		panic("negative PixMap height")
	}

	pix := make([]byte, width*height)

	return PixMap{
		pix:          pix,
		width:        width,
		height:       height,
		clip:         Region{W: width, H: height},
		zeroPix:      make([]byte, len(pix)),
		wholeLinePix: make([]byte, width),
	}
}

// NewPixMapWithPix creates new instance of PixMap using the slice of pixel colors
// as a source. pix slice contains colors for the entire PixMap.
// Pixels are organized from left to right, top to bottom. Slice element
// number 0 has pixel located in the top-left corner. Slice element number 1
// has pixel color on the right, and so on.
//
// The lineWidth is the width of PixMap. Height is calculated by dividing pixels
// by lineWith.
//
// This function is handy when you already have a pix slice and want to create a PixMap
// out of it. This function does not allocate anything on the heap.
func NewPixMapWithPix(pix []byte, lineWidth int) PixMap {
	if lineWidth < 0 {
		panic("PixMap lineWidth cant be negative")
	}

	if lineWidth == 0 {
		if len(pix) > 0 {
			panic("PixMap lineWidth cant be zero when pix slice is not empty")
		}

		return PixMap{}
	}

	if len(pix)%lineWidth != 0 {
		panic("invalid pixmap lineWidth. Length of pix slice must be multiple of lineWidth.")
	}

	height := len(pix) / lineWidth

	return PixMap{
		pix:          pix,
		width:        lineWidth,
		height:       height,
		clip:         Region{W: lineWidth, H: height},
		zeroPix:      make([]byte, len(pix)),
		wholeLinePix: make([]byte, lineWidth),
	}
}

// Pix return pixel colors. Pixels are organized from left to right,
// top to bottom. Slice element number 0 has pixel located
// in the top-left corner. Slice element number 1 has pixel color
// on the right and so on.
//
// Returned slice can be freely read and updated. Useful when
// you want to use your own functions for pixel manipulation.
func (p PixMap) Pix() []byte {
	return p.pix
}

// Width returns the width of PixMap (lineWidth), without taking
// into account the current clipping region.
func (p PixMap) Width() int {
	return p.width
}

// Height returns the height of PixMap, without taking into account
// the current clipping region.
func (p PixMap) Height() int {
	return p.height
}

// Clip returns the clipping region, which specifies which fragment
// of the PixMap can be accessed by its functions. By default, clipping
// region has a size of the entire PixMap (no clipping).
func (p PixMap) Clip() Region {
	return p.clip
}

// WithClip creates a new PixMap which has a different clipping region.
// The newly created PixMap still refers to the same pixels though.
func (p PixMap) WithClip(x, y, w, h int) PixMap {
	if x < 0 {
		w += x
		x = 0
	}

	if y < 0 {
		h += y
		y = 0
	}

	if x+w > p.width {
		w = p.width - x
	}

	if y+h > p.height {
		h = p.height - y
	}

	p.clip = Region{X: x, Y: y, W: w, H: h}

	return p
}

// Clear clears the entire PixMap with color 0. It does not take into account the clipping region.
func (p PixMap) Clear() {
	copy(p.pix, p.zeroPix)
}

// ClearCol clears the entire PixMap with specified color. It does not take into account the clipping region.
func (p PixMap) ClearCol(col byte) {
	line := p.lineOfColor(col, p.width)
	pix := p.pix

	copy(pix, line)
	for i := 1; i < p.height; i++ {
		pix = pix[p.width:]
		copy(pix, line)
	}
}

// Pointer finds the index of (x,y) coordinates and returns a pointer to pixel data at this position
// in Pointer.Pix.
//
// If the pixel is outside the clipping region, then the closest pixel to the bottom-right
// is returned. The difference in position is returned in Pointer.DeltaX and Pointer.DeltaY.
//
// If the pixel is either below or right after the clipping region then ok=false and empty
// Pointer are returned.
func (p PixMap) Pointer(x, y, w, h int) (ptr Pointer, ok bool) {
	if w <= 0 || h <= 0 {
		return ptr, false
	}

	clip := p.clip
	if x >= clip.X+clip.W {
		return ptr, false
	}

	if y >= clip.Y+clip.H {
		return ptr, false
	}

	if x+w <= clip.X {
		return ptr, false
	}

	if y+h <= clip.Y {
		return ptr, false
	}

	var dx, dy int

	if x < clip.X {
		dx = clip.X - x
		x += dx
		w -= dx
	}

	if y < clip.Y {
		dy = clip.Y - y
		y += dy
		h -= dy
	}

	pix := p.pix[y*p.width+x:]

	return Pointer{
		DeltaX:          dx,
		DeltaY:          dy,
		Pix:             pix,
		RemainingPixels: MinInt(w, clip.X+clip.W-x),
		RemainingLines:  MinInt(h, clip.Y+clip.H-y),
	}, true
}

// Pointer is a low-level struct for fast pixel processing created by PixMap.Pointer function.
type Pointer struct {
	// Pix is a slice of PixMap.Pix at the position specified when calling PixMap.Pointer function.
	Pix             []byte
	DeltaX, DeltaY  int
	RemainingPixels int // in line
	RemainingLines  int
}

// Copy copies the region specified by x, y, w, h into dst PixMap at dstX,dstY position.
func (p PixMap) Copy(x, y, w, h int, dst PixMap, dstX, dstY int) {
	dstPtr, srcPtr := p.pointersForCopy(x, y, w, h, dst, dstX, dstY)

	remainingLines := MinInt(dstPtr.RemainingLines, srcPtr.RemainingLines)

	if remainingLines == 0 {
		return
	}

	remainingPixels := MinInt(dstPtr.RemainingPixels, srcPtr.RemainingPixels)

	copy(dstPtr.Pix[:remainingPixels], srcPtr.Pix)
	for i := 1; i < remainingLines; i++ {
		dstPtr.Pix = dstPtr.Pix[dst.width:]
		srcPtr.Pix = srcPtr.Pix[p.width:]
		copy(dstPtr.Pix[:remainingPixels], srcPtr.Pix)
	}
}

// Merge merges destination with source by running merge operation for each destination line.
func (p PixMap) Merge(x, y, w, h int, dst PixMap, dstX, dstY int, merge func(dst, src []byte)) {
	dstPtr, srcPtr := p.pointersForCopy(x, y, w, h, dst, dstX, dstY)

	remainingLines := MinInt(dstPtr.RemainingLines, srcPtr.RemainingLines)

	if remainingLines == 0 {
		return
	}

	remainingPixels := MinInt(dstPtr.RemainingPixels, srcPtr.RemainingPixels)

	merge(dstPtr.Pix[:remainingPixels], srcPtr.Pix)
	for i := 1; i < remainingLines; i++ {
		dstPtr.Pix = dstPtr.Pix[dst.width:]
		srcPtr.Pix = srcPtr.Pix[p.width:]
		merge(dstPtr.Pix[:remainingPixels], srcPtr.Pix)
	}
}

// Foreach runs the update function on PixMap fragment specified by x, y, w and h.
//
// The update function accepts entire line to increase the performance.
func (p PixMap) Foreach(x, y, w, h int, update func(x, y int, dst []byte)) {
	if update == nil {
		return
	}

	ptr, ok := p.Pointer(x, y, w, h)
	if !ok {
		return
	}

	x = x + ptr.DeltaX
	y = y + ptr.DeltaY

	update(x, y, ptr.Pix[:ptr.RemainingPixels])
	for i := 1; i < ptr.RemainingLines; i++ {
		ptr.Pix = ptr.Pix[p.width:]
		update(x, y+i, ptr.Pix[:ptr.RemainingPixels])
	}
}

// Set sets the pixel color at given position.
func (p PixMap) Set(x, y int, col byte) {
	if x < p.clip.X {
		return
	}
	if y < p.clip.Y {
		return
	}
	if x >= p.clip.X+p.clip.W {
		return
	}
	if y >= p.clip.Y+p.clip.H {
		return
	}

	p.pix[y*p.width+x] = col
}

// Get returns the pixel color at given position. If coordinates
// are outside clipping region than color 0 is returned.
func (p PixMap) Get(x, y int) byte {
	if x < p.clip.X {
		return 0
	}
	if y < p.clip.Y {
		return 0
	}
	if x >= p.clip.X+p.clip.W {
		return 0
	}
	if y >= p.clip.Y+p.clip.H {
		return 0
	}

	return p.pix[y*p.width+x]
}

func (p PixMap) lineOfColor(col byte, length int) []byte {
	line := p.wholeLinePix[:length]
	for i := 0; i < len(line); i++ {
		p.wholeLinePix[i] = col
	}
	return line
}

func (p PixMap) pointersForCopy(srcX int, srcY int, w int, h int, dst PixMap, dstX int, dstY int) (Pointer, Pointer) {
	dstPtr, _ := dst.Pointer(dstX, dstY, w, h)
	srcPtr, _ := p.Pointer(srcX, srcY, w, h)

	// both maps must be moved by the same DeltaX and DeltaY
	if srcPtr.DeltaX > dstPtr.DeltaX {
		dstPtr = addX(dstPtr, srcPtr.DeltaX-dstPtr.DeltaX)
	} else {
		srcPtr = addX(srcPtr, dstPtr.DeltaX-srcPtr.DeltaX)
	}

	if srcPtr.DeltaY > dstPtr.DeltaY {
		dstPtr = addY(dstPtr, srcPtr.DeltaY-dstPtr.DeltaY, dst.width)
	} else {
		srcPtr = addY(srcPtr, dstPtr.DeltaY-srcPtr.DeltaY, p.width)
	}

	return dstPtr, srcPtr
}

func addX(p Pointer, n int) Pointer {
	p.DeltaX += n
	p.RemainingPixels -= n
	if p.RemainingPixels <= 0 {
		return Pointer{}
	}
	p.Pix = p.Pix[n:]
	return p
}

func addY(p Pointer, n int, lineWidth int) Pointer {
	p.DeltaY += n
	p.RemainingLines -= n
	if p.RemainingLines <= 0 {
		return Pointer{}
	}
	p.Pix = p.Pix[n*lineWidth:]
	return p
}

// String returns PixMap as string for debugging purposes.
func (p PixMap) String() string {
	return fmt.Sprintf("{width:%d, height:%d, clip:%+v, pix:%s}",
		p.width, p.height, p.clip, sfmt.FormatBigSlice(p.pix, 1024))
}
