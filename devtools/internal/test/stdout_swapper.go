// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

//go:build !js

package test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// SwapStdout swaps stdout with fake one.
func SwapStdout(t *testing.T) StdoutSwapper {
	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	oldStdout := os.Stdout
	os.Stdout = writer

	return StdoutSwapper{
		prevStdout: oldStdout,
		reader:     reader,
	}
}

type StdoutSwapper struct {
	prevStdout *os.File
	reader     *os.File
}

func (s *StdoutSwapper) BringStdoutBack() {
	_ = os.Stdout.Close()
	os.Stdout = s.prevStdout
}

func (s *StdoutSwapper) ReadOutput(t *testing.T) string {
	all, err := io.ReadAll(s.reader)
	require.NoError(t, err)
	return string(all)
}
