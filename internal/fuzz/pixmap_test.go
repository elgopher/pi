// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package fuzz_test

import (
	"testing"

	"github.com/elgopher/pi"
)

func FuzzPixMap_Pointer(f *testing.F) {
	pixMap := pi.NewPixMap(8, 9).WithClip(1, 1, 3, 5)

	f.Fuzz(func(t *testing.T, x, y, w, h int) {
		pixMap.Pointer(x, y, w, h)
	})
}

func FuzzPixMap_Foreach(f *testing.F) {
	pixMap := pi.NewPixMap(2, 3)
	f.Fuzz(func(t *testing.T, x, y, w, h, dstX, dstY int) {
		pixMap.Foreach(x, y, w, h, func(x, y int, dst []byte) {
			dst[0] = byte(x + y)
		})
	})
}

func FuzzPixMap_Copy_Src_Bigger(f *testing.F) {
	src := pi.NewPixMap(5, 4)
	dst := pi.NewPixMap(2, 3)

	f.Fuzz(func(t *testing.T, x, y, w, h, dstX, dstY int) {
		src.Copy(x, y, w, h, dst, dstX, dstY)
	})
}

func FuzzPixMap_Copy_Dst_Bigger(f *testing.F) {
	src := pi.NewPixMap(2, 3)
	dst := pi.NewPixMap(5, 4)

	f.Fuzz(func(t *testing.T, x, y, w, h, dstX, dstY int) {
		src.Copy(x, y, w, h, dst, dstX, dstY)
	})
}

func FuzzPixMap_Merge(f *testing.F) {
	src := pi.NewPixMap(2, 3)
	dst := pi.NewPixMap(4, 3)

	f.Fuzz(func(t *testing.T, x, y, w, h, dstX, dstY int) {
		src.Merge(x, y, w, h, dst, dstX, dstY, func(dst, src []byte) {
			copy(dst, src)
		})
	})
}
