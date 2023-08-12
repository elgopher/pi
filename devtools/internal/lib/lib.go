// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

// Package lib provides symbols to be used by Yaegi interpreter.
package lib

import (
	"reflect"
	"strings"
)

// Symbols variable stores the map of stdlib symbols per package.
var Symbols = map[string]map[string]reflect.Value{}

// MapTypes variable contains a map of functions which have an interface{} as parameter but
// do something special if the parameter implements a given interface.
var MapTypes = map[reflect.Value][]reflect.Type{}

func AllPackages() []Package {
	var packages []Package
	for pkgWithAlias := range Symbols {
		if strings.Contains(pkgWithAlias, "/") {
			pkg := pkgWithAlias[:strings.LastIndex(pkgWithAlias, "/")]
			alias := pkgWithAlias[strings.LastIndex(pkgWithAlias, "/")+1:]
			packages = append(packages, Package{Path: pkg, Alias: alias})
		}
	}
	return packages
}

type Package struct {
	Path  string // for example: github.com/elgopher/pi
	Alias string // for example: pi
}

func (p Package) IsPiPackage() bool {
	return strings.HasPrefix(p.Path, "github.com/elgopher/pi")
}

func (p Package) IsStdPackage() bool {
	return !p.IsPiPackage() && !strings.HasPrefix(p.Path, "github.com/traefik/yaegi")
}
