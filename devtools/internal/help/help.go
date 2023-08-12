// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package help

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/elgopher/pi/devtools/internal/lib"
)

var NotFound = fmt.Errorf("no help found")

func PrintHelp(topic string) error {
	switch topic {
	case "":
		fmt.Println("This is interactive terminal. " +
			"You can write Go code here, which will run immediately. " +
			"You can use all Pi packages: pi, key, state, snap, font, image and " +
			"selection of standard packages: " + strings.Join(stdPackages(), ", ") + ". " +
			"\n\n" +
			"Type help topic for more information. For example: help pi or help pi.Spr" +
			"\n\n" +
			"Available commands: help [h], pause [p], resume [r], undo [u]",
		)
		return nil
	default:
		return goDoc(topic)
	}
}

func stdPackages() []string {
	var packages []string
	for _, p := range lib.AllPackages() {
		if p.IsStdPackage() {
			packages = append(packages, p.Alias)
		}
	}
	sort.Strings(packages)
	return packages
}

func goDoc(symbol string) error {
	symbol = completeSymbol(symbol)
	if symbolNotSupported(symbol) {
		return NotFound
	}

	fmt.Println("###############################################################################")

	var args []string
	args = append(args, "doc")
	if shouldShowDetailedDescriptionForSymbol(symbol) {
		args = append(args, "-all")
	}
	args = append(args, symbol)
	command := exec.Command("go", args...)
	command.Stdout = bufio.NewWriter(os.Stdout)

	if err := command.Run(); err != nil {
		var exitErr *exec.ExitError
		if isExitErr := errors.As(err, &exitErr); isExitErr && exitErr.ExitCode() == 1 {
			return NotFound
		}

		return fmt.Errorf("problem getting help: %w", err)
	}

	return nil
}

func completeSymbol(symbol string) string {
	packages := lib.AllPackages()

	for _, p := range packages {
		if p.Alias == symbol {
			return p.Path
		}
	}

	for _, p := range packages {
		prefix := p.Alias + "."
		if strings.HasPrefix(symbol, prefix) {
			return p.Path + "." + symbol[len(prefix):]
		}
	}

	return symbol
}

func symbolNotSupported(symbol string) bool {
	packages := lib.AllPackages()

	for _, p := range packages {
		prefix := p.Path + "."
		if strings.HasPrefix(symbol, prefix) || symbol == p.Path {
			return false
		}
	}

	return true
}

var symbolsWithDetailedDescription = []string{
	"github.com/elgopher/pi.Button",
	"github.com/elgopher/pi.MouseButton",
	"github.com/elgopher/pi/key.Button",
}

func shouldShowDetailedDescriptionForSymbol(symbol string) bool {
	for _, s := range symbolsWithDetailedDescription {
		if symbol == s {
			return true
		}
	}

	return false
}
