// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package fileloc

import (
	"runtime"
	"strconv"
	"strings"
)

// Get operation is heavy :( 5 allocations.
func Get(prefixToSkip string) string {
	pcs := make([]uintptr, 5)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	frame, next := frames.Next()
	for ; next; frame, next = frames.Next() {
		name := frame.Func.Name()
		if !strings.HasPrefix(name, prefixToSkip) {
			return frame.File + ":" + strconv.Itoa(frame.Line)
		}
	}

	return ""
}
