// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package gameloop

var UpdateFunctions []func()

func Update() {
	for _, f := range UpdateFunctions {
		f()
	}
}
