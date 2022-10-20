// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package vm

var ClippingRegion Region

type Region struct {
	X, Y, W, H int
}
