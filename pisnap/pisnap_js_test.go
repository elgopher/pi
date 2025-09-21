// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pisnap_test

import (
	"testing"

	"github.com/elgopher/pi/pisnap"

	"github.com/stretchr/testify/assert"
)

func TestCaptureOrErr(t *testing.T) {
	file, err := pisnap.CaptureOrErr()
	assert.Error(t, err)
	assert.Empty(t, file)
}
