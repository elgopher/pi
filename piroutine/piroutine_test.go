// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package piroutine_test

import (
	"testing"

	"github.com/elgopher/pi/piroutine"
)

func TestPiroutine(t *testing.T) {
	r := piroutine.New(
		piroutine.Wait(2),
		piroutine.Printf("abc"),
		piroutine.Wait(5),
	)
	r.SetTracing(true)
	for r.Resume() {
	}
}
