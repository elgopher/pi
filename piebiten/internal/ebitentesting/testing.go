// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package ebitentesting

import (
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

type gameWithOneUpdate struct {
	m    *testing.M
	code int
}

func (g *gameWithOneUpdate) Update() error {
	g.code = g.m.Run()
	return ebiten.Termination
}

func (*gameWithOneUpdate) Draw(*ebiten.Image) {}

func (*gameWithOneUpdate) Layout(int, int) (int, int) {
	return 1, 1
}

// MainWithRunLoop should be Run from TestMain(*testing.M) in order to
// run all tests inside Ebitengine's game loop
func MainWithRunLoop(m *testing.M) {
	g := &gameWithOneUpdate{m: m, code: 1}
	if err := ebiten.RunGame(g); err != nil {
		panic(err) // RunGame does not return an error for ebiten.Termination
	}
	os.Exit(g.code)
}
