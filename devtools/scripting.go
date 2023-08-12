// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package devtools

import (
	"fmt"

	"github.com/elgopher/pi/devtools/internal/help"
	"github.com/elgopher/pi/devtools/internal/interpreter"
	"github.com/elgopher/pi/devtools/internal/snapshot"
	"github.com/elgopher/pi/devtools/internal/terminal"
)

var interpreterInstance interpreter.Instance

func init() {
	i, err := interpreter.New(help.PrintHelp)
	if err != nil {
		panic("problem creating interpreter instance: " + err.Error())
	}

	interpreterInstance = i
}

// Export makes the v visible in terminal. v can be a variable or a function.
// It will be visible under the specified name.
//
// If you want to be able to update the value of v in the terminal, please pass a pointer to v. For example:
//
//	v := 1
//	devtools.Export("v", &v)
//
// Then in the terminal write:
//
//	v = 2
func Export(name string, v any) {
	if err := interpreterInstance.Export(name, v); err != nil {
		panic("devtools.Export failed: " + err.Error())
	}
}

// ExportType makes the type T visible in terminal. It will be visible under the name "package.Name".
// For example type "github.com/a/b/c/pkg.SomeType" will be visible as "pkg.SomeType". This function
// also automatically imports the package.
func ExportType[T any]() {
	if err := interpreter.ExportType[T](interpreterInstance); err != nil {
		panic("devtools.ExportType failed: " + err.Error())
	}
}

func evaluateNextCommandFromTerminal() {
	select {
	case cmd := <-terminal.Commands:
		result, err := interpreterInstance.Eval(cmd)
		if err != nil {
			fmt.Println(err)
			terminal.CommandsProcessed <- terminal.Done

			return
		}

		switch result {
		case interpreter.GoCodeExecuted:
			snapshot.Take()
		case interpreter.Resumed:
			resumeGame()
		case interpreter.Paused:
			pauseGame()
		case interpreter.Undoed:
			snapshot.Undo()
			fmt.Println("Undoed last draw operation")
		}

		if result == interpreter.Continued {
			terminal.CommandsProcessed <- terminal.MoreInputNeeded
		} else {
			terminal.CommandsProcessed <- terminal.Done
		}
	default:
	}
}
