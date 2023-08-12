// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package terminal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
)

var (
	linerState        *liner.State
	Commands          = make(chan string, 1)
	CommandsProcessed = make(chan ProcessingResult, 1)
)

type ProcessingResult int

const (
	Done            ProcessingResult = 0
	MoreInputNeeded ProcessingResult = 1
)

func StartReadingCommands() {
	linerState = liner.NewLiner()
	linerState.SetCtrlCAborts(true)
	linerState.SetBeep(false)

	var prompt = "> "
	var cmd strings.Builder

	go func() {
		for {
			p, err := linerState.Prompt(prompt)
			cmd.WriteString(p)
			if err != nil {
				if errors.Is(err, liner.ErrPromptAborted) {
					fmt.Println("(^D to quit)")
					continue
				}
				if err == io.EOF {
					linerState.Close()
					os.Exit(0)
				}
				linerState.Close()
				panic(err)
			}

			linerState.AppendHistory(p)

			Commands <- cmd.String()
			result := <-CommandsProcessed

			if result == MoreInputNeeded {
				cmd.WriteRune('\n')
				prompt = "  "
			} else {
				cmd.Reset()
				prompt = "> "
			}
		}
	}()
}

func StopReadingCommandsFromStdin() {
	linerState.Close()
}
