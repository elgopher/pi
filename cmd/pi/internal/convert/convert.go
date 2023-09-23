// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package convert

import (
	"fmt"
	"strings"

	"github.com/elgopher/pi/cmd/pi/internal/convert/internal/p8"
)

type InputFormat string

const (
	InputFormatP8 = "p8"
)

type Command struct {
	InputFormat
	InputFile  string
	OutputFile string
}

func (o Command) Run() error {
	if o.InputFile == "" {
		return fmt.Errorf("input file not provided")
	}

	if o.InputFormat != "" && o.InputFormat != InputFormatP8 {
		return fmt.Errorf("input format %s not supported", o.InputFormat)
	}

	if o.InputFormat == "" {
		if strings.HasSuffix(o.InputFile, ".p8") {
			o.InputFormat = InputFormatP8
		} else {
			return fmt.Errorf("cannot deduct the format of %s input file", o.InputFile)
		}
	}

	if o.OutputFile == "" {
		return fmt.Errorf("output file not provided")
	}

	fmt.Printf("Converting %s to %s... ", o.InputFile, o.OutputFile)
	fmt.Printf("Using %s input format... ", o.InputFormat)
	if err := o.convert(); err != nil {
		return err
	}
	fmt.Println("Done")
	return nil
}

func (o Command) convert() error {
	if o.InputFormat == InputFormatP8 {
		if err := p8.ConvertToAudioSfx(o.InputFile, o.OutputFile); err != nil {
			return err
		}
	}

	return nil
}
