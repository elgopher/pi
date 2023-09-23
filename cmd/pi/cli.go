// (c) 2022-2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/elgopher/pi/cmd/pi/internal/convert"
)

func main() {
	app := cli.App{
		Usage: "Pi Command Line Interface",
		Commands: []*cli.Command{
			convertCmd(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
	}
}

func convertCmd() *cli.Command {
	return &cli.Command{
		Name:        "convert",
		Usage:       "Converts one file format into another one",
		Description: "Format of input and output file is deducted based on files extension.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Format of input file. Overrides what CLI deducted based on input file extension. For now, the only supported input format is p8.",
			},
		},
		ArgsUsage: "input.file output.file",
		Action: func(context *cli.Context) error {
			inputFile := context.Args().Get(0)
			outputFile := context.Args().Get(1)

			if context.Args().Len() > 2 {
				return fmt.Errorf("too many arguments")
			}

			command := convert.Command{
				InputFormat: convert.InputFormat(context.String("format")),
				InputFile:   inputFile,
				OutputFile:  outputFile,
			}
			return command.Run()
		},
	}
}
