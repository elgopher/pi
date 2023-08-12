package lib

import (
	"errors"
	"os"
	"strconv"
)

var errRestricted = errors.New("restricted")

// osExit invokes panic instead of exit.
func osExit(code int) { panic("os.Exit(" + strconv.Itoa(code) + ")") }

// osFindProcess returns os.FindProcess, except for self process.
func osFindProcess(pid int) (*os.Process, error) {
	if pid == os.Getpid() {
		return nil, errRestricted
	}
	return os.FindProcess(pid)
}
