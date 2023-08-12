// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package sfmt

import (
	"fmt"
	"strings"
)

func FormatBigSlice[T any](s []T, maxSize int) string {
	var out string
	l := len(s)
	if l > maxSize {
		var b strings.Builder
		for i := 0; i < maxSize; i++ {
			b.WriteString(fmt.Sprintf("%v", s[i]))
			b.WriteString(" ")
		}

		out = fmt.Sprintf("(%d)[%s...]", l, b.String())
	} else {
		out = fmt.Sprintf("%v", s)
	}

	return out
}
