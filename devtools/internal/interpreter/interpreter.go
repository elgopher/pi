// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package interpreter

import (
	"fmt"
	"go/build"
	"go/scanner"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/traefik/yaegi/interp"

	"github.com/elgopher/pi/devtools/internal/lib"
)

type Instance struct {
	yaegi                   *interp.Interpreter
	alreadyImportedPackages map[string]struct{}
	printHelp               func(topic string) error
}

func New(printHelp func(topic string) error) (Instance, error) {
	yaegi := interp.New(interp.Options{
		GoPath: gopath(), // if GoPath is set then Yaegi does not complain about setting GOPATH.
	})
	err := yaegi.Use(lib.Symbols)
	if err != nil {
		return Instance{}, fmt.Errorf("problem loading pi and stdlib symbols into Yaegi interpreter: %w", err)
	}

	yaegi.ImportUsed()

	instance := Instance{
		yaegi:                   yaegi,
		alreadyImportedPackages: map[string]struct{}{},
		printHelp:               printHelp,
	}

	err = ExportType[noResult](instance)
	if err != nil {
		return Instance{}, fmt.Errorf("problem exporting noResult type: %s", err)
	}

	return instance, nil
}

func gopath() string {
	p := os.Getenv("GOPATH")
	if p == "" {
		p = build.Default.GOPATH
	}
	return p
}

func (i Instance) SetUpdate(update *func()) error {
	return i.usePointerToFunc("github.com/elgopher/pi/pi", "Update", update)
}

func (i Instance) usePointerToFunc(packagePath string, variableName string, f *func()) error {
	err := i.yaegi.Use(interp.Exports{
		packagePath: map[string]reflect.Value{
			variableName: reflect.ValueOf(f).Elem(),
		},
	})

	if err != nil {
		return fmt.Errorf("problem loading %s symbol into Yaegi interpreter: %w", variableName, err)
	}

	return nil
}

func (i Instance) SetDraw(draw *func()) error {
	return i.usePointerToFunc("github.com/elgopher/pi/pi", "Draw", draw)
}

var goIdentifierRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

func (i Instance) Export(name string, f any) error {
	if !goIdentifierRegex.MatchString(name) {
		return ErrInvalidIdentifier{
			message: fmt.Sprintf("'%s' is not a valid Go identifier", name),
		}
	}
	value := reflect.ValueOf(f)

	if err := ensureNameDoesNotClashWithImportedPackages(name); err != nil {
		return err
	}

	if err := ensureNameDoesNotClashWithCommandNames(name); err != nil {
		return err
	}

	err := i.yaegi.Use(interp.Exports{
		"main/main": map[string]reflect.Value{
			name: value,
		},
	})
	if err != nil {
		return fmt.Errorf("problem loading %s symbol into Yaegi interpreter: %w", name, err)
	}

	_, err = i.yaegi.Eval(`import . "main"`)
	if err != nil {
		return fmt.Errorf("problem re-importing main package: %w", err)
	}

	return nil
}

func ensureNameDoesNotClashWithImportedPackages(name string) error {
	for _, pkg := range lib.AllPackages() {
		if name == pkg.Alias {
			return ErrInvalidIdentifier{
				message: fmt.Sprintf(`"%s" clashes with imported package name`, name),
			}
		}
	}

	return nil
}

var allCommandNames = []string{"h", "help", "p", "pause", "r", "resume", "u", "undo", "n", "next"}

func ensureNameDoesNotClashWithCommandNames(name string) error {
	for _, cmd := range allCommandNames {
		if cmd == name {
			return ErrInvalidIdentifier{
				message: fmt.Sprintf(`"%s" clashes with command name`, name),
			}
		}
	}

	return nil
}

func ExportType[T any](i Instance) error {
	var nilValue *T
	var t T
	fullTypeName := fmt.Sprintf("%T", t)

	slice := strings.Split(fullTypeName, ".")
	packageName := slice[0]
	typeName := slice[1]

	if packageName == "main" {
		if err := ensureNameDoesNotClashWithImportedPackages(typeName); err != nil {
			return err
		}

		if err := ensureNameDoesNotClashWithCommandNames(typeName); err != nil {
			return err
		}
	}

	err := i.yaegi.Use(interp.Exports{
		packageName + "/" + packageName: map[string]reflect.Value{
			typeName: reflect.ValueOf(nilValue),
		},
	})
	if err != nil {
		return fmt.Errorf("problem loading %s symbol into Yaegi interpreter: %w", typeName, err)
	}

	_, alreadyImported := i.alreadyImportedPackages[packageName]
	if alreadyImported {
		return nil
	}

	_, err = i.yaegi.Eval(`import "` + packageName + `"`)
	if err != nil {
		return fmt.Errorf("problem re-importing main package: %w", err)
	}

	i.alreadyImportedPackages[packageName] = struct{}{}

	return nil
}

type EvalResult int

const (
	HelpPrinted        EvalResult = 0
	GoCodeExecuted     EvalResult = 1
	Resumed            EvalResult = 2
	Paused             EvalResult = 3
	Undoed             EvalResult = 4
	Continued          EvalResult = 5
	NextFrameRequested EvalResult = 6
)

func (i Instance) Eval(cmd string) (EvalResult, error) {
	trimmedCmd := strings.Trim(cmd, " ")

	if isHelpCommand(trimmedCmd) {
		topic := strings.Trim(strings.TrimLeft(trimmedCmd, "help"), " ")
		return HelpPrinted, i.printHelp(topic)
	} else if trimmedCmd == "next" || trimmedCmd == "n" {
		return NextFrameRequested, nil
	} else if trimmedCmd == "resume" || trimmedCmd == "r" {
		return Resumed, nil
	} else if trimmedCmd == "pause" || trimmedCmd == "p" {
		return Paused, nil
	} else if trimmedCmd == "undo" || trimmedCmd == "u" {
		return Undoed, nil
	} else {
		return i.runGoCode(cmd)
	}
}

func isHelpCommand(trimmedCmd string) bool {
	return strings.HasPrefix(trimmedCmd, "help ") ||
		strings.HasPrefix(trimmedCmd, "h ") ||
		trimmedCmd == "help" ||
		trimmedCmd == "h"
}

func (i Instance) runGoCode(source string) (r EvalResult, e error) {
	defer func() {
		err := recover()
		if err != nil {
			r = GoCodeExecuted
			e = fmt.Errorf("panic when running Yaegi Go interpreter: %s", err)
		}
	}()

	// Yaegi returns the last result computed by the interpreter which is unfortunate, because we don't know
	// if source is a statement or an expression. If source is a statement and previously source was an expression
	// then Yaegi returns the same result again. That's why following code overrides last computed result to the value
	// used as a discriminator for no result.
	source = "interpreter.noResult{}; " + source

	res, err := i.yaegi.Eval(source)
	if err != nil {
		if shouldContinueOnError(err, source) {
			return Continued, nil
		}

		return GoCodeExecuted, convertYaegiError(err)
	}

	printResult(res)

	return GoCodeExecuted, nil
}

var (
	yaegiCfgErrorPattern     = regexp.MustCompile(`^(\d+:)(\d+:)(.*)`) // control flow graph generation error
	yaegiScannerErrorPattern = regexp.MustCompile(`^_\.go:(\d+:)(\d+:)(.*)`)
)

// Yaegi reports wrong character number in error messages. It is better to remove this information completely.
func convertYaegiError(err error) error {
	errStr := err.Error()

	switch {
	case yaegiCfgErrorPattern.MatchString(errStr):
		return fmt.Errorf(
			yaegiCfgErrorPattern.ReplaceAllString(errStr, "$1$3"),
		)
	case yaegiScannerErrorPattern.MatchString(errStr):
		return fmt.Errorf(
			yaegiScannerErrorPattern.ReplaceAllString(errStr, "$1$3"),
		)
	default:
		return err
	}
}

// shouldContinueOnError returns true if the error can be safely ignored
// to let the caller grab one more line before retrying to parse its input.
func shouldContinueOnError(err error, source string) bool {
	errorsList, ok := err.(scanner.ErrorList)
	if !ok || len(errorsList) < 1 {
		return false
	}

	e := errorsList[0]

	msg := e.Msg
	if strings.HasSuffix(msg, "found 'EOF'") {
		return true
	}

	if msg == "raw string literal not terminated" {
		return true
	}

	if strings.HasPrefix(msg, "expected operand, found '}'") && !strings.HasSuffix(source, "}") {
		return true
	}

	return false
}

func printResult(res reflect.Value) {
	if res.IsValid() {
		kind := res.Type().Kind()

		if kind == reflect.Struct {
			_, noResultFound := res.Interface().(noResult)
			if noResultFound {
				return
			}
		}

		fmt.Printf("%+v: %+v\n", res.Type(), res)
	}
}

// noResult is a special struct which is used to check if code run by the interpreter returned something
type noResult struct{}
