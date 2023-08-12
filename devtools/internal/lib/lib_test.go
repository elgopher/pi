// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package lib_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elgopher/pi/devtools/internal/lib"
)

var (
	piPackage  = lib.Package{Path: "github.com/elgopher/pi", Alias: "pi"}
	fmtPackage = lib.Package{Path: "fmt", Alias: "fmt"}
)

func TestAllPackages(t *testing.T) {
	t.Run("should return packages", func(t *testing.T) {
		packages := lib.AllPackages()
		assert.NotEmpty(t, packages)

		assert.Contains(t, packages, piPackage)
		assert.Contains(t, packages, fmtPackage)
	})
}

func TestPackage_IsStdPackage(t *testing.T) {
	assert.True(t, fmtPackage.IsStdPackage())
	assert.False(t, piPackage.IsStdPackage())
}

func TestPackage_IsPiPackage(t *testing.T) {
	assert.False(t, fmtPackage.IsPiPackage())
	assert.True(t, piPackage.IsPiPackage())
}
